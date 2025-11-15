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
	// Convert parts to Content array
	contents := make([]*xaiv1.Content, 0, len(parts))
	for _, p := range parts {
		contents = append(contents, &xaiv1.Content{
			Text: p.Content(),
		})
	}

	return &Message{
		proto: &xaiv1.Message{
			Role:    roleToProto(role),
			Content: contents,
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
	return roleFromProto(m.proto.GetRole())
}

// Content returns the content of the message as a single string.
func (m *Message) Content() string {
	if len(m.proto.Content) == 0 {
		return ""
	}
	// Concatenate all text content
	var result strings.Builder
	for _, c := range m.proto.Content {
		result.WriteString(c.Text)
	}
	return result.String()
}

// WithRole sets the role of the message.
func (m *Message) WithRole(role string) *Message {
	m.proto.Role = roleToProto(role)
	return m
}

// Parts returns the parts of the message.
func (m *Message) Parts() []Part {
	return m.parts
}
