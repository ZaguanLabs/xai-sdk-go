// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Message represents a chat message with role and content.
type Message struct {
	proto *xaiv1.Message
}

// NewMessage creates a new message with the given role and content.
func NewMessage(role, content string) Message {
	return Message{
		proto: &xaiv1.Message{
			Role:    role,
			Content: content,
		},
	}
}

// System creates a system message.
func System(content string) Message {
	return NewMessage("system", content)
}

// User creates a user message.
func User(content string) Message {
	return NewMessage("user", content)
}

// Assistant creates an assistant message.
func Assistant(content string) Message {
	return NewMessage("assistant", content)
}

// Proto returns the underlying protobuf message.
func (m Message) Proto() *xaiv1.Message {
	return m.proto
}

// Role returns the role of the message.
func (m Message) Role() string {
	return m.proto.GetRole()
}

// Content returns the content of the message.
func (m Message) Content() string {
	return m.proto.GetContent()
}

// WithRole sets the role of the message.
func (m Message) WithRole(role string) Message {
	m.proto.Role = role
	return m
}

// WithContent sets the content of the message.
func (m Message) WithContent(content string) Message {
	m.proto.Content = content
	return m
}