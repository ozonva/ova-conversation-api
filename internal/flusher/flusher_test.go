package flusher

import (
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ova-conversation-api/internal/domain"
	"ova-conversation-api/internal/repo"
)

var _ = Describe("Flusher", func() {
	var (
		mockCtrl  *gomock.Controller
		mockRepo  *repo.MockRepo
		f         Flusher
		entities  = listConversations()
		chunkSize = 5
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = repo.NewMockRepo(mockCtrl)
		f = NewFlusher(chunkSize, mockRepo)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Saving data", func() {
		When("can save all the data", func() {
			Context("entities count less than chunkSize", func() {
				BeforeEach(func() {
					mockRepo.EXPECT().AddEntities(entities[:4]).Return(nil).Times(1)
				})
				It("Flush should return nil", func() {
					Expect(f.Flush(entities[:4])).To(BeNil())
				})
			})
			Context("entities count more than chunkSize", func() {
				BeforeEach(func() {
					gomock.InOrder(
						mockRepo.EXPECT().AddEntities(entities[:5]).Return(nil).Times(1),
						mockRepo.EXPECT().AddEntities(entities[5:]).Return(nil).Times(1),
					)
				})
				It("Flush should return nil", func() {
					Expect(f.Flush(entities)).To(BeNil())
				})
			})
		})
		When("unable to save", func() {
			err := errors.New("error saving data")
			Context("all data", func() {
				BeforeEach(func() {
					gomock.InOrder(
						mockRepo.EXPECT().AddEntities(entities[:5]).Return(err).Times(1),
						mockRepo.EXPECT().AddEntities(entities[5:]).Return(err).Times(1),
					)
				})
				It("Flush should return all data", func() {
					Expect(f.Flush(entities)).To(Equal(entities))
				})
			})
			Context("part of data", func() {
				BeforeEach(func() {
					gomock.InOrder(
						mockRepo.EXPECT().AddEntities(entities[:5]).Return(err).Times(1),
						mockRepo.EXPECT().AddEntities(entities[5:]).Return(nil).Times(1),
					)
				})
				It("Flush should return only unsaved data", func() {
					Expect(f.Flush(entities)).To(Equal(entities[:5]))
				})
			})
		})
	})
})

func listConversations() []domain.Conversation {
	return []domain.Conversation{
		{ID: 1, UserID: 1, Text: "c1", Date: time.Now()},
		{ID: 2, UserID: 1, Text: "c2", Date: time.Now()},
		{ID: 3, UserID: 1, Text: "c3", Date: time.Now()},
		{ID: 4, UserID: 2, Text: "c4", Date: time.Now()},
		{ID: 5, UserID: 2, Text: "c5", Date: time.Now()},
		{ID: 6, UserID: 2, Text: "c6", Date: time.Now()},
	}
}
