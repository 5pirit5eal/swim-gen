package main

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/golobby/dotenv"

	"github.com/5pirit5eal/swim-rag/rag"
)

func main() {
	ctx := context.Background()
	config := rag.Config{}
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	ragServer, err := rag.NewRAGServer(ctx, config)
	defer ragServer.Close()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add/", ragServer.AddDocuments)
	mux.HandleFunc("POST /query/", ragServer.Query)
	mux.HandleFunc("GET /example/", func(w http.ResponseWriter, r *http.Request) {
		if err := rag.WriteResponseJSON(w, http.StatusOK, rag.Example(ctx, config)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	port := cmp.Or(os.Getenv("SERVERPORT"), "8080")
	address := "localhost:" + port
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
