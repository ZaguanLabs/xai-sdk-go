package chat

import (
	"context"
	"fmt"
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ServiceTier = xaiv1.ServiceTier

const (
	ServiceTierDefault  = xaiv1.ServiceTier_SERVICE_TIER_DEFAULT
	ServiceTierPriority = xaiv1.ServiceTier_SERVICE_TIER_PRIORITY
)

func WithServiceTier(tier ServiceTier) RequestOption {
	return func(r *Request) { r.proto.ServiceTier = tier }
}
func (r *Request) SetServiceTier(tier ServiceTier) *Request { r.proto.ServiceTier = tier; return r }
func (r *Response) ServiceTier() ServiceTier {
	tier := r.proto.GetServiceTier()
	if tier == xaiv1.ServiceTier_SERVICE_TIER_UNSPECIFIED {
		return ServiceTierDefault
	}
	return tier
}
func (c *Chunk) ServiceTier() ServiceTier {
	tier := c.proto.GetServiceTier()
	if tier == xaiv1.ServiceTier_SERVICE_TIER_UNSPECIFIED {
		return ServiceTierDefault
	}
	return tier
}
func (r *Response) OutputFiles() []*xaiv1.OutputFile { return r.proto.GetOutputFiles() }
func (c *Chunk) OutputFiles() []*xaiv1.OutputFile    { return c.proto.GetOutputFiles() }

type CompactionClient interface {
	CompactContext(context.Context, *xaiv1.CompactContextRequest, ...grpc.CallOption) (*xaiv1.CompactContextResponse, error)
}

func (c *Client) CompactContext(ctx context.Context, in *xaiv1.CompactContextRequest, opts ...grpc.CallOption) (*xaiv1.CompactContextResponse, error) {
	return c.grpcClient.CompactContext(ctx, in, opts...)
}

// Compact replaces the conversation only after a successful compaction RPC.
func (r *Request) Compact(ctx context.Context, client CompactionClient) (*xaiv1.CompactContextResponse, error) {
	if r == nil || r.proto == nil || client == nil {
		return nil, fmt.Errorf("request and compaction client required")
	}
	copy := proto.Clone(r.proto).(*xaiv1.GetCompletionsRequest)
	result, err := client.CompactContext(ctx, &xaiv1.CompactContextRequest{Model: copy.Model, Input: copy.Messages})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("empty compaction response")
	}
	r.ApplyCompaction(result)
	return result, nil
}
func (r *Request) ApplyCompaction(result *xaiv1.CompactContextResponse) *Request {
	if result != nil {
		r.proto.Messages = []*xaiv1.Message{{Role: xaiv1.MessageRole_ROLE_USER, EncryptedContent: result.EncryptedContent}}
	}
	return r
}
