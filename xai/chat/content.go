// Package chat provides chat completion functionality for the xAI SDK.
package chat

// Part represents a part of a message content.
// Note: Content parts are not yet implemented in the proto definitions.
type Part interface {
	// Placeholder interface until content parts are properly defined in proto
	Content() string
}

// TextPart represents a text part of a message.
type TextPart struct {
	text string
}

// Content returns the text content.
func (t *TextPart) Content() string {
	return t.text
}

// Text creates a new text part.
func Text(text string) Part {
	return &TextPart{text: text}
}
