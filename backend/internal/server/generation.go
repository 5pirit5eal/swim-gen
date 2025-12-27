package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

// GeneratePromptHandler handles the request to generate a prompt for the LLM.
// It uses the GoogleGenAIClient to generate a prompt based on the provided language.
// @Summary Generate a prompt for the LLM
// @Description Generate a prompt for the LLM based on the provided language
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param request body models.GeneratePromptRequest true "Request to generate a prompt"
// @Success 200 {object} models.GeneratedPromptResponse "Generated prompt response"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /generate-prompt [post]
func (rs *RAGService) GeneratePromptHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Generating prompt...")

	// Parse HTTP request from JSON.
	gpr := &models.GeneratePromptRequest{}
	err := models.GetRequestJSON(req, gpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate the prompt using the GoogleGenAIClient
	prompt, err := rs.db.Client.GeneratePrompt(req.Context(), *gpr)
	if err != nil {
		logger.Error("Error generating prompt", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &models.GeneratedPromptResponse{Prompt: prompt}
	logger.Info("Prompt generated successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, response); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// QueryHandler handles the RAG query request.
// It parses the request, queries the RAG, generating or choosing a plan, and returns the result as JSON.
// @Summary Query training plans
// @Description Query the RAG system for relevant training plans based on input
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param query body models.QueryRequest true "Query parameters"
// @Success 200 {object} models.RAGResponse "Query results"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /query [post]
func (rs *RAGService) QueryHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Querying the database...")
	// Parse HTTP request from JSON.

	qr := &models.QueryRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := req.Context().Value(models.UserIdCtxKey).(string)

	var userProfileStr string
	// Check if preferences should be used (default to true)
	usePreferences := true
	if qr.Preferences != nil {
		usePreferences = *qr.Preferences
	}

	if usePreferences && userId != "" {
		profile, err := rs.db.GetUserProfile(req.Context(), userId)
		if err != nil {
			logger.Warn("Failed to get user profile, proceeding without it", httplog.ErrAttr(err))
		} else {
			userProfileStr = rs.db.FormatUserProfile(profile)
		}
	}

	p, err := rs.db.Query(req.Context(), qr.Content, qr.Language, userProfileStr, qr.Filter, qr.Method, qr.PoolLength)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unsupported method:") {
			http.Error(w, "Method may only be 'choose' or 'generate', invalid choice.", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Recalculate the sums of the rows to be sure they are correct
	p.Table.UpdateSum()
	logger.Debug("Updated the table sums...", "sum", p.Table[len(p.Table)-1].Sum)

	// userId is already declared above
	if userId != "" {
		// Add a plan id to the newly created plan
		p.PlanID = uuid.NewString()
		httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(p.PlanID))
		logger.Info("Adding plan to user history", "user_id", userId, "plan_id", p.PlanID)
		err = rs.db.AddPlanToHistory(req.Context(), p, userId)
		if err != nil {
			logger.Error("Failed to add plan to user history", httplog.ErrAttr(err))
		} else {
			// Add the initial conversation to the memory
			// 1. User message
			userMsg, err := rs.db.Memory.AddMessage(req.Context(), p.PlanID, userId, models.RoleUser, qr.Content, nil, nil)
			if err != nil {
				logger.Error("Failed to add user message to memory", httplog.ErrAttr(err))
			} else {
				// 2. AI message with plan snapshot
				_, err = rs.db.Memory.AddMessage(req.Context(), p.PlanID, userId, models.RoleAI, p.Description, &userMsg.ID, p)
				if err != nil {
					logger.Error("Failed to add AI message to memory", httplog.ErrAttr(err))
				}
			}
		}
	}

	// Convert to response payload
	answer := &models.RAGResponse{
		PlanID:      p.PlanID,
		Title:       p.Title,
		Description: p.Description,
		Table:       p.Table,
	}

	logger.Info("Answer generated successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, answer); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}
