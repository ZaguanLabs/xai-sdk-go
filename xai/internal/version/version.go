// Package version provides version information for the xAI SDK.
package version

import (
	"fmt"
	"runtime"
	"strings"
)

// SDKVersion is the current version of the xAI SDK Go client.
// This is the single source of truth for the SDK version.
// Update this constant for new releases.
const SDKVersion = "0.10.0"

// BuildInfo contains build and runtime information.
type BuildInfo struct {
	SDKVersion string
	GoVersion  string
	GOOS       string
	GOARCH     string
	GitCommit  string
	BuildTime  string
}

// GetBuildInfo returns build information for the SDK.
func GetBuildInfo() BuildInfo {
	info := BuildInfo{
		SDKVersion: SDKVersion,
		GoVersion:  runtime.Version(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
	}

	// Try to extract build info from runtime/debug if available
	if buildInfo, ok := readBuildInfo(); ok {
		info.GitCommit = buildInfo.GitCommit
		info.BuildTime = buildInfo.BuildTime
	}

	return info
}

// String returns a formatted string representation of the build info.
func (bi BuildInfo) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("xai-sdk-go v%s", bi.SDKVersion))
	parts = append(parts, fmt.Sprintf("go %s", bi.GoVersion))
	parts = append(parts, fmt.Sprintf("%s/%s", bi.GOOS, bi.GOARCH))

	if bi.GitCommit != "" {
		parts = append(parts, fmt.Sprintf("git %s", bi.GitCommit))
	}

	if bi.BuildTime != "" {
		parts = append(parts, fmt.Sprintf("built %s", bi.BuildTime))
	}

	return strings.Join(parts, " ")
}

// GetSDKVersion returns the current SDK version.
func GetSDKVersion() string {
	return SDKVersion
}

// GetRuntimeInfo returns runtime information.
func GetRuntimeInfo() (goVersion, os, arch string) {
	return runtime.Version(), runtime.GOOS, runtime.GOARCH
}

// readBuildInfo attempts to read build information from the binary.
func readBuildInfo() (BuildInfo, bool) {
	// This is a simplified implementation
	// In a real scenario, this would parse build info from the binary
	return BuildInfo{}, false
}
