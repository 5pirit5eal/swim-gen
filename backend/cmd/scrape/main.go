package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/rag"
)

func main() {
	// Command line flags
	url := flag.String("url", "", "URL to scrape training plans from")
	envFile := flag.String("env", ".env", "path to .env file")
	help := flag.Bool("help", false, "display help information")

	flag.Parse()

	// Display help if requested
	if *help {
		fmt.Println("Scrape training plans from a given URL")
		fmt.Println("Usage: scrape --url <url> [--env <env_file>]")
		fmt.Println("  --url <url>     URL to scrape training plans from")
		fmt.Println("  --env <file>    Path to environment file (default: .env)")
		fmt.Println("  --help          Display this help information")
		os.Exit(0)
	}

	// Validate required parameters
	if *url == "" {
		log.Fatal("Error: URL parameter is required. Use --url to specify the URL to scrape.")
	}

	// Load configuration
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	cfg, err := config.LoadConfig(filepath.Join(projectRoot, *envFile), true)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Initialize context
	ctx := context.Background()

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

	// Perform scraping
	fmt.Printf("Starting to scrape URL: %s\n", *url)
	err = db.ScrapeURL(ctx, *url)
	if err != nil {
		log.Fatal("Error scraping URL:", err)
	}

	fmt.Println("Scraping completed successfully")
}
