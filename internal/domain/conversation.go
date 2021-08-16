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

func (c *Conversation) String() string {
	return fmt.Sprintf("Conversation ID: %d, userID: %d, text: %q, date: %s", c.ID, c.UserID, c.Text, c.Date.String())
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
