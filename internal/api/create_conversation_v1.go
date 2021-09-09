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

func (s *apiServer) CreateConversationV1(ctx context.Context, req *conversationApi.CreateConversationV1Request) (*emptypb.Empty, error) {
	const nameHandler = "CreateConversationV1"

	span, ctx := opentracing.StartSpanFromContext(ctx, nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("CreateConversationV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "CreateConversationV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("create conversation", req.Validate())
	if err != nil {
		return nil, err
	}

	id, err := s.repo.AddEntity(domain.Conversation{
		UserID: req.GetUserId(),
		Text:   req.GetText(),
	})
	if err != nil {
		log.Error().Err(err).Msg("create conversation")
		return nil, status.Error(codes.Internal, "internal error")
	}

	msg := kafka.Message{
		Type: kafka.Create,
		Body: map[string]interface{}{
			"id": id,
		},
	}

	err = s.kafkaProducer.Send(msg)
	if err != nil {
		return nil, err
	}

	promCreateCntr.Inc()

	return &emptypb.Empty{}, nil
}
