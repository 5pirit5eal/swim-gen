package main

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golobby/dotenv"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/5pirit5eal/swim-rag/internal/server"
)

func main() {
	// Configure log to write to stdout
	log.SetOutput(os.Stdout)
	log.Println("Starting server...")

	ctx := context.Background()
	config := models.Config{}
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(filepath.Join(projectRoot, ".env"))
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	ragServer, err := server.NewRAGServer(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer ragServer.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", ragServer.AddDocuments)
	mux.HandleFunc("POST /query", ragServer.Query)
	mux.HandleFunc("GET /scrape", ragServer.Scrape)
	// mux.HandleFunc("GET /example", func(w http.ResponseWriter, r *http.Request) {
	// 	if err := models.WriteResponseJSON(w, http.StatusOK, rag.Example(ctx, config)); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })
	port := cmp.Or(os.Getenv("SERVERPORT"), "8080")
	address := "localhost:" + port
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
