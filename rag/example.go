package rag

import (
	"context"
	"log"
	"os"

	"github.com/golobby/dotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
)

func Example(ctx context.Context, config Config) string {
	log.Println("Initializing example...")
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	llm, err := vertex.New(ctx, googleai.WithCloudProject(config.ProjectID), googleai.WithCloudLocation(config.Region), googleai.WithDefaultModel(config.Model))
	if err != nil {
		log.Fatal(err)
	}

	prompt := "Tell me a danish joke in danish"
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Example generated successfully")

	return answer

}
