// Package embed provides client functionality for the xAI Embeddings API.
package embed

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Client provides access to the xAI Embeddings API.
type Client struct {
	// client is the underlying gRPC client (to be added when service is defined)
}

// NewClient creates a new embeddings client.
func NewClient() *Client {
	return &Client{}
}

// Request represents an embeddings request.
type Request struct {
	proto *xaiv1.EmbedRequest
}

// NewRequest creates a new embeddings request.
func NewRequest(model string, inputs ...Input) *Request {
	embedInputs := make([]*xaiv1.EmbedInput, 0, len(inputs))
	for _, input := range inputs {
		embedInputs = append(embedInputs, input.Proto())
	}

	return &Request{
		proto: &xaiv1.EmbedRequest{
			Model: model,
			Input: embedInputs,
		},
	}
}

// WithEncodingFormat sets the encoding format for the embeddings.
func (r *Request) WithEncodingFormat(format xaiv1.EmbedEncodingFormat) *Request {
	r.proto.EncodingFormat = format
	return r
}

// WithUser sets the user identifier for the request.
func (r *Request) WithUser(user string) *Request {
	r.proto.User = user
	return r
}

// Proto returns the underlying protobuf request.
func (r *Request) Proto() *xaiv1.EmbedRequest {
	return r.proto
}

// Input represents an embedding input (text or image).
type Input struct {
	proto *xaiv1.EmbedInput
}

// Text creates a text input for embedding.
func Text(text string) Input {
	return Input{
		proto: &xaiv1.EmbedInput{
			String_: text,
		},
	}
}

// Image creates an image input for embedding.
func Image(imageURL string, detail xaiv1.ImageDetail) Input {
	return Input{
		proto: &xaiv1.EmbedInput{
			ImageUrl: &xaiv1.ImageUrlContent{
				ImageUrl: imageURL,
				Detail:   detail,
			},
		},
	}
}

// Proto returns the underlying protobuf input.
func (i Input) Proto() *xaiv1.EmbedInput {
	return i.proto
}

// Response represents an embeddings response.
type Response struct {
	proto *xaiv1.EmbedResponse
}

// ID returns the response ID.
func (r *Response) ID() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.Id
}

// Model returns the model used.
func (r *Response) Model() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.Model
}

// SystemFingerprint returns the system fingerprint.
func (r *Response) SystemFingerprint() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.SystemFingerprint
}

// Embeddings returns all embeddings in the response.
func (r *Response) Embeddings() []Embedding {
	if r.proto == nil || len(r.proto.Embeddings) == 0 {
		return nil
	}

	embeddings := make([]Embedding, len(r.proto.Embeddings))
	for i, e := range r.proto.Embeddings {
		embeddings[i] = Embedding{proto: e}
	}
	return embeddings
}

// Usage returns the usage information.
func (r *Response) Usage() *xaiv1.EmbeddingUsage {
	if r.proto == nil {
		return nil
	}
	return r.proto.Usage
}

// Proto returns the underlying protobuf response.
func (r *Response) Proto() *xaiv1.EmbedResponse {
	return r.proto
}

// Embedding represents a single embedding result.
type Embedding struct {
	proto *xaiv1.Embedding
}

// Index returns the index of this embedding.
func (e *Embedding) Index() int32 {
	if e.proto == nil {
		return 0
	}
	return e.proto.Index
}

// Vectors returns the feature vectors for this embedding.
func (e *Embedding) Vectors() []FeatureVector {
	if e.proto == nil || len(e.proto.Embeddings) == 0 {
		return nil
	}

	vectors := make([]FeatureVector, len(e.proto.Embeddings))
	for i, v := range e.proto.Embeddings {
		vectors[i] = FeatureVector{proto: v}
	}
	return vectors
}

// Proto returns the underlying protobuf embedding.
func (e *Embedding) Proto() *xaiv1.Embedding {
	return e.proto
}

// FeatureVector represents an embedding vector.
type FeatureVector struct {
	proto *xaiv1.FeatureVector
}

// FloatArray returns the embedding as a float array.
func (f *FeatureVector) FloatArray() []float32 {
	if f.proto == nil {
		return nil
	}
	return f.proto.FloatArray
}

// Base64Array returns the embedding as a base64-encoded string.
func (f *FeatureVector) Base64Array() string {
	if f.proto == nil {
		return ""
	}
	return f.proto.Base64Array
}

// Proto returns the underlying protobuf feature vector.
func (f *FeatureVector) Proto() *xaiv1.FeatureVector {
	return f.proto
}

// Generate generates embeddings for the given request.
// Note: This is a placeholder. Actual implementation requires gRPC client setup.
func (c *Client) Generate(ctx context.Context, req *Request) (*Response, error) {
	// TODO: Implement actual gRPC call when service is defined
	return nil, nil
}
