package api

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *apiServer) ListConversationsV1(ctx context.Context, req *conversationApi.ListConversationsV1Request) (*conversationApi.ListConversationsV1Response, error) {
	const nameHandler = "ListConversationsV1"

	span, ctx := opentracing.StartSpanFromContext(ctx, nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("ListConversationsV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "ListConversationsV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("list of conversations", req.Validate())
	if err != nil {
		return nil, err
	}

	entities, err := s.repo.ListEntities(req.GetLimit(), req.GetOffset())
	if err != nil {
		log.Error().Err(err).Msg("list of conversations")
		return nil, status.Error(codes.Internal, "internal error")
	}

	result := make([]*conversationApi.DescribeConversationV1Response, 0, len(entities))
	for i := range entities {
		result = append(result, entityToResponse(&entities[i]))
	}

	return &conversationApi.ListConversationsV1Response{Conversations: result}, nil
}
