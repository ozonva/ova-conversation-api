package api

import (
	"context"
	"encoding/binary"

	"github.com/opentracing/opentracing-go"
	openTrLog "github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/utils"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

const (
	batchSize = 2
)

func (s *apiServer) MultiCreateConversationsV1(ctx context.Context, req *conversationApi.MultiCreateConversationsV1Request) (*emptypb.Empty, error) {
	nameHandler := "CreateConversationV1"

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan(nameHandler)

	defer span.Finish()

	log.Info().Msg(nameHandler)
	if req == nil {
		log.Info().Msg("MultiCreateConversationsV1Request is null")
		return nil, status.Error(codes.InvalidArgument, "MultiCreateConversationsV1Request is null")
	}
	log.Info().Msgf("request: %s", req.String())

	err := checkValidateError("create conversation", req.Validate())
	if err != nil {
		return nil, err
	}

	reqCreate := req.GetCreateConversation()
	entities := make([]domain.Conversation, 0, len(reqCreate))
	for i := range reqCreate {
		entities = append(entities, domain.Conversation{UserID: reqCreate[i].GetUserId(), Text: reqCreate[i].GetText()})
	}

	batches, err := utils.MakeSliceOfSlicesConversation(entities, batchSize)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}

	for i := range batches {
		childSpan := tracer.StartSpan("multi create conversions by batches", opentracing.ChildOf(span.Context()))
		childSpan.LogFields(openTrLog.Int("size", binary.Size(batches[i])))

		err = s.repo.AddEntities(batches[i])
		if err != nil {
			childSpan.Finish()

			log.Error().Msgf("repo: create conversation: %s", err)
			return nil, status.Error(codes.Internal, "internal error")
		}

		childSpan.Finish()
	}

	return &emptypb.Empty{}, nil
}
