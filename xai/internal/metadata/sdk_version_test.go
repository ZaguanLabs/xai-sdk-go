package metadata

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/version"
	"google.golang.org/grpc/metadata"
)

func TestSDKVersionMetadata(t *testing.T) {
	m := NewSDKMetadata("test-api-key")
	md := m.ToMetadata()

	// Check xai-sdk-version
	sdkVersion := md.Get(SDKVersionKey)
	switch {
	case len(sdkVersion) == 0:
		t.Error("xai-sdk-version should be set")
	case !strings.HasPrefix(sdkVersion[0], "go/"):
		t.Errorf("xai-sdk-version should start with %q, got %q", "go/", sdkVersion[0])
	case !strings.Contains(sdkVersion[0], version.GetSDKVersion()):
		// Should contain the actual version
		t.Errorf("xai-sdk-version should contain %q, got %q", version.GetSDKVersion(), sdkVersion[0])
	}

	// Check xai-sdk-language
	sdkLanguage := md.Get(SDKLanguageKey)
	switch {
	case len(sdkLanguage) == 0:
		t.Error("xai-sdk-language should be set")
	case !strings.HasPrefix(sdkLanguage[0], "go/"):
		t.Errorf("xai-sdk-language should start with %q, got %q", "go/", sdkLanguage[0])
	case !strings.Contains(sdkLanguage[0], runtime.Version()):
		// Should contain the Go version
		t.Errorf("xai-sdk-language should contain %q, got %q", runtime.Version(), sdkLanguage[0])
	}
}

func TestSDKVersionInContext(t *testing.T) {
	m := NewSDKMetadata("test-api-key")
	ctx := context.Background()
	ctx = m.AddToOutgoingContext(ctx)

	// Extract metadata from context
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("Failed to extract metadata from context")
	}

	// Check xai-sdk-version
	sdkVersion := md.Get(SDKVersionKey)
	if len(sdkVersion) == 0 {
		t.Error("xai-sdk-version should be set in context")
	} else if !strings.HasPrefix(sdkVersion[0], "go/") {
		t.Errorf("xai-sdk-version should start with 'go/', got %q", sdkVersion[0])
	}

	// Check xai-sdk-language
	sdkLanguage := md.Get(SDKLanguageKey)
	if len(sdkLanguage) == 0 {
		t.Error("xai-sdk-language should be set in context")
	} else if !strings.HasPrefix(sdkLanguage[0], "go/") {
		t.Errorf("xai-sdk-language should start with 'go/', got %q", sdkLanguage[0])
	}
}

func TestSDKVersionFormat(t *testing.T) {
	// Test that the format matches Python SDK pattern
	m := NewSDKMetadata("test-api-key")
	md := m.ToMetadata()

	sdkVersion := md.Get(SDKVersionKey)[0]
	sdkLanguage := md.Get(SDKLanguageKey)[0]

	// Python format: "python/1.4.0" and "python/3.11"
	// Go format should be: "go/0.6.0" and "go/go1.23.0"

	// Verify format: language/version
	if !strings.Contains(sdkVersion, "/") {
		t.Errorf("xai-sdk-version should contain '/', got %q", sdkVersion)
	}

	if !strings.Contains(sdkLanguage, "/") {
		t.Errorf("xai-sdk-language should contain '/', got %q", sdkLanguage)
	}

	// Verify language prefix
	parts := strings.Split(sdkVersion, "/")
	if len(parts) != 2 || parts[0] != "go" {
		t.Errorf("xai-sdk-version should have format 'go/version', got %q", sdkVersion)
	}

	parts = strings.Split(sdkLanguage, "/")
	if len(parts) < 2 || parts[0] != "go" {
		t.Errorf("xai-sdk-language should have format 'go/version', got %q", sdkLanguage)
	}
}

func TestSDKVersionConsistency(t *testing.T) {
	// Verify that both ToMetadata and AddToOutgoingContext produce the same values
	m := NewSDKMetadata("test-api-key")

	// Get from ToMetadata
	md1 := m.ToMetadata()
	sdkVersion1 := md1.Get(SDKVersionKey)[0]
	sdkLanguage1 := md1.Get(SDKLanguageKey)[0]

	// Get from AddToOutgoingContext
	ctx := m.AddToOutgoingContext(context.Background())
	md2, _ := metadata.FromOutgoingContext(ctx)
	sdkVersion2 := md2.Get(SDKVersionKey)[0]
	sdkLanguage2 := md2.Get(SDKLanguageKey)[0]

	// They should match
	if sdkVersion1 != sdkVersion2 {
		t.Errorf("xai-sdk-version mismatch: ToMetadata=%q, AddToOutgoingContext=%q", sdkVersion1, sdkVersion2)
	}

	if sdkLanguage1 != sdkLanguage2 {
		t.Errorf("xai-sdk-language mismatch: ToMetadata=%q, AddToOutgoingContext=%q", sdkLanguage1, sdkLanguage2)
	}
}
