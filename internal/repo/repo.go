package repo

import "ova-conversation-api/internal/domain"

type Repo interface {
	AddEntities(entities []domain.Conversation) error
	ListEntities(limit, offset uint64) ([]domain.Conversation, error)
	DescribeEntity(entityId uint64) (*domain.Conversation, error)
}
