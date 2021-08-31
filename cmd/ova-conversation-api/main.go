package main

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"ova-conversation-api/internal/api"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

const (
	portGRPC = ":8082"
)

func runServerGRPC() {
	l, err := net.Listen("tcp", portGRPC)
	if err != nil {
		log.Fatal().Msgf("failed to listen TCP: %s", err)
	}
	log.Info().Msgf("ova-conversation-api: gRPC server started on the port %s", portGRPC)

	service := grpc.NewServer()
	conversationApi.RegisterConversationApiServer(service, api.NewConversationApiServer())
	if err = service.Serve(l); err != nil {
		log.Fatal().Msgf("failed to serve: %s", err)
	}
}

func main() {
	runServerGRPC()
}
