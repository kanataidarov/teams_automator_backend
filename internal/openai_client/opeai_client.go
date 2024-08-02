package openai_client

import (
	"context"
	"errors"
	"github.com/kanataidarov/teams_automator/config"
	ai "github.com/sashabaranov/go-openai"
	"log"
	"os"
)

func Transcribe(ctx context.Context, cfg *config.Config) string {
	filePath := cfg.Grpc.InputFile

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Input file does not exist %s. Error: %v", filePath, err)
		return ""
	}
	log.Printf("Transcribing file - %s", filePath)

	client := ai.NewClient(cfg.OpenAi.Secret)
	resp, err := client.CreateTranscription(
		ctx,
		ai.AudioRequest{
			Model:    cfg.OpenAi.Model,
			FilePath: filePath,
		})
	if err != nil {
		log.Printf("Transcription error. Error: %v", err)
		return ""
	}

	log.Printf("Transcription successful - %s", resp.Text)
	return resp.Text
}
