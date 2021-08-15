package domain

import (
	"fmt"
	"time"
)

type Conversation struct {
	ID     uint64
	UserID uint64
	Text   string
	Date   time.Time
}

func NewConversation(id uint64, userID uint64, text string) *Conversation {
	return &Conversation{
		ID:     id,
		UserID: userID,
		Text:   text,
		Date:   time.Now(),
	}
}

func (c *Conversation) String() string {
	return fmt.Sprintf("Conversation text: %s, date: %s", c.Text, c.Date.String())
}

func (c *Conversation) IsValid() bool {
	if len(c.Text) == 0 {
		return false
	}

	if c.Date.After(time.Now()) {
		return false
	}

	return true
}
