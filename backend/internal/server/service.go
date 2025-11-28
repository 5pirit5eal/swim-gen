package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/pdf"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type RAGService struct {
	// Background context for the server
	ctx context.Context
	// Database client used for storing and querying documents
	db *rag.RAGDB
	// Supabase client for authentication
	auth *supabase.Client
	// Configuration for the RAG server
	cfg config.Config
}

// Initializes a new RAG service with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGService and an error if any occurred during
// initialization.
func NewRAGService(ctx context.Context, cfg config.Config) (*RAGService, error) {
	slog.Info("Initializing RAG server with config", "cfg", slog.AnyValue(cfg))
	db, err := rag.NewGoogleAIStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("Created database connection successfully")

	auth, err := supabase.NewClient(cfg.SB.ApiUrl, cfg.SB.AnonKey, nil)
	if err != nil {
		fmt.Println("Failed to initialize the client: ", err)
	}

	slog.Info("Initialized Supabase client successfully")

	return &RAGService{
		ctx:  ctx,
		cfg:  cfg,
		db:   db,
		auth: auth,
	}, nil
}

// Closes the database connection and LLM client.
// It is important to call this method when the service is no longer needed
// to release resources and avoid memory leaks.
func (rs *RAGService) Close() {
	slog.Info("Closing RAG server...")
	if err := rs.db.Store.Close(); err != nil {
		slog.Error("Error closing database connection", "err", err.Error())
	}
	slog.Info("RAG server closed successfully")
}

