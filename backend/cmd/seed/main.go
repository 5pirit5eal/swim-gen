package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/tmc/langchaingo/schema"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func main() {
	// Command line flags
	path := flag.String("path", "", "Path to the directory containing training drills")
	var lang arrayFlags
	flag.Var(&lang, "lang", "Language for the training drills")
	envFile := flag.String("env", ".env", "Name of .env file")
	help := flag.Bool("help", false, "display help information")

	flag.Parse()

	// Display help if requested
	if *help {
		fmt.Println("Upload training drills to the exercise database")
		fmt.Println("Usage: seed --lang <language> [--env <env_file>]")
		fmt.Println("  --lang <language>  Languages for the training drills, can be specified multiple times")
		fmt.Println("  --env <file>       Name of environment file (default: .env)")
		fmt.Println("  --help             Display this help information")
		os.Exit(0)
	}

	// Validate required parameters
	if len(lang) == 0 {
		// Check the data directory for available languages and use all of them
		entries, err := os.ReadDir(*path)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			lang = append(lang, strings.TrimSuffix(strings.ToLower(e.Name()), ".json"))
		}
	}

	// Load configuration
	// Traverse two levels up from the drills path to get the project root (drills directory -> parent -> project root).
	drillsParentDir := filepath.Dir(*path)
	projectRoot := filepath.Dir(drillsParentDir)

	// Expect the .env file to be in the backend directory
	cfg, err := config.LoadConfig(filepath.Join(projectRoot, "backend", *envFile), true)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Initialize context
	ctx := context.Background()
	defer ctx.Done()

	// Initialize RAG database
	db, err := rag.NewGoogleAIStore(ctx, cfg)
	if err != nil {
		log.Fatal("Error initializing RAG database:", err)
	}
	defer func() {
		if err := db.PlanStore.Close(); err != nil {
			log.Printf("Error closing plan store connection: %v", err)
		}
		if err := db.DrillStore.Close(); err != nil {
			log.Printf("Error closing drill store connection: %v", err)
		}
	}()

	// Upload training drills
	for _, l := range lang {
		fmt.Printf("Processing language: %s\n", l)

		// Read the JSON file
		filePath := filepath.Join(*path, l+".json")
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
			continue
		}

		var drills []models.Drill
		if err := json.Unmarshal(content, &drills); err != nil {
			log.Printf("Error unmarshaling JSON from %s: %v", filePath, err)
			continue
		}

		var documents []schema.Document
		for _, drill := range drills {
			// Skip empty drills
			if drill.Title == "" {
				continue
			}

			doc := schema.Document{
				PageContent: fmt.Sprintf("Title: %s\nDescription: %s\nShort Description: %s\nTargets: %s\nStyles: %s\nDifficulty: %s",
					drill.Title,
					strings.Join(drill.Description, " "),
					drill.ShortDescription,
					strings.Join(drill.Targets, ", "),
					strings.Join(drill.Styles, ", "),
					drill.Difficulty,
				),
				Metadata: map[string]any{
					"slug":              drill.Slug,
					"targets":           drill.Targets,
					"short_description": drill.ShortDescription,
					"img_name":          drill.ImgName,
					"img_description":   drill.ImgDescription,
					"title":             drill.Title,
					"description":       drill.Description,
					"video_url":         drill.VideoURL,
					"styles":            drill.Styles,
					"difficulty":        drill.Difficulty,
					"target_groups":     drill.TargetGroups,
					"language":          l,
				},
			}
			documents = append(documents, doc)
		}

		if len(documents) > 0 {
			fmt.Printf("Uploading %d drills for language %s...\n", len(documents), l)
			if _, err := db.DrillStore.AddDocuments(ctx, documents); err != nil {
				log.Printf("Error uploading documents for language %s: %v", l, err)
			} else {
				fmt.Printf("Successfully uploaded %d drills for language %s\n", len(documents), l)
			}
		} else {
			fmt.Printf("No valid drills found for language %s\n", l)
		}
	}

	fmt.Println("Upload completed successfully")
}
