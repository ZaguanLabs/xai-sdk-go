// Package sample provides a client for the xAI Sample/Completion API (legacy).
package sample

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Sample/Completion API.
// Note: This is a legacy API. The Chat API is recommended for new applications.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Sample API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// Request represents a text sampling request.
type Request struct {
	Prompts          []string
	Model            string
	LogProbs         bool
	TopLogProbs      int32
	MaxTokens        int32
	N                int32
	PresencePenalty  float32
	Seed             int32
	Stop             []string
	FrequencyPenalty float32
	Temperature      float32
	TopP             float32
	User             string
}

// Choice represents a single completion choice.
type Choice struct {
	FinishReason string
	Index        int32
	Text         string
}

// Response represents the sampling response.
type Response struct {
	Choices []*Choice
	Model   string
}

// NewRequest creates a new sampling request.
func NewRequest(model string, prompts ...string) *Request {
	return &Request{
		Model:       model,
		Prompts:     prompts,
		MaxTokens:   100,
		Temperature: 1.0,
		N:           1,
	}
}

// WithMaxTokens sets the maximum number of tokens to generate.
func (r *Request) WithMaxTokens(maxTokens int32) *Request {
	r.MaxTokens = maxTokens
	return r
}

// WithTemperature sets the sampling temperature.
func (r *Request) WithTemperature(temperature float32) *Request {
	r.Temperature = temperature
	return r
}

// Sample generates text completions.
func (c *Client) Sample(ctx context.Context, req *Request) (*Response, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	protoReq := &xaiv1.SampleTextRequest{
		Prompt:           req.Prompts,
		Model:            req.Model,
		Logprobs:         req.LogProbs,
		TopLogprobs:      req.TopLogProbs,
		MaxTokens:        req.MaxTokens,
		N:                req.N,
		PresencePenalty:  req.PresencePenalty,
		Seed:             req.Seed,
		Stop:             req.Stop,
		FrequencyPenalty: req.FrequencyPenalty,
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		User:             req.User,
	}

	jsonData, err := protojson.Marshal(protoReq)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/completions", jsonData)
	if err != nil {
		return nil, err
	}

	var sampleResp xaiv1.SampleTextResponse
	if err := protojson.Unmarshal(resp.Body, &sampleResp); err != nil {
		return nil, err
	}

	choices := make([]*Choice, len(sampleResp.Choices))
	for i, ch := range sampleResp.Choices {
		choices[i] = &Choice{
			FinishReason: ch.FinishReason.String(),
			Index:        ch.Index,
			Text:         ch.Text,
		}
	}

	return &Response{
		Choices: choices,
		Model:   sampleResp.Model,
	}, nil
}
