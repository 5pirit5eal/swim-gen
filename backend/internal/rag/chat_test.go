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
		buildContext: func(_ context.Context, _ string, _ *models.Plan) ([]schema.Document, error) {
			calls.buildContext++
			return nil, nil
		},
		queryMode: func() {
			calls.queryMode++
		},
		chatRefine: func(_ context.Context, _ string, _ *models.Plan, _, _ string, _ any, _ []schema.Document) (*models.ChatResponse, error) {
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
}

func TestChatWithContextPreservesPlanUpdateAndSnapshotScope(t *testing.T) {
	deps, calls := testChatDependencies(t)
	planID := "00000000-0000-0000-0000-000000000001"
	userID := "user-a"
	deps.chatRefine = func(_ context.Context, _ string, _ *models.Plan, _, _ string, _ any, _ []schema.Document) (*models.ChatResponse, error) {
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
