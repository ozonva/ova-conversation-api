package api

import (
	"context"
	"database/sql"
	"net"
	"sync"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/kafka"
	"ova-conversation-api/internal/repo"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

var _ = Describe("API", func() {

	const bufSize = 2048 * 2048

	var (
		bufListener   *bufconn.Listener
		mockCtrl      *gomock.Controller
		mockRepo      *repo.MockRepo
		mockKafkaProd *kafka.MockProducer
		grpcServer    *grpc.Server
		ctx           context.Context
		startWG       sync.WaitGroup
		connect       *grpc.ClientConn
		client        conversationApi.ConversationApiClient
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = repo.NewMockRepo(mockCtrl)
		mockKafkaProd = kafka.NewMockProducer(mockCtrl)
		grpcServer = grpc.NewServer()
		conversationApi.RegisterConversationApiServer(grpcServer, NewConversationApiServer(mockRepo, mockKafkaProd))
		bufListener = bufconn.Listen(bufSize)
		startWG.Add(1)
		go func() {
			startWG.Done()
			Expect(grpcServer.Serve(bufListener)).To(BeNil())
		}()
		startWG.Wait()
		ctx = context.Background()
		conn, err := grpc.DialContext(
			ctx, "bufnet",
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return bufListener.Dial()
			}), grpc.WithInsecure())
		Expect(err).To(BeNil())
		connect = conn
		client = conversationApi.NewConversationApiClient(conn)
	})
	AfterEach(func() {
		mockCtrl.Finish()
		connect.Close()
		grpcServer.Stop()
	})

	Describe("Create new conversation", func() {
		When("Can create", func() {
			expectedConversation := domain.Conversation{UserID: 1, Text: "test"}
			expectedId := uint64(1)
			BeforeEach(func() {
				mockRepo.EXPECT().AddEntity(expectedConversation).Return(expectedId, nil).Times(1)
				mockKafkaProd.EXPECT().Send(kafka.Message{
					Type: kafka.Create,
					Body: map[string]interface{}{"id": expectedId}}).
					Times(1)
			})
			It("should return empty response", func() {
				req := conversationApi.CreateConversationV1Request{
					UserId: expectedConversation.UserID,
					Text:   expectedConversation.Text,
				}
				response, err := client.CreateConversationV1(ctx, &req)
				Expect(err).To(BeNil())
				Expect(response.String()).To(Equal(""))
			})
		})
		When("Can`t create because internal problem", func() {
			expectedConversation := domain.Conversation{UserID: 1, Text: "test"}
			expectedId := uint64(1)
			BeforeEach(func() {
				mockRepo.EXPECT().AddEntity(expectedConversation).Return(uint64(0), sql.ErrNoRows).Times(1)
				mockKafkaProd.EXPECT().Send(kafka.Message{
					Type: kafka.Create,
					Body: map[string]interface{}{"id": expectedId}}).
					Times(0)
			})
			It("internal error", func() {
				req := conversationApi.CreateConversationV1Request{
					UserId: expectedConversation.UserID,
					Text:   expectedConversation.Text,
				}
				response, err := client.CreateConversationV1(ctx, &req)
				Expect(err).To(Equal(status.Error(codes.Internal, "internal error")))
				Expect(response).To(BeNil())
			})
		})
	})

	Describe("Multi create conversations", func() {
		When("Can create", func() {
			expectedConversations := []domain.Conversation{
				{UserID: 1, Text: "test1"},
				{UserID: 2, Text: "test2"},
				{UserID: 3, Text: "test3"},
			}
			BeforeEach(func() {
				gomock.InOrder(
					mockRepo.EXPECT().AddEntities(expectedConversations[:2]).Return(nil).Times(1),
					mockRepo.EXPECT().AddEntities(expectedConversations[2:]).Return(nil).Times(1),
				)
			})
			It("should return empty response", func() {
				creates := make([]*conversationApi.CreateConversationV1Request, 0, len(expectedConversations))
				for _, expectedConversation := range expectedConversations {
					creates = append(creates, &conversationApi.CreateConversationV1Request{
						UserId: expectedConversation.UserID,
						Text:   expectedConversation.Text,
					})
				}

				req := conversationApi.MultiCreateConversationsV1Request{CreateConversation: creates}

				response, err := client.MultiCreateConversationsV1(ctx, &req)
				Expect(response.String()).To(Equal(""))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Describe conversation", func() {
		When("Can get", func() {
			expectedConversation := domain.Conversation{ID: 1, UserID: 17, Text: "test"}
			expectedId := expectedConversation.ID
			BeforeEach(func() {
				mockRepo.EXPECT().DescribeEntity(expectedId).Return(&expectedConversation, nil).Times(1)
			})
			It("should return describe", func() {
				req := conversationApi.DescribeConversationV1Request{
					Id: expectedId,
				}
				response, err := client.DescribeConversationV1(ctx, &req)

				Expect(err).To(BeNil())
				Expect(response.Id).To(Equal(expectedConversation.ID))
				Expect(response.UserId).To(Equal(expectedConversation.UserID))
				Expect(response.Text).To(Equal(expectedConversation.Text))
			})
		})

		When("Can`t get", func() {
			expectedConversation := domain.Conversation{ID: 1, UserID: 17, Text: "test"}
			expectedId := expectedConversation.ID
			BeforeEach(func() {
				mockRepo.EXPECT().DescribeEntity(expectedId).Return(nil, sql.ErrNoRows).Times(1)
			})
			It("no return describe", func() {
				req := conversationApi.DescribeConversationV1Request{
					Id: expectedId,
				}
				response, err := client.DescribeConversationV1(ctx, &req)

				Expect(err).To(Equal(status.Error(codes.Internal, "internal error")))
				Expect(response).To(BeNil())
			})
		})
	})
})
