# Critical Bugs Found in Go SDK vs Python SDK

## 🔴 CRITICAL BUG #1: AppendMessage doesn't support Response objects

**Python SDK**:
```python
chat.append(response)  # Accepts Response objects
```

**Go SDK**:
```go
req.AppendMessage(msg)  // Only accepts Message objects
```

**Impact**: Cannot properly append assistant responses to conversation history, which breaks multi-turn conversations with tool calls.

**What's missing**:
1. `AppendResponse()` method that accepts `*Response`
2. When appending Response, must extract and include:
   - `tool_calls` from the response
   - `reasoning_content` from the response
   - `encrypted_content` from the response

---

## 🔴 CRITICAL BUG #2: Response missing ReasoningContent() and EncryptedContent() accessors

**Python SDK**:
```python
response.reasoning_content  # Property accessor
response.encrypted_content  # Property accessor
response.tool_calls  # Property accessor (we fixed this!)
```

**Go SDK**:
```go
response.Content()    // ✅ Exists
response.ToolCalls()  // ✅ Fixed in this session
response.ReasoningContent()   // ❌ MISSING
response.EncryptedContent()   // ❌ MISSING
```

**Impact**: Cannot access reasoning content or encrypted content from responses, which breaks:
- Reasoning model interactions
- Zero Data Retention (ZDR) workflows
- Conversation continuity with encrypted content

---

## 🔴 CRITICAL BUG #3: Chunk missing ReasoningContent() and EncryptedContent() accessors

**Python SDK**:
```python
chunk.reasoning_content  # Property accessor
chunk.encrypted_content  # Property accessor (for streaming)
```

**Go SDK**:
```go
chunk.Content()    // ✅ Exists
chunk.ToolCalls()  // ✅ Fixed in this session
chunk.ReasoningContent()   // ❌ MISSING
chunk.EncryptedContent()   // ❌ MISSING
```

**Impact**: Cannot access reasoning or encrypted content during streaming.

---

## 🔴 CRITICAL BUG #4: Message doesn't support tool_calls, reasoning_content, encrypted_content

**Proto Definition** (chat.proto):
```protobuf
message Message {
  repeated Content content = 1;
  MessageRole role = 2;
  string name = 3;
  repeated ToolCall tool_calls = 4;        // ❌ Not exposed in Go
  string reasoning_content = 5;             // ❌ Not exposed in Go
  string encrypted_content = 6;             // ❌ Not exposed in Go
}
```

**Go SDK Message**:
```go
type Message struct {
    proto *xaiv1.Message
    parts []Part
}
// Only exposes: Role(), Content(), Parts()
// Missing: ToolCalls(), ReasoningContent(), EncryptedContent()
```

**Impact**: Cannot create messages with tool calls, reasoning content, or encrypted content when building conversation history.

---

## 🟡 MEDIUM BUG #5: No helper to create assistant messages with tool calls

**Python SDK**:
```python
# When appending a response, it automatically creates an assistant message with:
# - content
# - tool_calls
# - reasoning_content
# - encrypted_content
```

**Go SDK**: No equivalent helper function.

**Impact**: Users must manually construct complex message structures.

---

## 🟡 MEDIUM BUG #6: FinishReason() returns string instead of enum

**Python SDK**:
```python
response.finish_reason  # Returns sample_pb2.FinishReason enum
```

**Go SDK**:
```go
response.FinishReason()  // Returns string
```

**Impact**: Less type-safe, harder to check specific finish reasons programmatically.

---

## 🟡 MEDIUM BUG #7: No multi-output support (N > 1)

**Python SDK**:
```python
# Automatically detects multi-output mode when N > 1
if message._index is None:
    # Every single output should be appended for agentic tool call responses.
    for output in message.proto.outputs:
        # Append each output as a separate message
```

**Go SDK**: No equivalent logic for handling multiple outputs.

**Impact**: When using `N > 1` parameter, responses with multiple choices aren't properly handled.

---

## Summary

| Bug | Severity | Status | Impact |
|-----|----------|--------|--------|
| AppendMessage doesn't support Response | 🔴 Critical | ❌ Not Fixed | Breaks multi-turn conversations |
| Response missing ReasoningContent() | 🔴 Critical | ❌ Not Fixed | Cannot access reasoning |
| Response missing EncryptedContent() | 🔴 Critical | ❌ Not Fixed | Breaks ZDR workflows |
| Chunk missing ReasoningContent() | 🔴 Critical | ❌ Not Fixed | Cannot stream reasoning |
| Chunk missing EncryptedContent() | 🔴 Critical | ❌ Not Fixed | Cannot stream encrypted content |
| Message missing tool_calls support | 🔴 Critical | ❌ Not Fixed | Cannot build proper conversation history |
| Message missing reasoning_content | 🔴 Critical | ❌ Not Fixed | Cannot include reasoning in messages |
| Message missing encrypted_content | 🔴 Critical | ❌ Not Fixed | Cannot include encrypted content |
| No assistant message helper with tool calls | 🟡 Medium | ❌ Not Fixed | Poor developer experience |
| FinishReason returns string not enum | 🟡 Medium | ❌ Not Fixed | Less type-safe |
| No multi-output support | 🟡 Medium | ❌ Not Fixed | N > 1 doesn't work properly |

**Total Critical Bugs**: 8  
**Total Medium Bugs**: 3  
**Total Bugs**: 11

---

## Action Plan

### Phase 1: Response/Chunk Accessors (Critical)
1. Add `ReasoningContent()` to Response
2. Add `EncryptedContent()` to Response
3. Add `ReasoningContent()` to Chunk
4. Add `EncryptedContent()` to Chunk

### Phase 2: Message Enhancement (Critical)
5. Add `WithToolCalls()` to Message
6. Add `WithReasoningContent()` to Message
7. Add `WithEncryptedContent()` to Message
8. Add `ToolCalls()` accessor to Message
9. Add `ReasoningContent()` accessor to Message
10. Add `EncryptedContent()` accessor to Message

### Phase 3: AppendResponse (Critical)
11. Add `AppendResponse(*Response)` method to Request
12. Implement proper extraction of tool_calls, reasoning_content, encrypted_content
13. Handle multi-output mode (when N > 1)

### Phase 4: Helpers (Medium)
14. Add `AssistantWithToolCalls()` helper function
15. Consider enum-based FinishReason

---

## Testing Requirements

Each fix must include:
1. Unit tests verifying the functionality
2. Integration tests with real proto messages
3. Documentation updates
4. Example code demonstrating usage
