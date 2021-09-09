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
	const nameHandler = "CreateConversationV1"

	reqCreate := req.GetCreateConversation()

	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, nameHandler)
	parentSpan.LogFields(
		openTrLog.Int("Total conversations count", len(reqCreate)),
		openTrLog.Int("Num conversations in batch", batchSize))

	defer parentSpan.Finish()

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
		err = s.repo.AddEntities(batches[i])
		if err != nil {
			log.Error().Err(err).Msg("create conversations")
			return nil, status.Error(codes.Internal, "internal error")
		}

		childSpan := opentracing.StartSpan("multi create conversions by batches", opentracing.ChildOf(parentSpan.Context()))
		childSpan.LogFields(openTrLog.Int("size of batch", binary.Size(batches[i])))
		childSpan.Finish()
	}

	return &emptypb.Empty{}, nil
}
