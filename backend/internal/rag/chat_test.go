package rag

import (
	"context"
	"errors"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
)

func testChatDependencies(t *testing.T) (chatDependencies, *chatDependencyCalls) {
	t.Helper()
	calls := &chatDependencyCalls{}
	plan := &models.Plan{PlanID: "00000000-0000-0000-0000-000000000001", Title: "Owned plan"}
	deps := chatDependencies{
		getPlanForUser: func(_ context.Context, planID, userID string) (*models.Plan, error) {
			calls.getPlan = append(calls.getPlan, []string{planID, userID})
			return plan, nil
		},
		getConversation: func(_ context.Context, planID, userID string) ([]models.Message, error) {
			calls.getConversation = append(calls.getConversation, []string{planID, userID})
			return nil, nil
		},
		buildContext: func(_ context.Context, _ string, _ *models.Plan, _ models.Language) ([]schema.Document, []schema.Document, error) {
			calls.buildContext++
			assert.Equal(t, 1, calls.queryMode, "retrieval must use query embeddings")
			return nil, nil, nil
		},
		queryMode: func() {
			calls.queryMode++
		},
		chatRefine: func(_ context.Context, _ string, _ *models.Plan, _, _ string, _ any, _, _ []schema.Document) (*models.ChatResponse, error) {
			calls.chatRefine++
			return &models.ChatResponse{Response: "response"}, nil
		},
		addMessage: func(_ context.Context, planID, userID string, role models.Role, content string, _ *string, snapshot *models.Plan) (*models.Message, error) {
			calls.addMessage++
			return &models.Message{ID: "message", PlanID: planID, UserID: userID, Role: role, Content: content, PlanSnapshot: snapshot}, nil
		},
		upsertPlan: func(_ context.Context, plan models.Plan, userID string) (string, error) {
			calls.upsertPlan = append(calls.upsertPlan, []string{plan.PlanID, userID})
			return plan.PlanID, nil
		},
	}
	return deps, calls
}

type chatDependencyCalls struct {
	getPlan         [][]string
	getConversation [][]string
	buildContext    int
	queryMode       int
	chatRefine      int
	addMessage      int
	upsertPlan      [][]string
}

func TestChatWithContextStopsBeforeDownstreamWorkForUnauthorizedPlan(t *testing.T) {
	deps, calls := testChatDependencies(t)
	deps.getPlanForUser = func(context.Context, string, string) (*models.Plan, error) {
		return nil, pgx.ErrNoRows
	}

	db := &RAGDB{}
	_, _, err := db.chatWithContext(
		context.Background(),
		"00000000-0000-0000-0000-000000000001",
		"user-b",
		"change it",
		models.LanguageEN,
		25,
		deps,
	)

	require.ErrorIs(t, err, ErrChatPlanNotFound)
	assert.Empty(t, calls.getConversation)
	assert.Zero(t, calls.buildContext)
	assert.Zero(t, calls.queryMode)
	assert.Zero(t, calls.chatRefine)
	assert.Zero(t, calls.addMessage)
	assert.Empty(t, calls.upsertPlan)
}

func TestChatWithContextRequiresPlanBeforeAnyDependency(t *testing.T) {
	deps := chatDependencies{}
	db := &RAGDB{}

	_, _, err := db.chatWithContext(context.Background(), "", "user-a", "hello", models.LanguageEN, 25, deps)

	require.ErrorIs(t, err, ErrChatPlanRequired)
}

func TestChatWithContextUsesAuthenticatedUserAndPlanForOwnerPath(t *testing.T) {
	deps, calls := testChatDependencies(t)
	planID := "00000000-0000-0000-0000-000000000001"
	userID := "user-a"
	db := &RAGDB{}

	updatedPlan, aiMessage, err := db.chatWithContext(
		context.Background(),
		planID,
		userID,
		"make it harder",
		models.LanguageEN,
		25,
		deps,
	)

	require.NoError(t, err)
	require.NotNil(t, updatedPlan)
	require.NotNil(t, aiMessage)
	assert.Equal(t, [][]string{{planID, userID}}, calls.getPlan)
	assert.Equal(t, [][]string{{planID, userID}}, calls.getConversation)
	assert.Equal(t, 1, calls.buildContext)
	assert.Equal(t, 1, calls.queryMode)
	assert.Equal(t, 1, calls.chatRefine)
	assert.Equal(t, 2, calls.addMessage)
	assert.Empty(t, calls.upsertPlan)
	assert.Equal(t, "Owned plan", updatedPlan.Title)
	assert.Same(t, updatedPlan, aiMessage.PlanSnapshot)
}

