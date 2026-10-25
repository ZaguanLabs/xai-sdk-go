package chat

import (
	"context"
	"errors"
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc"
	"testing"
)

type compactStub struct{ fail bool }

func (s compactStub) CompactContext(_ context.Context, r *xaiv1.CompactContextRequest, _ ...grpc.CallOption) (*xaiv1.CompactContextResponse, error) {
	if s.fail {
		return nil, errors.New("failed")
	}
	return &xaiv1.CompactContextResponse{EncryptedContent: "compact", DroppedMessageCount: uint32(len(r.Input))}, nil
}
func TestCompactionAndInlineFiles(t *testing.T) {
	data := []byte("content")
	m := User(FileData(data, "file.txt", "text/plain"), FileURL("https://example.test/file"))
	data[0] = 'x'
	if string(m.Proto().Content[0].GetFile().Data) != "content" || m.Proto().Content[1].GetFile().Url != "https://example.test/file" {
		t.Fatal(m.Proto())
	}
	r := NewRequest("model", WithServiceTier(ServiceTierPriority), WithInclude(IncludeCodeExecutionFilesOutput))
	r.AppendMessage(*m)
	if _, err := r.Compact(context.Background(), compactStub{fail: true}); err == nil || len(r.Proto().Messages) != 1 || r.Proto().Messages[0].EncryptedContent != "" {
		t.Fatal("failed compaction mutated input")
	}
	result, err := r.Compact(context.Background(), compactStub{})
	if err != nil || result.DroppedMessageCount != 1 {
		t.Fatal(result, err)
	}
	if r.Proto().Messages[0].Role != xaiv1.MessageRole_ROLE_USER || r.Proto().Messages[0].EncryptedContent != "compact" || r.Proto().ServiceTier != ServiceTierPriority {
		t.Fatal(r.Proto())
	}
}
