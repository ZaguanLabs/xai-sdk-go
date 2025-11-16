// Package chat provides chat completion functionality for the xAI SDK.
package chat

// Part represents a part of a message content.
// This interface allows for different content types (text, images, files, etc.)
// to be used in messages.
type Part interface {
	// Content returns the string representation of the part.
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