func TestChatWithContextPreservesPlanUpdateAndSnapshotScope(t *testing.T) {
	deps, calls := testChatDependencies(t)
	planID := "00000000-0000-0000-0000-000000000001"
	userID := "user-a"
	deps.chatRefine = func(_ context.Context, _ string, _ *models.Plan, _, _ string, _ any, _, _ []schema.Document) (*models.ChatResponse, error) {
		calls.chatRefine++
		return &models.ChatResponse{
			Response: "updated",
			Plan:     &models.GeneratedPlan{Title: "Updated plan"},
		}, nil
	}
	var snapshots []*models.Plan
	deps.addMessage = func(_ context.Context, receivedPlanID, receivedUserID string, role models.Role, content string, _ *string, snapshot *models.Plan) (*models.Message, error) {
		calls.addMessage++
		snapshots = append(snapshots, snapshot)
		return &models.Message{ID: "message", PlanID: receivedPlanID, UserID: receivedUserID, Role: role, Content: content}, nil
	}

	db := &RAGDB{}
	updatedPlan, _, err := db.chatWithContext(context.Background(), planID, userID, "update", models.LanguageEN, 25, deps)

	require.NoError(t, err)
	assert.Equal(t, planID, updatedPlan.PlanID)
	assert.Equal(t, [][]string{{planID, userID}}, calls.upsertPlan)
	require.Len(t, snapshots, 2)
	assert.Nil(t, snapshots[0])
	assert.Same(t, updatedPlan, snapshots[1])
	assert.Equal(t, planID, snapshots[1].PlanID)
}

func TestChatWithContextPropagatesUnexpectedAuthorizationError(t *testing.T) {
	deps, _ := testChatDependencies(t)
	backendErr := errors.New("database unavailable")
	deps.getPlanForUser = func(context.Context, string, string) (*models.Plan, error) {
		return nil, backendErr
	}

	_, _, err := (&RAGDB{}).chatWithContext(context.Background(), "plan", "user", "hello", models.LanguageEN, 25, deps)

	require.Error(t, err)
	assert.ErrorIs(t, err, backendErr)
}

func TestChatWithContextPassesGroundingToRefinement(t *testing.T) {
	deps, calls := testChatDependencies(t)
	plans := []schema.Document{{PageContent: "Reference plan"}}
	drills := []schema.Document{{PageContent: "Freestyle drill", Metadata: map[string]any{"slug": "Catch-up", "img_name": "catch-up.webp"}}}
	deps.buildContext = func(_ context.Context, query string, plan *models.Plan, lang models.Language) ([]schema.Document, []schema.Document, error) {
		assert.Equal(t, 1, calls.queryMode)
		assert.Equal(t, "add technique", query)
		assert.Equal(t, "Owned plan", plan.Title)
		assert.Equal(t, models.LanguageEN, lang)
		return plans, drills, nil
	}
	deps.chatRefine = func(_ context.Context, _ string, _ *models.Plan, _, _ string, _ any, planDocs, drillDocs []schema.Document) (*models.ChatResponse, error) {
		calls.chatRefine++
		assert.Equal(t, plans, planDocs)
		assert.Equal(t, drills, drillDocs)
		return &models.ChatResponse{Response: "Which technique?"}, nil
	}

	_, _, err := (&RAGDB{}).chatWithContext(context.Background(), "plan", "user", "add technique", models.LanguageEN, 25, deps)
	require.NoError(t, err)
	assert.Equal(t, 1, calls.chatRefine)
}

func TestChatWithContextStopsWhenGroundingFails(t *testing.T) {
	deps, calls := testChatDependencies(t)
	retrievalErr := errors.New("drill retrieval failed")
	deps.buildContext = func(context.Context, string, *models.Plan, models.Language) ([]schema.Document, []schema.Document, error) {
		return nil, nil, retrievalErr
	}

	_, _, err := (&RAGDB{}).chatWithContext(context.Background(), "plan", "user", "add technique", models.LanguageEN, 25, deps)
	require.ErrorIs(t, err, retrievalErr)
	assert.Zero(t, calls.chatRefine)
	assert.Zero(t, calls.addMessage)
	assert.Empty(t, calls.upsertPlan)
}

func TestBuildChatContextRespectsDisabledRAG(t *testing.T) {
	plans, drills, err := (&RAGDB{}).buildChatContext(context.Background(), "add technique", &models.Plan{}, models.LanguageEN)
	require.NoError(t, err)
	assert.Empty(t, plans)
	assert.Empty(t, drills)
}
