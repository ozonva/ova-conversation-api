package api

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *apiServer) DescribeConversationV1(ctx context.Context, req *conversationApi.DescribeConversationV1Request) (*conversationApi.DescribeConversationV1Response, error) {
	nameHandler := "DescribeConversationV1"

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan(nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("DescribeConversationV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "DescribeConversationV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("describe conversation", req.Validate())
	if err != nil {
		return nil, err
	}

	entity, err := s.repo.DescribeEntity(req.GetId())
	if err != nil {
		log.Error().Msgf("repo: describe conversation: %s", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if entity == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return entityToResponse(entity), nil
}
