package saver

import (
	"log"
	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/flusher"
	"sync"
	"time"
)

type Saver interface {
	Init()
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
	return &saver{
		duration: duration,
		capacity: capacity,
		flusher:  flusher,
		data:     make([]domain.Conversation, 0, capacity),
	}
}

func (s *saver) Init() {
	ticker := time.NewTicker(time.Duration(s.duration) * time.Millisecond)
	s.chanTicker = make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

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

	wg.Wait()
}

func (s *saver) Save(entity domain.Conversation) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if (cap(s.data) == 0) || (len(s.data) == cap(s.data)-1) {
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
