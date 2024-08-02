package main

import (
	"context"
	"github.com/kanataidarov/teams_automator/config"
	"github.com/kanataidarov/teams_automator/internal/handler"
	"log"
)

func main() {
	cfg := config.Load()
	log.Println("Starting Teams Automator")

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	err := handler.Serve(ctx, cfg)
	if err != nil {
		log.Fatalln("Error serving Grpc", err)
	}

	log.Println("Stopping interview_automator backend")
}
