// Package tokenizer provides a client for the xAI Tokenization API.
package tokenizer

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Tokenization API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Tokenization API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// Token represents a single token.
type Token struct {
	TokenID     uint32
	StringToken string
	TokenBytes  []byte
}

// Response represents the tokenization response.
type Response struct {
	Tokens []*Token
	Model  string
}

// Tokenize tokenizes the given text.
func (c *Client) Tokenize(ctx context.Context, text, model, user string) (*Response, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.TokenizeTextRequest{
		Text:  text,
		Model: model,
		User:  user,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/tokenize", jsonData)
	if err != nil {
		return nil, err
	}

	var tokenResp xaiv1.TokenizeTextResponse
	if err := protojson.Unmarshal(resp.Body, &tokenResp); err != nil {
		return nil, err
	}

	tokens := make([]*Token, len(tokenResp.Tokens))
	for i, t := range tokenResp.Tokens {
		tokens[i] = &Token{
			TokenID:     t.TokenId,
			StringToken: t.StringToken,
			TokenBytes:  t.TokenBytes,
		}
	}

	return &Response{
		Tokens: tokens,
		Model:  tokenResp.Model,
	}, nil
}
