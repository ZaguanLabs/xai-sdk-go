// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Helper functions to convert between user-friendly strings and proto enums

// roleToProto converts a string role to MessageRole enum
func roleToProto(role string) xaiv1.MessageRole {
	switch role {
	case "system":
		return xaiv1.MessageRole_ROLE_SYSTEM
	case "user":
		return xaiv1.MessageRole_ROLE_USER
	case "assistant":
		return xaiv1.MessageRole_ROLE_ASSISTANT
	case "function":
		return xaiv1.MessageRole_ROLE_FUNCTION
	case "tool":
		return xaiv1.MessageRole_ROLE_TOOL
	default:
		return xaiv1.MessageRole_INVALID_ROLE
	}
}

// roleFromProto converts MessageRole enum to string
func roleFromProto(role xaiv1.MessageRole) string {
	switch role {
	case xaiv1.MessageRole_ROLE_SYSTEM:
		return "system"
	case xaiv1.MessageRole_ROLE_USER:
		return "user"
	case xaiv1.MessageRole_ROLE_ASSISTANT:
		return "assistant"
	case xaiv1.MessageRole_ROLE_FUNCTION:
		return "function"
	case xaiv1.MessageRole_ROLE_TOOL:
		return "tool"
	default:
		return ""
	}
}

// reasoningEffortToProto converts a string to ReasoningEffort enum
func reasoningEffortToProto(effort string) xaiv1.ReasoningEffort {
	switch effort {
	case "low":
		return xaiv1.ReasoningEffort_EFFORT_LOW
	case "medium":
		return xaiv1.ReasoningEffort_EFFORT_MEDIUM
	case "high":
		return xaiv1.ReasoningEffort_EFFORT_HIGH
	default:
		return xaiv1.ReasoningEffort_INVALID_EFFORT
	}
}

// reasoningEffortFromProto converts ReasoningEffort enum to string
func reasoningEffortFromProto(effort xaiv1.ReasoningEffort) string {
	switch effort {
	case xaiv1.ReasoningEffort_EFFORT_LOW:
		return "low"
	case xaiv1.ReasoningEffort_EFFORT_MEDIUM:
		return "medium"
	case xaiv1.ReasoningEffort_EFFORT_HIGH:
		return "high"
	default:
		return ""
	}
}

// toolModeToProto converts a string to ToolMode enum
func toolModeToProto(mode string) xaiv1.ToolMode {
	switch mode {
	case "auto":
		return xaiv1.ToolMode_TOOL_MODE_AUTO
	case "none":
		return xaiv1.ToolMode_TOOL_MODE_NONE
	case "required":
		return xaiv1.ToolMode_TOOL_MODE_REQUIRED
	default:
		return xaiv1.ToolMode_TOOL_MODE_INVALID
	}
}

// formatTypeToProto converts a string to FormatType enum
func formatTypeToProto(format string) xaiv1.FormatType {
	switch format {
	case "text":
		return xaiv1.FormatType_FORMAT_TYPE_TEXT
	case "json_object":
		return xaiv1.FormatType_FORMAT_TYPE_JSON_OBJECT
	case "json_schema":
		return xaiv1.FormatType_FORMAT_TYPE_JSON_SCHEMA
	default:
		return xaiv1.FormatType_FORMAT_TYPE_INVALID
	}
}

// searchModeToProto converts a string to SearchMode enum
func searchModeToProto(mode string) xaiv1.SearchMode {
	switch mode {
	case "off":
		return xaiv1.SearchMode_OFF_SEARCH_MODE
	case "on":
		return xaiv1.SearchMode_ON_SEARCH_MODE
	case "auto":
		return xaiv1.SearchMode_AUTO_SEARCH_MODE
	default:
		return xaiv1.SearchMode_INVALID_SEARCH_MODE
	}
}
