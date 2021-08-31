package api

import (
	"context"

	"github.com/rs/zerolog/log"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *serverGRPC) ListConversationsV1(ctx context.Context, req *conversationApi.ListConversationsV1Request) (*conversationApi.ListConversationsV1Response, error) {
	log.Info().Msg("ListConversationsV1")

	return &conversationApi.ListConversationsV1Response{}, nil
}
