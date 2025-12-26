package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/pdf"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
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
	if err := rs.db.PlanStore.Close(); err != nil {
		slog.Error("Error closing database connection", "err", err.Error())
	}
	if err := rs.db.DrillStore.Close(); err != nil {
		slog.Error("Error closing drill store connection", "err", err.Error())
	}
	slog.Info("RAG server closed successfully")
}

// getMimeTypeFromFilename returns the MIME type based on the file extension.
// Returns an error if the file type is not supported.
func getMimeTypeFromFilename(filename string) (string, error) {
	filename = strings.ToLower(filename)
	switch {
	case strings.HasSuffix(filename, ".png"):
		return "image/png", nil
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg", nil
	case strings.HasSuffix(filename, ".pdf"):
		return "application/pdf", nil
	default:
		return "", fmt.Errorf("unsupported file type: %s. Supported formats: PNG, JPEG, PDF", filename)
	}
}

// UploadPlanHandler handles the HTTP request to upload a private training plan to the database.
// It parses the request, stores the documents and their embeddings in the
// database, and responds with a success message.
// @Summary Upload a new private training plan
// @Description Upload and store a new user created swim training plan in the database
// @Tags Upload
// @Accept json
// @Produce json
// @Param plan body models.UploadPlanRequest true "Training plan data"
// @Success 200 {string} string "Plan added successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /add [post]
func (rs *RAGService) UploadPlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding uploaded plan to the users history...")
	// Parse HTTP request from JSON.

	dpr := &models.UploadPlanRequest{}

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
		var err error
		desc, err = rs.db.Client.DescribeTable(req.Context(), &dpr.Table)
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
		desc.Title = dpr.Title
		desc.Text = dpr.Description
	}

	// Create a donated plan
	plan := &models.DonatedPlan{
		UserID:       req.Context().Value(models.UserIdCtxKey).(string),
		PlanID:       uuid.NewString(),
		Title:        desc.Title,
		Description:  desc.Text,
		Table:        dpr.Table,
		AllowSharing: dpr.AllowSharing,
	}

	// Store the plan in the database
	err = rs.db.AddUploadedPlan(req.Context(), plan)
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

