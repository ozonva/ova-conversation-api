package saver

import (
	"log"
	"sync"
	"time"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/flusher"
)

type Saver interface {
	Save(entity domain.Conversation)
	Close()
}

type saver struct {
	mux        sync.Mutex
	chanTicker chan struct{}
	duration   uint
	capacity   uint
	flusher    flusher.Flusher
	data       []domain.Conversation
}

func NewSaver(duration uint, capacity uint, flusher flusher.Flusher) Saver {
	s := &saver{
		duration: duration,
		capacity: capacity,
		flusher:  flusher,
		data:     make([]domain.Conversation, 0, capacity),
	}
	s.init()

	return s
}

func (s *saver) init() {
	ticker := time.NewTicker(time.Duration(s.duration) * time.Millisecond)
	s.chanTicker = make(chan struct{})

	go func() {
		for {
			select {
			case _, ok := <-s.chanTicker:
				if !ok {
					ticker.Stop()

					return
				}
			case <-ticker.C:
				s.flush()
			}
		}
	}()
}

func (s *saver) Save(entity domain.Conversation) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if (s.capacity == 0) || (len(s.data) == int(s.capacity)-1) {
		s.data = append(s.data, entity)
		s.flush()
	}

	s.data = append(s.data, entity)
}

func (s *saver) Close() {
	s.flush()

	close(s.chanTicker)
}

func (s *saver) flush() {
	if len(s.data) == 0 {
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	notFlushed := s.flusher.Flush(s.data)
	if len(notFlushed) > 0 {
		s.data = make([]domain.Conversation, len(notFlushed), s.capacity)
		copy(s.data, notFlushed)
		log.Printf("saver did not flush the amount of data: %d %T", len(notFlushed), s.data)
	} else {
		s.data = make([]domain.Conversation, 0, s.capacity)
	}
}
