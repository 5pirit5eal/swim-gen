package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/golobby/dotenv"

	"github.com/5pirit5eal/swim-rag/rag"
)

func main() {
	// Configure log to write to stdout
	log.SetOutput(os.Stdout)
	log.Println("Starting server...")

	// ctx := context.Background()
	config := rag.Config{}
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	// ragServer, err := rag.NewRAGServer(ctx, config)
	// defer ragServer.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// mux := http.NewServeMux()
	// mux.HandleFunc("POST /add/", ragServer.AddDocuments)
	// mux.HandleFunc("POST /query/", ragServer.Query)
	// mux.HandleFunc("GET /example/", func(w http.ResponseWriter, r *http.Request) {
	// 	if err := rag.WriteResponseJSON(w, http.StatusOK, rag.Example(ctx, config)); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })
	// port := cmp.Or(os.Getenv("SERVERPORT"), "8080")
	// address := "localhost:" + port
	// log.Println("listening on", address)
	// log.Fatal(http.ListenAndServe(address, mux))
	plans, err := rag.Scrape(make([]string, 0), "https://docswim.de/index.php/2017/07/10/trainingsplan-01-grundlagen-fundament-3-700m/", "https://docswim.de/index.php/2019/09/02/trainingsplan-99-kraulschwimmen-lernen-der-kraul-kurs-teil-2-2-1-700m/")
	log.Println("Scraping produced", plans.Len(), "plans")
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for kvp := range plans.Range() {
		if kvp.Plan.Description == "" {
			// log.Println("Skipping plan with empty title for url: ", kvp.URL)
			continue
		}
		log.Println("Found plan for url: ", kvp.URL)
		fileName := fmt.Sprintf("plan_%d.json", i)
		i++
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(kvp.Plan); err != nil {
			log.Fatal(err)
		}
		log.Printf("Written plan to %s\n", fileName)
	}
}
