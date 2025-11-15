// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	"strings"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Message represents a chat message with role and content.
type Message struct {
	proto *xaiv1.Message
	parts []Part
}

// NewMessage creates a new message with the given role and content parts.
func NewMessage(role string, parts ...Part) *Message {
	var content strings.Builder
	for _, p := range parts {
		content.WriteString(p.Content())
	}
	return &Message{
		proto: &xaiv1.Message{
			Role:    role,
			Content: content.String(),
		},
		parts: parts,
	}
}

// System creates a system message.
func System(parts ...Part) *Message {
	return NewMessage("system", parts...)
}

// User creates a user message.
func User(parts ...Part) *Message {
	return NewMessage("user", parts...)
}

// Assistant creates an assistant message.
func Assistant(parts ...Part) *Message {
	return NewMessage("assistant", parts...)
}

// Proto returns the underlying protobuf message.
func (m *Message) Proto() *xaiv1.Message {
	return m.proto
}

// Role returns the role of the message.
func (m *Message) Role() string {
	return m.proto.GetRole()
}

// Content returns the content of the message.
func (m *Message) Content() string {
	return m.proto.GetContent()
}

// WithRole sets the role of the message.
func (m *Message) WithRole(role string) *Message {
	m.proto.Role = role
	return m
}

// Parts returns the parts of the message.
func (m *Message) Parts() []Part {
	return m.parts
}