// FileToPlanHandler handles the request to convert a file (image or PDF) of a plan to a plan
// The file is sent as form data. Supported formats: PNG, JPEG, PDF
// @Summary Convert a file (image or PDF) of a plan to a plan
// @Description Convert a file containing a training plan to a structured plan. Supports PNG, JPEG, and PDF formats.
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "File containing a plan (PNG, JPEG, or PDF)"
// @Success 200 {object} models.RAGResponse "Plan ID of the converted plan"
// @Failure 400 {string} string "Bad request or unsupported file type"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /file-to-plan [post]
func (rs *RAGService) FileToPlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Request for file to plan received...")

	// 1. tell Go to parse the incoming multipart stream
	err := req.ParseMultipartForm(20 << 20) // 20 MB max memory
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. retrieve the file (form field name must match the client's key)
	file, header, err := req.FormFile("file")
	logger.Debug("Filename", "filename", header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	// 3. Detect MIME type from filename
	mimeType, err := getMimeTypeFromFilename(header.Filename)
	if err != nil {
		logger.Error("Unsupported file type", "filename", header.Filename, httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Debug("Detected MIME type", "mimeType", mimeType)

	// read the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("Converting file to plan")
	// Get language from form data
	language := req.FormValue("language")
	if language == "" {
		language = "en"
	}

	resp, err := rs.db.Client.FileToPlan(req.Context(), fileBytes, header.Filename, mimeType, models.Language(language))
	if err != nil {
		logger.Error("Failed to convert file to plan in the database", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answer := &models.RAGResponse{
		Title:       resp.Title,
		Description: resp.Description,
		Table:       resp.Table,
	}

	logger.Info("Image converted to plan successfully", "plan_id", resp)
	if err := models.WriteResponseJSON(w, http.StatusOK, answer); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// GetUploadedPlansHandler handles the request to get all uploaded plans for a user.
// @Summary Get uploaded plans
// @Description Get all plans uploaded by the authenticated user
// @Tags Upload
// @Accept json
// @Produce json
// @Success 200 {array} models.DonatedPlan
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /uploads [get]
func (rs *RAGService) GetUploadedPlansHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Getting uploaded plans...")

	userId := req.Context().Value(models.UserIdCtxKey).(string)
	if userId == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}
	httplog.LogEntrySetField(req.Context(), "user_id", slog.StringValue(userId))

	plans, err := rs.db.GetUploadedPlans(req.Context(), userId)
	if err != nil {
		logger.Error("Failed to get uploaded plans", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Uploaded plans retrieved successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, plans); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// GetUploadedPlanHandler handles the request to get a specific uploaded plan.
// @Summary Get a uploaded plan
// @Description Get a specific plan uploaded by the authenticated user
// @Tags Upload
// @Accept json
// @Produce json
// @Param plan_id path string true "Plan ID"
// @Success 200 {object} models.DonatedPlan
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Plan not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /uploads/{plan_id} [get]
func (rs *RAGService) GetUploadedPlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Getting uploaded plan...")

	planID := chi.URLParam(req, "plan_id")
	if planID == "" {
		http.Error(w, "Plan ID is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(planID))

	plan, err := rs.db.GetUploadedPlan(req.Context(), planID)
	if err != nil {
		logger.Error("Failed to get uploaded plan", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Uploaded plan retrieved successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, plan); err != nil {
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
		httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(qr.PlanID))
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

	// Determine storage path
	var username string
	userID := req.Context().Value(models.UserIdCtxKey).(string)
	if userID != "" {
		profile, err := rs.db.GetUserProfile(req.Context(), userID)
		if err != nil {
			logger.Warn("Failed to get user profile for PDF path generation", httplog.ErrAttr(err))
		} else if profile != nil {
			username = profile.Username
		}
	}

	storagePath := pdf.GenerateStoragePath(username, qr.PlanID, qr.Title)

	// Upload the PDF to cloud storage
	uri, err := pdf.UploadPDF(req.Context(), rs.cfg.Bucket.ServiceAccount, rs.cfg.Bucket.Name, storagePath, planPDF)
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
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(resp))
	logger.Info("Plan upserted successfully")
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
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(plan.PlanID))

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
	logger.Info("Plan added to history successfully")
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
	if spr.PlanID != "" {
		httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(spr.PlanID))
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

// FeedbackHandler handles the request to submit feedback for a training plan.
// @Summary Submit feedback for a training plan
// @Description Submit a rating, was_swam status, and difficulty rating for a training plan
// @Tags Feedback
// @Accept json
// @Produce json
// @Param request body models.FeedbackRequest true "Feedback data"
// @Success 200 {string} string "Feedback submitted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /feedback [post]
func (rs *RAGService) FeedbackHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Submitting feedback...")

	// Parse HTTP request from JSON.
	fr := &models.FeedbackRequest{}
	err := models.GetRequestJSON(req, fr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if fr.PlanID != "" {
		httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(fr.PlanID))
	}

	userId := req.Context().Value(models.UserIdCtxKey).(string)
	if userId == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	feedback := &models.Feedback{
		UserID:           userId,
		PlanID:           fr.PlanID,
		Rating:           fr.Rating,
		WasSwam:          fr.WasSwam,
		DifficultyRating: fr.DifficultyRating,
		Comment:          fr.Comment,
	}

	// Check if feedback already exists
	existingFeedback, err := rs.db.GetFeedback(req.Context(), userId, fr.PlanID)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		logger.Error("Error checking for existing feedback", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingFeedback != nil {
		// Update existing feedback
		err = rs.db.UpdateFeedback(req.Context(), feedback)
	} else {
		// Add new feedback
		err = rs.db.AddFeedback(req.Context(), feedback)
	}

	if err != nil {
		logger.Error("Failed to submit feedback", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Feedback submitted successfully")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Feedback submitted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// DeletePlanHandler handles the request to delete a plan from the user's history.
// @Summary Delete a training plan
// @Description Delete a training plan from the authenticated user's history and remove all associated data
// @Tags Training Plans
// @Accept json
// @Produce json
// @Param plan_id path string true "Plan ID to delete"
// @Success 200 {string} string "Plan deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Plan not found"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /plan/{plan_id} [delete]
func (rs *RAGService) DeletePlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Deleting plan...")

	planID := chi.URLParam(req, "plan_id")
	if planID == "" {
		http.Error(w, "Plan ID is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(planID))

	userID := req.Context().Value(models.UserIdCtxKey).(string)
	if userID == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	err := rs.db.DeletePlan(req.Context(), planID, userID)
	if err != nil {
		if err.Error() == "plan not found in user history or user does not own the plan" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Error("Failed to delete plan", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Plan deleted successfully")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Plan deleted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// DeleteUserHandler handles the request to delete a user account and all associated data.
// This operation deletes the user from auth.users which triggers CASCADE deletion of all related data.
// @Summary Delete user account
// @Description Permanently delete the authenticated user's account and all associated data
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {string} string "User deleted successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /user [delete]
func (rs *RAGService) DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Deleting user account...")

	userID := req.Context().Value(models.UserIdCtxKey).(string)
	if userID == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	// Parse UUID for the admin auth client
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error("Invalid user ID format", httplog.ErrAttr(err))
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Delete user via Supabase Admin API using service role key
	// This will CASCADE delete all related data in: profiles, history, donations, feedback, shared_plans, shared_history, memory
	err = rs.auth.Auth.WithToken(rs.cfg.SB.ServiceRoleKey).AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: userUUID,
	})
	if err != nil {
		logger.Error("Failed to delete user", httplog.ErrAttr(err))
		http.Error(w, "Failed to delete user account", http.StatusInternalServerError)
		return
	}

	logger.Info("User deleted successfully", "user_id", userID)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("User deleted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}
