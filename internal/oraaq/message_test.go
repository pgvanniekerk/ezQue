package oraaq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaw(t *testing.T) {
	const content = "testContent"
	var id [16]byte

	message := &Message{
		ID:      id,
		Content: content,
	}

	raw := message.Raw()

	// check if the ID and Content properties of the raw message match those of the original message
	require.Equal(t, message.ID, raw.ID, "Raw ID does not match the original message's ID")
	require.Equal(t, message.Content, raw.Content, "Raw Content does not match the original message's Content")
}

func TestText(t *testing.T) {
	const content = "testContent"

	message := &Message{
		Content: content,
	}

	// Text() should return the Content property of the message
	require.Equal(t, content, message.Text(), "Text does not return the correct message Content")
}
