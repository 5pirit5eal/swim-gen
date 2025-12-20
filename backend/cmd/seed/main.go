package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/rag"
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
		fmt.Println("Usage: scrape --lang <language> [--env <env_file>]")
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
	projectRoot := filepath.Dir(*path)

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
		if err := db.Store.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Upload training drills
	// TODO: Implement language specific upload here

	fmt.Println("Upload completed successfully")
}
