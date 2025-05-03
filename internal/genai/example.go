package genai

import (
	"context"

	"github.com/go-chi/httplog/v2"
	"google.golang.org/genai"
)

func (c *GoogleGenAIClient) GenerateExample(ctx context.Context) (string, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("Initializing example...")

	prompt := "Tell me a danish joke in danish"
	answer, err := c.gc.Models.GenerateContent(ctx, c.cfg.Model, genai.Text(prompt), c.gcfg)
	if err != nil {
		logger.Error("Error generating answer", httplog.ErrAttr(err))
		return "", err
	}

	logger.Info("Example generated successfully")
	return answer.Text(), nil
}
