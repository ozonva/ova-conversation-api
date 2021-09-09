package api

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/kafka"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *apiServer) UpdateConversationV1(ctx context.Context, req *conversationApi.UpdateConversationV1Request) (*emptypb.Empty, error) {
	const nameHandler = "UpdateConversationV1"

	span, ctx := opentracing.StartSpanFromContext(ctx, nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("UpdateConversationV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "UpdateConversationV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("update conversation", req.Validate())
	if err != nil {
		return nil, err
	}

	id, err := s.repo.UpdateEntity(domain.Conversation{
		ID:   req.GetId(),
		Text: req.GetText(),
	})
	if err != nil {
		log.Error().Err(err).Msg("update conversation")
		return nil, status.Error(codes.Internal, "internal error")
	}
	if id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	msg := kafka.Message{
		Type: kafka.Update,
		Body: map[string]interface{}{
			"id": id,
		},
	}

	err = s.kafkaProducer.Send(msg)
	if err != nil {
		return nil, err
	}

	promUpdateCntr.Inc()

	return &emptypb.Empty{}, nil
}
