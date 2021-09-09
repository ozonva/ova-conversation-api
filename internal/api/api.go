package api

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/kafka"
	"ova-conversation-api/internal/repo"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

type apiServer struct {
	conversationApi.ConversationApiServer
	repo          repo.Repo
	kafkaProducer kafka.Producer
}

func NewConversationApiServer(r repo.Repo, kp kafka.Producer) conversationApi.ConversationApiServer {
	return &apiServer{repo: r, kafkaProducer: kp}
}

func checkValidateError(handler string, err error) error {
	if err != nil {
		log.Error().Msgf("%s, error: %s", handler, err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func entityToResponse(e *domain.Conversation) *conversationApi.DescribeConversationV1Response {
	response := conversationApi.DescribeConversationV1Response{
		Id:     e.ID,
		UserId: e.UserID,
		Text:   e.Text,
		Date:   timestamppb.New(e.Date),
	}

	return &response
}
