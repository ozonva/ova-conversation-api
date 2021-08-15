package domain

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConversation_String(t *testing.T) {
	text := "Test"
	c := NewConversation(1, 1, text)
	dateString := c.Date.String()

	expected := fmt.Sprintf("Conversation text: %s, date: %s", text, dateString)
	require.Equal(t, expected, c.String())
}

func TestConversation_IsValid(t *testing.T) {
	c1 := NewConversation(1, 1, "Test")
	c2 := NewConversation(1, 1, "")
	c3 := NewConversation(1, 1, "Test")
	c3.Date = time.Unix(0xFFFFFFFFFFFFFF, 0)

	require.True(t, c1.IsValid())
	require.False(t, c2.IsValid())
	require.False(t, c3.IsValid())
}