// DonatePlanHandler handles the HTTP request to donate a training plan to the database.
// It parses the request, stores the documents and their embeddings in the
// database, and responds with a success message.
// @Summary Donate a new training plan
// @Description Upload and store a new user created swim training plan in the RAG system
// @Tags Donation
// @Accept json
// @Produce json
// @Param plan body models.DonatePlanRequest true "Training plan data"
// @Success 200 {string} string "Plan added successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /add [post]
func (rs *RAGService) DonatePlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding documents to the database...")
	// Parse HTTP request from JSON.

	dpr := &models.DonatePlanRequest{}

	err := models.GetRequestJSON(req, dpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the table is filled
	if len(dpr.Table) == 0 {
		http.Error(w, "Table is empty", http.StatusBadRequest)
		return
	}

	desc := &models.Description{}
	// Check if description is empty and generate one if needed
	if dpr.Description == "" || dpr.Title == "" {
		// Generate a description for the plan
		desc, err := rs.db.Client.DescribeTable(req.Context(), &dpr.Table)
		if err != nil {
			logger.Error("Error when generating description with LLM", httplog.ErrAttr(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch {
		case dpr.Title != "":
			desc.Title = dpr.Title
			fallthrough
		case dpr.Description != "":
			desc.Text = dpr.Description
		}
	} else {
		// Generate metadata with improve plan
		m, err := rs.db.Client.GenerateMetadata(req.Context(), &models.Plan{Title: dpr.Title, Description: dpr.Description, Table: dpr.Table})
		if err != nil {
			logger.Error("Error when generating metadata with LLM", httplog.ErrAttr(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		desc.Meta = m
	}

	// Create a donated plan
	plan := &models.DonatedPlan{
		UserID:      req.Context().Value(models.UserIdCtxKey).(string),
		PlanID:      uuid.NewString(),
		CreatedAt:   time.Now().Format(time.DateTime),
		Title:       desc.Title,
		Description: desc.Text,
		Table:       dpr.Table,
	}

	// Store the plan in the database
	err = rs.db.AddDonatedPlan(req.Context(), plan, desc.Meta)
	if err != nil {
		logger.Error("Failed to store plan in the database", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Scraping completed successfully")); err != nil {
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

	p, err := rs.db.Query(req.Context(), qr.Content, qr.Language, qr.Filter, qr.Method, qr.PoolLength)
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

	userId := req.Context().Value(models.UserIdCtxKey).(string)
	if userId != "" {
		// Add a plan id to the newly created plan
		p.PlanID = uuid.NewString()
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

// PlanToPDFHandler handles the Plan to PDF export request.
// @Summary Export training plan to PDF
// @Description Generate and download a PDF version of a training plan
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param plan body models.PlanToPDFRequest true "Training plan data to export"
// @Success 200 {object} models.PlanToPDFResponse "PDF export response with URI"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /export-pdf [post]
func (rs *RAGService) PlanToPDFHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Exporting table to PDF...")

	// Parse HTTP request from JSON.
	qr := &models.PlanToPDFRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Increment the export count for the user profile if UserID is provided
	if qr.PlanID != "" {
		err = rs.db.IncrementExportCount(req.Context(), req.Context().Value(models.UserIdCtxKey).(string), qr.PlanID)
		if err != nil {
			logger.Error("Failed to increment export count", httplog.ErrAttr(err))
		}
	}

	// Convert the table to PDF
	planPDF, err := pdf.PlanToPDF(
		&models.Plan{
			Title:       qr.Title,
			Description: qr.Description,
			Table:       qr.Table,
		},
		qr.Horizontal,
		qr.LargeFont,
		qr.Language,
	)
	if err != nil {
		logger.Error("Table generation failed", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload the PDF to cloud storage
	uri, err := pdf.UploadPDF(req.Context(), rs.cfg.Bucket.ServiceAccount, rs.cfg.Bucket.Name, pdf.GenerateFilename(), planPDF)
	if err != nil {
		logger.Error("PDF upload failed", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answer := &models.PlanToPDFResponse{URI: uri}

	logger.Info("Answer generated successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, answer); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

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

// UpsertPlan upserts a plan into the users history.
// If the plan exists and belongs to the user, it updates the plan.
// Otherwise it inserts a new plan for the user.
// @Summary Update or insert a training plan into a user's history
// @Description Update an existing training plan if it belongs to the user, or insert a new one
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param request body models.UpsertPlanRequest true "Request to upsert a training plan"
// @Success 200 {object} models.UpsertPlanResponse "Plan ID of the upserted training plan"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /upsert-plan [post]
func (rs *RAGService) UpsertPlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Upserting plan into the database...")
	upr := &models.UpsertPlanRequest{}

	err := models.GetRequestJSON(req, upr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := req.Context().Value(models.UserIdCtxKey).(string)
	if userId == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}
	logger.Debug("Upserting plan into db")
	resp, err := rs.db.UpsertPlan(req.Context(), models.Plan{
		PlanID:      upr.PlanID,
		Title:       upr.Title,
		Description: upr.Description,
		Table:       upr.Table,
	}, userId)
	if err != nil {
		logger.Error("Failed to upsert plan in the database", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the plan ID
	answer := &models.UpsertPlanResponse{PlanID: resp}
	logger.Info("Plan upserted successfully", "plan_id", resp)
	if err := models.WriteResponseJSON(w, http.StatusOK, answer); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// AddPlanToHistoryHandler handles the request to add a new plan to a user's history.
// This is used when a user wants to save a plan snapshot from conversation history with a new PlanID.
// @Summary Add a plan to user history
// @Description Add a plan to the authenticated user's history with a new id
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param request body models.AddPlanToHistoryRequest true "Plan to add to history"
// @Success 200 {object} models.AddPlanToHistoryResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /add-plan-to-history [post]
func (rs *RAGService) AddPlanToHistoryHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding plan to user history...")

	// Parse request body
	var plan models.AddPlanToHistoryRequest
	err := models.GetRequestJSON(req, &plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := req.Context().Value(models.UserIdCtxKey).(string)
	if userID == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	// Generate a new PlanID for the snapshot
	plan.PlanID = uuid.NewString()

	// Add to user history
	err = rs.db.AddPlanToHistory(req.Context(), plan.Plan(), userID)
	if err != nil {
		logger.Error("Failed to add plan to user history", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success and the new PlanID
	response := models.AddPlanToHistoryResponse{
		Message: "Plan added to history successfully",
		PlanID:  plan.PlanID,
	}
	logger.Info("Plan added to history successfully", "plan_id", plan.PlanID)
	if err := models.WriteResponseJSON(w, http.StatusOK, response); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// SharePlanHandler handles the request to share a training plan.
// It generates a shareable url_hash or processes email sharing based on the method provided.
// @Summary Share a training plan
// @Description Share a training plan via link or email. Email sharing is not implemented yet.
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param request body models.SharePlanRequest true "Request to share a training plan"
// @Success 200 {object} models.SharePlanResponse "Share plan response with URI"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /share-plan [post]
func (rs *RAGService) SharePlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Sharing plan...")
	spr := &models.SharePlanRequest{}

	err := models.GetRequestJSON(req, spr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := req.Context().Value(models.UserIdCtxKey).(string)
	if userId == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	url_hash, err := rs.db.SharePlan(req.Context(), spr.PlanID, userId, spr.Method)
	if err != nil {
		logger.Error("Failed to share plan", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the shareable URI
	answer := &models.SharePlanResponse{URLHash: url_hash}
	logger.Info("Plan shared successfully", "uri", url_hash)
	if err := models.WriteResponseJSON(w, http.StatusOK, answer); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}
