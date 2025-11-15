package tokenizer

import (
	"context"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

// mockTokenizerServiceClient implements TokenizerServiceClient for testing
type mockTokenizerServiceClient struct {
	encodeResp  *xaiv1.EncodeTextResponse
	decodeResp  *xaiv1.DecodeTokensResponse
	countResp   *xaiv1.CountTokensResponse
	err         error
}

func (m *mockTokenizerServiceClient) EncodeText(ctx context.Context, req *xaiv1.EncodeTextRequest, opts ...grpc.CallOption) (*xaiv1.EncodeTextResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.encodeResp != nil {
		return m.encodeResp, nil
	}
	return &xaiv1.EncodeTextResponse{
		Tokens:     []int32{1, 2, 3, 4, 5},
		TokenCount: 5,
	}, nil
}

func (m *mockTokenizerServiceClient) DecodeTokens(ctx context.Context, req *xaiv1.DecodeTokensRequest, opts ...grpc.CallOption) (*xaiv1.DecodeTokensResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.decodeResp != nil {
		return m.decodeResp, nil
	}
	return &xaiv1.DecodeTokensResponse{
		Text:       "Hello, world!",
		TokenCount: 5,
	}, nil
}

func (m *mockTokenizerServiceClient) CountTokens(ctx context.Context, req *xaiv1.CountTokensRequest, opts ...grpc.CallOption) (*xaiv1.CountTokensResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.countResp != nil {
		return m.countResp, nil
	}
	return &xaiv1.CountTokensResponse{
		TokenCount:     5,
		CharacterCount: 13,
	}, nil
}

func TestEncode(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	tokens, err := client.Encode(context.Background(), "Hello, world!", "gpt-4")

	if err != nil {
		t.Fatalf("Encode() returned error: %v", err)
	}

	if len(tokens) != 5 {
		t.Errorf("Expected 5 tokens, got %d", len(tokens))
	}

	expectedTokens := []int32{1, 2, 3, 4, 5}
	for i, token := range tokens {
		if token != expectedTokens[i] {
			t.Errorf("Expected token[%d] to be %d, got %d", i, expectedTokens[i], token)
		}
	}
}

func TestDecode(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	tokens := []int32{1, 2, 3, 4, 5}
	text, err := client.Decode(context.Background(), tokens, "gpt-4")

	if err != nil {
		t.Fatalf("Decode() returned error: %v", err)
	}

	if text != "Hello, world!" {
		t.Errorf("Expected text to be 'Hello, world!', got '%s'", text)
	}
}

func TestCount(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	count, err := client.Count(context.Background(), "Hello, world!", "gpt-4")

	if err != nil {
		t.Fatalf("Count() returned error: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count to be 5, got %d", count)
	}
}

func TestCountWithDetails(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	tokenCount, charCount, err := client.CountWithDetails(context.Background(), "Hello, world!", "gpt-4")

	if err != nil {
		t.Fatalf("CountWithDetails() returned error: %v", err)
	}

	if tokenCount != 5 {
		t.Errorf("Expected token count to be 5, got %d", tokenCount)
	}

	if charCount != 13 {
		t.Errorf("Expected character count to be 13, got %d", charCount)
	}
}

func TestEncodeEmptyText(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Encode(context.Background(), "", "gpt-4")

	if err == nil {
		t.Fatal("Expected error for empty text, got nil")
	}
}

func TestEncodeEmptyModel(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Encode(context.Background(), "Hello, world!", "")

	if err == nil {
		t.Fatal("Expected error for empty model, got nil")
	}
}

func TestDecodeEmptyTokens(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Decode(context.Background(), []int32{}, "gpt-4")

	if err == nil {
		t.Fatal("Expected error for empty tokens, got nil")
	}
}

func TestCountEmptyText(t *testing.T) {
	mockClient := &mockTokenizerServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Count(context.Background(), "", "gpt-4")

	if err == nil {
		t.Fatal("Expected error for empty text, got nil")
	}
}