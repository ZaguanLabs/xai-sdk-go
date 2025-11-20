// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	"encoding/json"
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
		switch p.Type() {
		case PartTypeImage:
			// Handle image parts
			if img, ok := p.(*ImagePart); ok {
				contents = append(contents, &xaiv1.Content{
					ImageUrl: &xaiv1.ImageUrlContent{
						ImageUrl: img.ImageURL(),
						Detail:   img.Detail(),
					},
				})
			}
		case PartTypeFile:
			// Handle file parts
			if file, ok := p.(*FilePart); ok {
				contents = append(contents, &xaiv1.Content{
					File: &xaiv1.FileContent{
						FileId: file.FileID(),
					},
				})
			}
		default:
			// Handle text parts (default)
			contents = append(contents, &xaiv1.Content{
				Text: p.Content(),
			})
		}
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

// Name returns the name of the message sender.
// This is useful for multi-user conversations to identify participants.
func (m *Message) Name() string {
	if m.proto == nil {
		return ""
	}
	return m.proto.Name
}

// WithName sets the name of the message sender.
func (m *Message) WithName(name string) *Message {
	if m.proto != nil {
		m.proto.Name = name
	}
	return m
}

// ToolCalls returns the tool calls in the message.
func (m *Message) ToolCalls() []*ToolCall {
	if m.proto == nil || len(m.proto.ToolCalls) == 0 {
		return nil
	}

	result := make([]*ToolCall, 0, len(m.proto.ToolCalls))
	for _, protoCall := range m.proto.ToolCalls {
		toolCall := parseToolCall(protoCall)
		if toolCall != nil {
			result = append(result, toolCall)
		}
	}
	return result
}

// ReasoningContent returns the reasoning content of the message.
func (m *Message) ReasoningContent() string {
	if m.proto == nil {
		return ""
	}
	return m.proto.ReasoningContent
}

// EncryptedContent returns the encrypted content of the message.
func (m *Message) EncryptedContent() string {
	if m.proto == nil {
		return ""
	}
	return m.proto.EncryptedContent
}

// WithToolCalls sets the tool calls for the message.
func (m *Message) WithToolCalls(toolCalls []*ToolCall) *Message {
	if m.proto == nil {
		return m
	}
	m.proto.ToolCalls = make([]*xaiv1.ToolCall, 0, len(toolCalls))
	for _, tc := range toolCalls {
		if tc == nil {
			continue
		}
		// Convert to proto ToolCall
		argsJSON, _ := json.Marshal(tc.Arguments())
		m.proto.ToolCalls = append(m.proto.ToolCalls, &xaiv1.ToolCall{
			Id:           tc.ID(),
			Type:         xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL,
			Status:       parseToolCallStatus(tc.Status()),
			ErrorMessage: tc.ErrorMessage(),
			Function: &xaiv1.FunctionCall{
				Name:      tc.Name(),
				Arguments: string(argsJSON),
			},
		})
	}
	return m
}

// parseToolCallStatus converts a status string to ToolCallStatus enum.
func parseToolCallStatus(status string) xaiv1.ToolCallStatus {
	switch status {
	case "TOOL_CALL_STATUS_IN_PROGRESS":
		return xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS
	case "TOOL_CALL_STATUS_COMPLETED":
		return xaiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED
	case "TOOL_CALL_STATUS_INCOMPLETE":
		return xaiv1.ToolCallStatus_TOOL_CALL_STATUS_INCOMPLETE
	case "TOOL_CALL_STATUS_FAILED":
		return xaiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED
	default:
		return xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS
	}
}

// WithReasoningContent sets the reasoning content for the message.
func (m *Message) WithReasoningContent(reasoning string) *Message {
	if m.proto != nil {
		m.proto.ReasoningContent = reasoning
	}
	return m
}

// WithEncryptedContent sets the encrypted content for the message.
func (m *Message) WithEncryptedContent(encrypted string) *Message {
	if m.proto != nil {
		m.proto.EncryptedContent = encrypted
	}
	return m
}
