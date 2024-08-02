package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kanataidarov/teams_automator/config"
	mc "github.com/kanataidarov/teams_automator/internal/msgraph_client"
	"github.com/kanataidarov/teams_automator/internal/openai_client"
	pb "github.com/kanataidarov/teams_automator/pkg/grpc"
	mapi "github.com/kanataidarov/teams_automator/pkg/model/rest/msgraph_api"
	"github.com/kanataidarov/teams_automator/pkg/types"
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

	doTeams(ctx, cfg, transcription)

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

func doTeams(ctx context.Context, cfg *config.Config, msgContent string) {
	user, err := mc.Get[mapi.User](ctx, cfg, mc.ProfileUrl())
	if err != nil {
		log.Printf("Error getting user profile. Error: %v", err)
	}
	userId := user.Id
	log.Println("User id", userId)

	chats, err := mc.Get[mapi.ChatsResponse](ctx, cfg, mc.ChatsUrl(userId))
	if err != nil {
		log.Printf("Error getting chats. Error: %v", err)
	}

	directChat := mapi.Chat{}
	for _, chatItem := range chats.Value {
		if chatItem.ChatType == types.OneOnOne && chatItem.LastMsg.LastMsgFrom.User.Id == userId {
			log.Printf("My Last Direct Msg Id: " + chatItem.LastMsg.Id)
			directChat = chatItem
			break
		}
	}

	msgRequest := mapi.MsgRequest{MsgBody: mapi.MsgBody{Content: msgContent, ContentType: "text"}}
	msgResponse, err := mc.Post[mapi.MsgResponse](ctx, cfg, mc.PostMsgUrl(userId, directChat.Id), msgRequest)
	if err != nil {
		log.Printf("Error posting msg. Error: %v", err)
	}
	log.Println("MsgResponse", msgResponse)
}
