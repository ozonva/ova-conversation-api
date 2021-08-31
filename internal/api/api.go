package api

import conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"

type serverGRPC struct {
	conversationApi.ConversationApiServer
}

func NewConversationApiServer() conversationApi.ConversationApiServer {
	return &serverGRPC{}
}
