// Package chat provides chat completion functionality for the xAI SDK.
package chat

import xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"

// PartType represents the type of content part.
type PartType int

const (
	// PartTypeText represents text content.
	PartTypeText PartType = iota
	// PartTypeImage represents image content.
	PartTypeImage
	// PartTypeFile represents file content.
	PartTypeFile
)

// Part represents a part of a message content.
// This interface allows for different content types (text, images, files, etc.)
// to be used in messages.
type Part interface {
	// Content returns the string representation of the part.
	Content() string
	// Type returns the type of this content part.
	Type() PartType
}

// TextPart represents a text part of a message.
type TextPart struct {
	text string
}

// Content returns the text content.
func (t *TextPart) Content() string {
	return t.text
}

// Type returns PartTypeText.
func (t *TextPart) Type() PartType {
	return PartTypeText
}

// Text creates a new text part.
func Text(text string) Part {
	return &TextPart{text: text}
}

// ImagePart represents an image part of a message.
type ImagePart struct {
	url    string
	detail xaiv1.ImageDetail
}

// Content returns the image URL.
func (i *ImagePart) Content() string {
	return i.url
}

// Type returns PartTypeImage.
func (i *ImagePart) Type() PartType {
	return PartTypeImage
}

// ImageURL returns the image URL.
func (i *ImagePart) ImageURL() string {
	return i.url
}

// Detail returns the image detail level.
func (i *ImagePart) Detail() xaiv1.ImageDetail {
	return i.detail
}

// ImageDetail represents the level of detail for image processing.
type ImageDetail string

const (
	// ImageDetailAuto lets the system select an appropriate resolution (default).
	ImageDetailAuto ImageDetail = "auto"
	// ImageDetailLow uses low-resolution image, reducing token usage and increasing speed.
	ImageDetailLow ImageDetail = "low"
	// ImageDetailHigh uses high-resolution image, increasing token usage but capturing more detail.
	ImageDetailHigh ImageDetail = "high"
)

// Image creates an image content part.
// The imageURL can be:
// - A URL to an image (PNG or JPG)
// - A base64-encoded data URI (e.g., "data:image/jpeg;base64,...")
//
// The detail parameter controls image resolution:
// - "auto": System selects appropriate resolution (default)
// - "low": Low resolution, faster, uses fewer tokens
// - "high": High resolution, slower, uses more tokens but captures more detail
//
// Image requirements:
// - Formats: PNG or JPG only
// - Size: Maximum 10 MiB
// - Fetch timeout: 5 seconds (for URLs)
// - User agent: "XaiImageApiFetch/1.0" (for URLs)
func Image(imageURL string, detail ...ImageDetail) Part {
	d := ImageDetailAuto
	if len(detail) > 0 {
		d = detail[0]
	}

	var pbDetail xaiv1.ImageDetail
	switch d {
	case ImageDetailLow:
		pbDetail = xaiv1.ImageDetail_DETAIL_LOW
	case ImageDetailHigh:
		pbDetail = xaiv1.ImageDetail_DETAIL_HIGH
	default:
		pbDetail = xaiv1.ImageDetail_DETAIL_AUTO
	}

	return &ImagePart{
		url:    imageURL,
		detail: pbDetail,
	}
}

// FilePart represents a file part of a message.
type FilePart struct {
	fileID string
}

// Content returns the file ID.
func (f *FilePart) Content() string {
	return f.fileID
}

// Type returns PartTypeFile.
func (f *FilePart) Type() PartType {
	return PartTypeFile
}

// FileID returns the file ID.
func (f *FilePart) FileID() string {
	return f.fileID
}

// File creates a file content part.
// The fileID should be the ID of a previously uploaded file.
// You can obtain this ID by uploading a file using the Files API.
func File(fileID string) Part {
	return &FilePart{fileID: fileID}
}
