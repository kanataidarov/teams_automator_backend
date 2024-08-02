package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kanataidarov/teams_automator/config"
	"github.com/kanataidarov/teams_automator/internal/openai_client"
	pb "github.com/kanataidarov/teams_automator/pkg/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
)

type server struct {
	pb.UnimplementedOpenAiApiServer
}

func Serve(_ context.Context, cfg *config.Config) error {
	port := cfg.Grpc.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d. Error: %v", port, err)
	}

	srv := grpc.NewServer()
	pb.RegisterOpenAiApiServer(srv, &server{})

	log.Printf("Grpc server listening on port %d", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on port %d. Error: %v", port, err)
	}

	return nil
}

func (s *server) Transcribe(ctx context.Context, req *pb.TranscribeRequest) (*pb.TranscribeResponse, error) {
	cfg := config.Load()

	var (
		header *pb.FileHeader
		buf    bytes.Buffer
	)

	header = req.Header
	fileName := filepath.Base(header.Name)
	log.Printf("File recevied - %s, size should be - %d", fileName, header.Size)

	if data := req.Data; data != nil {
		buf.Write(data)
	}
	log.Printf("File received. Size %d Kb", buf.Len()/1024)

	transcription := doTranscribe(ctx, cfg, buf)

	return &pb.TranscribeResponse{Transcription: transcription}, nil

}

func doTranscribe(ctx context.Context, cfg *config.Config, buf bytes.Buffer) string {

	prepareFile(cfg, buf)

	return openai_client.Transcribe(ctx, cfg)
}

func prepareFile(cfg *config.Config, buf bytes.Buffer) *os.File {
	filePath := cfg.Grpc.InputFile

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating file %s. Error: %v", filePath, err)
	}

	length, err := file.Write(buf.Bytes())
	if err != nil {
		_ = file.Close()
		log.Fatalf("Couldn't write to file %s. Error: %v", filePath, err)
	}
	log.Printf("Wrote %d Kb to file %s", length/1024, filePath)

	err = file.Close()
	if err != nil {
		log.Fatalf("Couldn't close file %s. Error: %v", filePath, err)
	}

	return file
}
