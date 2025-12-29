package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/go-chi/httplog/v2"
)

// GetDrillHandler handles the request to get a single drill by img_name and language.
// @Summary Get a single drill
// @Description Get a drill exercise by its image name identifier and language
// @Tags Drills
// @Accept json
// @Produce json
// @Param id query string true "Drill image name (unique identifier)"
// @Param lang query string true "Language code (en, de)"
// @Success 200 {object} models.Drill
// @Failure 400 {string} string "Bad request - missing parameters"
// @Failure 404 {string} string "Drill not found"
// @Failure 500 {string} string "Internal server error"
// @Router /drill [get]
func (rs *RAGService) GetDrillHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Getting drill...")

	// Get query parameters
	imgName := req.URL.Query().Get("id")
	lang := req.URL.Query().Get("lang")

	if imgName == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}
	if lang == "" {
		lang = "en" // Default to English
	}

	httplog.LogEntrySetField(req.Context(), "drill_id", slog.StringValue(imgName))
	httplog.LogEntrySetField(req.Context(), "lang", slog.StringValue(lang))

	drill, err := rs.db.GetDrillByImgName(req.Context(), imgName, lang)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Drill not found", http.StatusNotFound)
			return
		}
		logger.Error("Failed to get drill", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Drill retrieved successfully")
	if err := models.WriteResponseJSON(w, http.StatusOK, drill); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// SearchDrillsHandler handles the request to search drills with filters and pagination.
// @Summary Search drills
// @Description Search drill exercises with optional filters for language, target groups, styles, and difficulty
// @Tags Drills
// @Accept json
// @Produce json
// @Param lang query string true "Language code (en, de)"
// @Param target_groups query []string false "Target groups filter (e.g., Beginner, Competitive Swimmer)"
// @Param styles query []string false "Styles filter (e.g., Freestyle, Backstroke)"
// @Param difficulty query string false "Difficulty filter (Easy, Medium, Hard)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Results per page (default: 20, max: 100)"
// @Success 200 {object} rag.DrillSearchResult
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /drills/search [get]
func (rs *RAGService) SearchDrillsHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Searching drills...")

	// Get query parameters
	lang := req.URL.Query().Get("lang")
	if lang == "" {
		http.Error(w, "lang parameter is required", http.StatusBadRequest)
		return
	}

	// Parse target_groups[] array
	targetGroups := req.URL.Query()["target_groups"]
	// Also support target_groups[] format
	if len(targetGroups) == 0 {
		targetGroups = req.URL.Query()["target_groups[]"]
	}

	// Parse styles[] array
	styles := req.URL.Query()["styles"]
	// Also support styles[] format
	if len(styles) == 0 {
		styles = req.URL.Query()["styles[]"]
	}

	difficulty := req.URL.Query().Get("difficulty")
	searchQuery := req.URL.Query().Get("q")

	// Parse pagination
	page := 1
	if pageStr := req.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := req.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	params := rag.DrillSearchParams{
		Language:     lang,
		TargetGroups: targetGroups,
		Styles:       styles,
		Difficulty:   difficulty,
		SearchQuery:  searchQuery,
		Page:         page,
		Limit:        limit,
	}

	httplog.LogEntrySetField(req.Context(), "lang", slog.StringValue(lang))
	httplog.LogEntrySetField(req.Context(), "page", slog.IntValue(page))
	httplog.LogEntrySetField(req.Context(), "limit", slog.IntValue(limit))
	if searchQuery != "" {
		httplog.LogEntrySetField(req.Context(), "q", slog.StringValue(searchQuery))
	}

	result, err := rs.db.SearchDrills(req.Context(), params)
	if err != nil {
		logger.Error("Failed to search drills", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Drills search completed", "total", result.Total)
	if err := models.WriteResponseJSON(w, http.StatusOK, result); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}
