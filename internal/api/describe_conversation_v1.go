package api

import (
	"context"

	"github.com/rs/zerolog/log"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *serverGRPC) DescribeConversationV1(ctx context.Context, req *conversationApi.DescribeConversationV1Request) (*conversationApi.DescribeConversationV1Response, error) {
	log.Info().Msg("DescribeConversationV1")

	return &conversationApi.DescribeConversationV1Response{}, nil
}
