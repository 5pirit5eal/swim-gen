package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
)

// DrillSearchParams contains parameters for searching drills
type DrillSearchParams struct {
	Language     string   `json:"language"`
	TargetGroups []string `json:"target_groups,omitempty"`
	Styles       []string `json:"styles,omitempty"`
	Difficulty   string   `json:"difficulty,omitempty"`
	Page         int      `json:"page"`
	Limit        int      `json:"limit"`
}

// DrillSearchResult contains paginated drill search results
type DrillSearchResult struct {
	Drills []models.Drill `json:"drills"`
	Total  int            `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

// GetDrillByImgName retrieves a single drill by img_name and language
func (db *RAGDB) GetDrillByImgName(ctx context.Context, imgName, lang string) (*models.Drill, error) {
	logger := httplog.LogEntry(ctx)
	logger.Debug("Getting drill by img_name", "img_name", imgName, "lang", lang)

	query := `
		SELECT cmetadata FROM drill_embeddings
		WHERE cmetadata->>'img_name' = $1
		  AND cmetadata->>'language' = $2
		LIMIT 1
	`

	var metadataJSON []byte
	err := db.Conn.QueryRow(ctx, query, imgName, lang).Scan(&metadataJSON)
	if err != nil {
		logger.Error("Failed to get drill", "error", err)
		return nil, fmt.Errorf("drill not found: %w", err)
	}

	var drill models.Drill
	if err := json.Unmarshal(metadataJSON, &drill); err != nil {
		logger.Error("Failed to unmarshal drill metadata", "error", err)
		return nil, fmt.Errorf("failed to parse drill data: %w", err)
	}

	return &drill, nil
}

// SearchDrills performs a custom SQL query with array filters and pagination
func (db *RAGDB) SearchDrills(ctx context.Context, params DrillSearchParams) (*DrillSearchResult, error) {
	logger := httplog.LogEntry(ctx)
	logger.Debug("Searching drills", "params", params)

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	// Build dynamic query
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Language filter (required)
	conditions = append(conditions, fmt.Sprintf("cmetadata->>'language' = $%d", argIndex))
	args = append(args, params.Language)
	argIndex++

	// Difficulty filter (optional)
	if params.Difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("cmetadata->>'difficulty' = $%d", argIndex))
		args = append(args, params.Difficulty)
		argIndex++
	}

	// Target groups filter (array containment)
	if len(params.TargetGroups) > 0 {
		targetGroupsJSON, _ := json.Marshal(params.TargetGroups)
		conditions = append(conditions, fmt.Sprintf("cmetadata->'target_groups' @> $%d::jsonb", argIndex))
		args = append(args, string(targetGroupsJSON))
		argIndex++
	}

	// Styles filter (array containment)
	if len(params.Styles) > 0 {
		stylesJSON, _ := json.Marshal(params.Styles)
		conditions = append(conditions, fmt.Sprintf("cmetadata->'styles' @> $%d::jsonb", argIndex))
		args = append(args, string(stylesJSON))
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total matching drills
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM drill_embeddings WHERE %s", whereClause)
	var total int
	err := db.Conn.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		logger.Error("Failed to count drills", "error", err)
		return nil, fmt.Errorf("failed to count drills: %w", err)
	}

	// Fetch paginated results
	offset := (params.Page - 1) * params.Limit
	dataQuery := fmt.Sprintf(`
		SELECT cmetadata FROM drill_embeddings
		WHERE %s
		ORDER BY cmetadata->>'title'
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	rows, err := db.Conn.Query(ctx, dataQuery, args...)
	if err != nil {
		logger.Error("Failed to search drills", "error", err)
		return nil, fmt.Errorf("failed to search drills: %w", err)
	}
	defer rows.Close()

	var drills []models.Drill
	for rows.Next() {
		var metadataJSON []byte
		if err := rows.Scan(&metadataJSON); err != nil {
			logger.Warn("Failed to scan drill row", "error", err)
			continue
		}

		var drill models.Drill
		if err := json.Unmarshal(metadataJSON, &drill); err != nil {
			logger.Warn("Failed to unmarshal drill metadata", "error", err)
			continue
		}
		drills = append(drills, drill)
	}

	if drills == nil {
		drills = []models.Drill{}
	}

	return &DrillSearchResult{
		Drills: drills,
		Total:  total,
		Page:   params.Page,
		Limit:  params.Limit,
	}, nil
}

// QueryDrillsForPlan searches for relevant drills based on user query for plan generation
func (db *RAGDB) QueryDrillsForPlan(ctx context.Context, query string, lang string, limit int) ([]models.Drill, error) {
	logger := httplog.LogEntry(ctx)
	logger.Debug("Querying drills for plan", "query", query, "lang", lang, "limit", limit)

	// Set query mode for embedder
	db.Client.QueryMode()

	// Use the DrillStore for similarity search
	docs, err := db.DrillStore.SimilaritySearch(ctx, query, limit)
	if err != nil {
		logger.Error("Failed to search drills", "error", err)
		return nil, fmt.Errorf("failed to search drills: %w", err)
	}

	// Convert documents to drills, filtering by language
	drills := make([]models.Drill, 0, len(docs))
	for _, doc := range docs {
		drill := docToDrill(doc)
		if drill != nil && drill.Language == lang {
			drills = append(drills, *drill)
		}
	}

	// If we didn't get enough drills in the target language, search more
	if len(drills) < limit && len(docs) > 0 {
		// Try a larger search to get more results
		moreDocs, err := db.DrillStore.SimilaritySearch(ctx, query, limit*3)
		if err == nil {
			for _, doc := range moreDocs {
				drill := docToDrill(doc)
				if drill != nil && drill.Language == lang {
					// Check if we already have this drill
					isDuplicate := false
					for _, existing := range drills {
						if existing.ImgName == drill.ImgName {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						drills = append(drills, *drill)
						if len(drills) >= limit {
							break
						}
					}
				}
			}
		}
	}

	logger.Debug("Found drills for plan", "count", len(drills))
	return drills, nil
}

// docToDrill converts a langchain document to a Drill
func docToDrill(doc schema.Document) *models.Drill {
	metadata := doc.Metadata
	if metadata == nil {
		return nil
	}

	drill := &models.Drill{
		Slug:             getStringFromMetadata(metadata, "slug"),
		ShortDescription: getStringFromMetadata(metadata, "short_description"),
		ImgName:          getStringFromMetadata(metadata, "img_name"),
		ImgDescription:   getStringFromMetadata(metadata, "img_description"),
		Title:            getStringFromMetadata(metadata, "title"),
		Difficulty:       getStringFromMetadata(metadata, "difficulty"),
		Language:         getStringFromMetadata(metadata, "language"),
		Targets:          getStringSliceFromMetadata(metadata, "targets"),
		Description:      getStringSliceFromMetadata(metadata, "description"),
		VideoURL:         getStringSliceFromMetadata(metadata, "video_url"),
		Styles:           getStringSliceFromMetadata(metadata, "styles"),
		TargetGroups:     getStringSliceFromMetadata(metadata, "target_groups"),
	}

	return drill
}

func getStringFromMetadata(metadata map[string]any, key string) string {
	if val, ok := metadata[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getStringSliceFromMetadata(metadata map[string]any, key string) []string {
	if val, ok := metadata[key]; ok {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}
