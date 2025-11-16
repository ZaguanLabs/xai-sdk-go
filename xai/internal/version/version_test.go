package version

import (
	"runtime"
	"testing"
	"time"
)

func TestSDKVersion(t *testing.T) {
	version := GetSDKVersion()
	if version == "" {
		t.Error("SDK version should not be empty")
	}

	expectedVersion := "0.2.1"
	if version != expectedVersion {
		t.Errorf("Expected SDK version %s, got %s", expectedVersion, version)
	}
}

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()

	if info.SDKVersion == "" {
		t.Error("BuildInfo SDKVersion should not be empty")
	}

	if info.GoVersion == "" {
		t.Error("BuildInfo GoVersion should not be empty")
	}

	if info.GOOS == "" {
		t.Error("BuildInfo GOOS should not be empty")
	}

	if info.GOARCH == "" {
		t.Error("BuildInfo GOARCH should not be empty")
	}

	// Verify GoVersion matches runtime version
	if info.GoVersion != runtime.Version() {
		t.Errorf("Expected GoVersion %s, got %s", runtime.Version(), info.GoVersion)
	}

	// Verify GOOS matches runtime GOOS
	if info.GOOS != runtime.GOOS {
		t.Errorf("Expected GOOS %s, got %s", runtime.GOOS, info.GOOS)
	}

	// Verify GOARCH matches runtime GOARCH
	if info.GOARCH != runtime.GOARCH {
		t.Errorf("Expected GOARCH %s, got %s", runtime.GOARCH, info.GOARCH)
	}
}

func TestBuildInfoString(t *testing.T) {
	info := BuildInfo{
		SDKVersion: "0.2.1",
		GoVersion:  "go1.21.0",
		GOOS:       "linux",
		GOARCH:     "amd64",
		GitCommit:  "abc123",
		BuildTime:  time.Now().Format(time.RFC3339),
	}

	str := info.String()

	// Check that the string contains expected parts
	if !containsSubstring(str, "xai-sdk-go v0.2.1") {
		t.Errorf("Expected string to contain SDK version, got: %s", str)
	}

	if !containsSubstring(str, "go1.21.0") {
		t.Errorf("Expected string to contain Go version, got: %s", str)
	}

	if !containsSubstring(str, "linux/amd64") {
		t.Errorf("Expected string to contain OS/ARCH, got: %s", str)
	}

	if !containsSubstring(str, "git abc123") {
		t.Errorf("Expected string to contain git commit, got: %s", str)
	}
}

func TestGetRuntimeInfo(t *testing.T) {
	goVersion, os, arch := GetRuntimeInfo()

	if goVersion == "" {
		t.Error("Go version should not be empty")
	}

	if os == "" {
		t.Error("OS should not be empty")
	}

	if arch == "" {
		t.Error("ARCH should not be empty")
	}

	// Verify values match runtime
	if goVersion != runtime.Version() {
		t.Errorf("Expected Go version %s, got %s", runtime.Version(), goVersion)
	}

	if os != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, os)
	}

	if arch != runtime.GOARCH {
		t.Errorf("Expected ARCH %s, got %s", runtime.GOARCH, arch)
	}
}

func TestEmptyBuildInfo(t *testing.T) {
	info := BuildInfo{}
	str := info.String()

	// Should still contain SDK version at minimum
	if !containsSubstring(str, "xai-sdk-go") {
		t.Errorf("Expected string to contain SDK version, got: %s", str)
	}

	// Should contain go version
	if !containsSubstring(str, "go ") {
		t.Errorf("Expected string to contain Go version, got: %s", str)
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkGetBuildInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetBuildInfo()
	}
}

func BenchmarkGetRuntimeInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = GetRuntimeInfo()
	}
}
