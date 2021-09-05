package flusher

import (
	"log"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/repo"
	"ova-conversation-api/internal/utils"
)

type Flusher interface {
	Flush(entities []domain.Conversation) []domain.Conversation
}

func NewFlusher(chunkSize int, entityRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize:  chunkSize,
		entityRepo: entityRepo,
	}
}

type flusher struct {
	chunkSize  int
	entityRepo repo.Repo
}

func (f flusher) Flush(entities []domain.Conversation) []domain.Conversation {
	chunks, err := utils.MakeSliceOfSlicesConversation(entities, f.chunkSize)
	if err != nil {
		log.Println(err)

		return entities
	}

	notFlushed := make([]domain.Conversation, 0, len(entities))
	for _, chunk := range chunks {
		err := f.entityRepo.AddEntities(chunk)
		if err != nil {
			notFlushed = append(notFlushed, chunk...)
		}
	}

	if len(notFlushed) == 0 {
		return nil
	}

	return notFlushed
}
