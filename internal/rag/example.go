package rag

import (
	"context"
	"os"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/golobby/dotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
)

func Example(ctx context.Context, config models.Config) string {
	logger := httplog.LogEntry(ctx)
	logger.Info("Initializing example...")
	file, err := os.Open(".env")
	if err != nil {
		logger.Error("Error opening .env file", httplog.ErrAttr(err))
		panic(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		logger.Error("Error decoding .env file", httplog.ErrAttr(err))
		panic(err)
	}

	llm, err := vertex.New(ctx, googleai.WithCloudProject(config.ProjectID), googleai.WithCloudLocation(config.Region), googleai.WithDefaultModel(config.Model))
	if err != nil {
		logger.Error("Error initializing LLM", httplog.ErrAttr(err))
		panic(err)
	}

	prompt := "Tell me a danish joke in danish"
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		logger.Error("Error generating answer", httplog.ErrAttr(err))
		panic(err)
	}

	logger.Info("Example generated successfully")

	return answer

}
