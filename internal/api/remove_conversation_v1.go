package api

import (
	"context"

	"github.com/rs/zerolog/log"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *serverGRPC) RemoveConversationV1(ctx context.Context, req *conversationApi.RemoveConversationV1Request) (*conversationApi.RemoveConversationV1Response, error) {
	log.Info().Msg("RemoveConversationV1")

	return &conversationApi.RemoveConversationV1Response{}, nil
}
