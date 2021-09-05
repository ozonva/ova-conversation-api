package api

import (
	"context"

	"github.com/rs/zerolog/log"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *serverGRPC) CreateConversationV1(ctx context.Context, req *conversationApi.CreateConversationV1Request) (*conversationApi.CreateConversationV1Response, error) {
	log.Info().Msg("CreateConversationV1")

	return &conversationApi.CreateConversationV1Response{}, nil
}
