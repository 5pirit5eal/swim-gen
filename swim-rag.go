package main

import (
	"context"
	"fmt"
	"log"
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

	fmt.Println(rag.Example(ctx, config))
	dbPass, err := rag.GetDBPass(ctx, config.DB.PassLocation)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB Password:", dbPass)
}
