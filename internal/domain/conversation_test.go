package domain

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConversation_String(t *testing.T) {
	id := 1
	userID := 1
	text := "Test"
	c := Conversation{1, 1, text, time.Now()}
	dateString := c.Date.String()

	expected := fmt.Sprintf("Conversation ID: %d, userID: %d, text: %q, date: %s", id, userID, text, dateString)
	require.Equal(t, expected, c.String())
}

func TestConversation_IsValid(t *testing.T) {
	c1 := Conversation{1, 1, "Test", time.Now()}
	c2 := Conversation{1, 1, "", time.Now()}
	c3 := Conversation{1, 1, "Test", time.Now()}
	c3.Date = time.Unix(0xFFFFFFFFFFFFFF, 0)

	require.True(t, c1.IsValid())
	require.False(t, c2.IsValid())
	require.False(t, c3.IsValid())
}
