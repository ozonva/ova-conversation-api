package api

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"ova-conversation-api/internal/kafka"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func (s *apiServer) RemoveConversationV1(ctx context.Context, req *conversationApi.RemoveConversationV1Request) (*emptypb.Empty, error) {
	const nameHandler = "RemoveConversationV1"

	span, ctx := opentracing.StartSpanFromContext(ctx, nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("RemoveConversationV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "RemoveConversationV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("remove conversation", req.Validate())
	if err != nil {
		return nil, err
	}

	id, err := s.repo.RemoveEntity(req.GetId())
	if err != nil {
		log.Error().Err(err).Msg("remove conversation")
		return nil, status.Error(codes.Internal, "internal error")
	}

	msg := kafka.Message{
		Type: kafka.Remove,
		Body: map[string]interface{}{
			"id": id,
		},
	}

	err = s.kafkaProducer.Send(msg)
	if err != nil {
		return nil, err
	}

	promRemoveCntr.Inc()

	return &emptypb.Empty{}, nil
}
