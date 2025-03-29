package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"sync"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type RequestQueue struct {
	sync.Mutex
	URLs  []string
	Plans []models.Plan
}

func (rq *RequestQueue) Add(url string, plan models.Plan) {
	rq.Lock()
	defer rq.Unlock()
	rq.URLs = append(rq.URLs, url)
	rq.Plans = append(rq.Plans, plan)
}

func ImprovePlan(ctx context.Context, model llms.Model, p models.Plan, c chan schema.Document, ec chan error) {
	ms, err := models.MetadataSchema()
	if err != nil {
		log.Println("Failed in retrieving Schema")
		ec <- err
		return
	}
	// Enhance scraped documents with gemini and create meaningful metadata
	query := fmt.Sprintf(scrapeTemplateStr, p.Title, p.Description, p.Table.String(), ms)
	answer, err := llms.GenerateFromSinglePrompt(ctx, model, query, llms.WithResponseMIMEType("application/json"))
	if err != nil {
		log.Println("Error when generating answer with LLM:", err)
		ec <- fmt.Errorf("LLM generation error: %w", err)
		return
	}
	// log.Println("Successful answer from LLM: ", answer)

	planMap := p.Map()
	// Parse the answer as JSON
	var metadata models.Metadata
	err = json.Unmarshal([]byte(answer), &metadata)
	if err != nil {
		log.Println("Error parsing LLM response:", err)
		ec <- fmt.Errorf("JSON unmarshal error: %w", err)
		return
	}

	// Add the results to the map
	maps.Copy(planMap, models.StructToMap(metadata))
	// Create request body by converting the models.plans into documents
	c <- schema.Document{
		PageContent: p.String(),
		Metadata:    planMap,
	}
}
