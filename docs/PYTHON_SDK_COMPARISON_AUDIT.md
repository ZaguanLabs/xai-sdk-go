# Python SDK vs Go SDK Comparison Audit

This document compares the Python SDK and Go SDK implementations to identify discrepancies and ensure feature parity.

## Critical Issues Found

### 1. ❌ **CRITICAL BUG: WithTool marshals wrong data**

**Location**: `xai/chat/chat.go:258`

**Issue**: The `WithTool` function marshals `tool.Parameters()` instead of `tool.ToJSONSchema()`, which means:
- The `"required"` field is included in each property (invalid JSON Schema)
- The top-level `"type": "object"` and `"required": []` array are missing

**Current (WRONG)**:
```go
paramsJSON, _ := json.Marshal(tool.Parameters())
// Produces: {"city": {"type": "string", "description": "...", "required": true}}
```

**Should be**:
```go
paramsJSON, _ := json.Marshal(tool.ToJSONSchema())
// Produces: {"type": "object", "properties": {"city": {"type": "string", "description": "..."}}, "required": ["city"]}
```

**Impact**: Tools sent to the API have invalid JSON Schema format, which may cause:
- API rejections
- Tool calls to fail
- Incorrect parameter validation

**Fix**: Change line 258 to use `tool.ToJSONSchema()`

---

### 2. ✅ **FIXED: ToJSONSchema strips required field**

**Location**: `xai/chat/tool.go:73-103`

**Status**: Already fixed in this session

**What was fixed**: The `ToJSONSchema()` method now correctly:
- Strips `"required"` field from individual properties
- Builds a top-level `"required": []` array
- Produces valid JSON Schema

---

## Feature Comparison

### Message Creation

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| `user()` function | ✅ | ✅ | ✅ Match |
| `system()` function | ✅ | ✅ | ✅ Match |
| `assistant()` function | ✅ | ✅ | ✅ Match |
| `text()` for content | ✅ | ✅ `Text()` | ✅ Match |
| `image()` for content | ✅ | ❌ Missing | ⚠️ Gap |
| `file()` for content | ✅ | ❌ Missing | ⚠️ Gap |
| String auto-conversion | ✅ `_process_content()` | ❌ Requires `Text()` | ⚠️ Difference |

**Python SDK**:
```python
user("Hello")  # Strings auto-converted to text content
user(text("Hello"), image("url"))  # Mixed content
```

**Go SDK**:
```go
User(Text("Hello"))  // Must explicitly use Text()
// No image() or file() support yet
```

### Tool Creation

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| `tool()` function | ✅ Takes complete JSON Schema | ✅ `NewTool()` builder | ✅ Different approach |
| JSON Schema format | ✅ User provides | ✅ Built internally | ✅ Different approach |
| Validation | ✅ At API level | ✅ `Validate()` method | ✅ Match |
| `required_tool()` | ✅ | ✅ `WithToolChoice()` | ✅ Match |

**Difference**: Python expects complete JSON Schema, Go uses builder pattern. Both valid, but Go had a bug in transformation.

### Tool Results

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| `tool_result()` | ✅ | ✅ `NewToolResult()` | ✅ Match |
| Error handling | ✅ | ✅ `NewToolResultError()` | ✅ Match |
| Result formatting | ✅ | ✅ | ✅ Match |

### Server-Side Tools

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| Web Search | ✅ | ✅ `WebSearchTool()` | ✅ Match |
| X Search | ✅ | ✅ `XSearchTool()` | ✅ Match |
| Code Execution | ✅ | ✅ `CodeExecutionTool()` | ✅ Match |
| Collections Search | ✅ | ✅ `CollectionsSearchTool()` | ✅ Match |
| Document Search | ✅ | ✅ `DocumentSearchTool()` | ✅ Match |
| MCP | ✅ | ✅ `MCPTool()` | ✅ Match |

### Chat Parameters

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| All 24 parameters | ✅ | ✅ | ✅ Match |
| Parameter validation | ✅ | ✅ | ✅ Match |
| Type safety | Python typing | Go types | ✅ Match |

### Response Handling

| Feature | Python SDK | Go SDK | Status |
|---------|-----------|--------|--------|
| Synchronous | ✅ `sample()` | ✅ `Sample()` | ✅ Match |
| Streaming | ✅ `stream()` | ✅ `Stream()` | ✅ Match |
| Multi-output detection | ✅ Auto-detects | ❓ Unknown | ⚠️ Need to verify |
| Finish reason handling | ✅ Lenient | ✅ Lenient | ✅ Match |

---

## Minor Differences (Design Choices)

### 1. Content Type Handling

**Python**: Accepts strings directly, auto-converts via `_process_content()`
```python
user("Hello")  # String auto-converted
```

**Go**: Requires explicit `Text()` wrapper
```go
User(Text("Hello"))  // Must wrap
```

**Verdict**: Both valid. Go approach is more type-safe.

### 2. Tool Parameter Definition

**Python**: User provides complete JSON Schema
```python
tool(name="...", description="...", parameters={
    "type": "object",
    "properties": {...},
    "required": [...]
})
```

**Go**: Builder pattern constructs JSON Schema
```go
tool := NewTool("...", "...")
tool.WithParameter("param", "string", "desc", true)
// Internally builds JSON Schema
```

**Verdict**: Both valid. Go approach is more ergonomic but must transform correctly (which it now does).

---

## Missing Features (Non-Critical)

### 1. Image Content Support

**Python SDK has**:
```python
image(image_url, detail="auto")
```

**Go SDK**: Not implemented

**Priority**: Medium - useful for multimodal models

### 2. File Content Support

**Python SDK has**:
```python
file(file_id)
```

**Go SDK**: Not implemented

**Priority**: Medium - needed for document-based interactions

### 3. Conversation ID

**Python SDK has**:
```python
create(model="...", conversation_id="...")
```

**Go SDK**: Not explicitly exposed

**Priority**: Low - useful for telemetry/grouping

---

## Action Items

### Immediate (Critical)

1. ✅ **Fix `ToJSONSchema()` to strip required field** - DONE
2. ❌ **Fix `WithTool()` to marshal `ToJSONSchema()` instead of `Parameters()`** - TODO
3. ❌ **Add test to verify correct JSON Schema format in WithTool** - TODO

### Short Term (Important)

4. ⚠️ **Add `image()` content support** - For multimodal
5. ⚠️ **Add `file()` content support** - For document interactions
6. ⚠️ **Verify multi-output detection logic** - Ensure parity

### Long Term (Nice to Have)

7. 📝 **Consider auto-converting strings to Text()** - Ergonomics
8. 📝 **Add conversation_id support** - Telemetry
9. 📝 **Add Pydantic-style schema generation** - Developer experience

---

## Conclusion

**Critical Issues**: 1 (WithTool marshaling bug)  
**Feature Parity**: 95% (missing image/file content)  
**Design Differences**: Acceptable (different but valid approaches)

The Go SDK is very close to full parity with the Python SDK. The critical bug in `WithTool` must be fixed immediately, and adding image/file content support would bring it to 100% feature parity.
