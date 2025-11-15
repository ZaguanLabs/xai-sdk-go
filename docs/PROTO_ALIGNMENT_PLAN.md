# Proto Alignment Plan

**Goal**: Achieve 100% alignment between Go SDK and official xAI Python SDK v1.4.0

**Status**: In Progress  
**Started**: 2025-11-15  
**Target Completion**: TBD

---

## Executive Summary

We have extracted all proto definitions from the official xAI Python SDK v1.4.0 and identified gaps and mismatches with our Go SDK. This document tracks the systematic alignment effort.

### Overall Statistics

- **Total Proto Files in Python SDK**: 14
- **Existing in Go SDK**: 7
- **Missing in Go SDK**: 7
- **Total Messages to Align**: 107
- **Total Enums to Align**: 19

---

## Phase 1: Update Existing Proto Files

Update existing Go SDK proto files to match official Python SDK definitions.

### 1.1 auth.proto

**Status**: ❌ NEEDS UPDATE  
**Priority**: Medium  
**Complexity**: Low

**Current State**:
- Has custom `ValidateKey` service and messages
- Does not match Python SDK structure

**Python SDK Has**:
- `ApiKey` message (1 message, 0 enums)
- Fields: redacted_api_key, user_id, name, create_time, team_id, acls, api_key_id

**Action Items**:
- [ ] Replace custom auth messages with official `ApiKey` message
- [ ] Remove custom `ValidateKey` service (not in Python SDK)
- [ ] Update field numbers and types to match
- [ ] Regenerate Go code
- [ ] Update SDK auth wrapper code

**Estimated Effort**: 1-2 hours

---

### 1.2 chat.proto

**Status**: ✅ MOSTLY ALIGNED (v0.1.5)  
**Priority**: Critical  
**Complexity**: High

**Current State**:
- Message proto fixed in v0.1.5 to match Python SDK
- Has custom service definitions

**Python SDK Has**:
- 37 messages, 6 enums
- All message definitions extracted

**Known Differences**:
- Missing some enums (FormatType, etc.)
- Service definitions may differ
- Some message fields may be incomplete

**Action Items**:
- [ ] Compare all 37 messages field-by-field
- [ ] Add missing enums
- [ ] Verify all field numbers, types, and order
- [ ] Update service definitions if needed
- [ ] Regenerate Go code
- [ ] Test with proxy

**Estimated Effort**: 4-6 hours

---

### 1.3 collections.proto

**Status**: ❌ NEEDS MAJOR UPDATE  
**Priority**: Low (not currently used)  
**Complexity**: High

**Current State**:
- Has custom service and message definitions
- Significant differences from Python SDK

**Python SDK Has**:
- 23 messages, 3 enums
- Much more comprehensive than current implementation

**Action Items**:
- [ ] Replace entire file with extracted definitions
- [ ] Add all 23 messages
- [ ] Add all 3 enums
- [ ] Regenerate Go code
- [ ] Create SDK wrapper (future)

**Estimated Effort**: 6-8 hours

---

### 1.4 files.proto

**Status**: ❌ NEEDS UPDATE  
**Priority**: Medium  
**Complexity**: Medium

**Current State**:
- Has basic file operations
- Missing many messages from Python SDK

**Python SDK Has**:
- 12 messages, 2 enums
- More comprehensive file handling

**Action Items**:
- [ ] Compare existing messages with Python SDK
- [ ] Add missing messages
- [ ] Add missing enums
- [ ] Update field numbers and types
- [ ] Regenerate Go code
- [ ] Update SDK files wrapper

**Estimated Effort**: 3-4 hours

---

### 1.5 images.proto → image.proto

**Status**: ❌ NEEDS UPDATE + RENAME  
**Priority**: Low (not currently used)  
**Complexity**: Medium

**Current State**:
- Named `images.proto` (should be `image.proto`)
- May have differences from Python SDK

**Python SDK Has**:
- 4 messages, 2 enums
- Named `image.proto`

**Action Items**:
- [ ] Rename `images.proto` to `image.proto`
- [ ] Compare all messages with Python SDK
- [ ] Update field numbers and types
- [ ] Add missing enums
- [ ] Regenerate Go code
- [ ] Update SDK image wrapper

**Estimated Effort**: 2-3 hours

---

### 1.6 models.proto

**Status**: ⚠️ NEEDS VERIFICATION  
**Priority**: High (actively used)  
**Complexity**: Medium

**Current State**:
- Fixed in v0.1.1 for field order
- May still have differences

**Python SDK Has**:
- 7 messages, 1 enum
- LanguageModel, EmbeddingModel, ImageGenerationModel, etc.

**Action Items**:
- [ ] Compare all 7 messages field-by-field
- [ ] Verify field numbers and types
- [ ] Check enum definitions
- [ ] Regenerate Go code if needed
- [ ] Test with models API

**Estimated Effort**: 2-3 hours

---

### 1.7 tokenizer.proto → tokenize.proto

**Status**: ❌ NEEDS UPDATE + RENAME  
**Priority**: Low  
**Complexity**: Low

**Current State**:
- Named `tokenizer.proto` (should be `tokenize.proto`)
- May have differences from Python SDK

**Python SDK Has**:
- 3 messages, 0 enums
- Named `tokenize.proto`

**Action Items**:
- [ ] Rename `tokenizer.proto` to `tokenize.proto`
- [ ] Compare all messages with Python SDK
- [ ] Update field numbers and types
- [ ] Regenerate Go code
- [ ] Update SDK tokenizer wrapper

**Estimated Effort**: 1-2 hours

---

## Phase 2: Add Missing Proto Files

Create new proto files that exist in Python SDK but not in Go SDK.

### 2.1 deferred.proto

**Status**: ❌ MISSING  
**Priority**: Medium  
**Complexity**: Low

**Python SDK Has**:
- 2 messages, 1 enum
- `GetDeferredRequest`, `StartDeferredResponse`

**Action Items**:
- [ ] Create `deferred.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Create SDK deferred wrapper
- [ ] Add examples

**Estimated Effort**: 3-4 hours

---

### 2.2 documents.proto

**Status**: ❌ MISSING  
**Priority**: Low  
**Complexity**: Low

**Python SDK Has**:
- 4 messages, 0 enums
- Document search functionality

**Action Items**:
- [ ] Create `documents.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Create SDK documents wrapper (future)

**Estimated Effort**: 2-3 hours

---

### 2.3 embed.proto

**Status**: ❌ MISSING  
**Priority**: Medium  
**Complexity**: Medium

**Python SDK Has**:
- 5 messages, 1 enum
- Embeddings functionality

**Action Items**:
- [ ] Create `embed.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Create SDK embed wrapper
- [ ] Add examples

**Estimated Effort**: 4-5 hours

---

### 2.4 sample.proto

**Status**: ❌ MISSING  
**Priority**: Low  
**Complexity**: Low

**Python SDK Has**:
- 3 messages, 1 enum
- Text sampling functionality

**Action Items**:
- [ ] Create `sample.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Create SDK sample wrapper (future)

**Estimated Effort**: 2-3 hours

---

### 2.5 shared.proto

**Status**: ❌ MISSING  
**Priority**: High (shared types)  
**Complexity**: Low

**Python SDK Has**:
- 0 messages, 1 enum
- Shared enums and types

**Action Items**:
- [ ] Create `shared.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Update imports in other proto files
- [ ] Regenerate Go code

**Estimated Effort**: 1-2 hours

---

### 2.6 types.proto

**Status**: ❌ MISSING  
**Priority**: Medium  
**Complexity**: Low

**Python SDK Has**:
- 5 messages, 0 enums
- Type configurations (CharsConfiguration, ChunkConfiguration, etc.)

**Action Items**:
- [ ] Create `types.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Update related wrappers

**Estimated Effort**: 2-3 hours

---

### 2.7 usage.proto

**Status**: ❌ MISSING  
**Priority**: High (usage tracking)  
**Complexity**: Low

**Python SDK Has**:
- 2 messages, 1 enum
- `EmbeddingUsage`, `SamplingUsage`

**Action Items**:
- [ ] Create `usage.proto` from extracted definition
- [ ] Add to buf.yaml
- [ ] Regenerate Go code
- [ ] Update chat wrapper to use proper usage types

**Estimated Effort**: 2-3 hours

---

## Phase 3: Service Definitions

Add gRPC service definitions to match Python SDK.

### 3.1 Identify All Services

**Action Items**:
- [ ] Extract service definitions from Python SDK gRPC files
- [ ] Compare with current service definitions
- [ ] Document all RPC methods

**Estimated Effort**: 2-3 hours

---

### 3.2 Update Service Definitions

**Action Items**:
- [ ] Add missing services
- [ ] Update method signatures
- [ ] Verify request/response types
- [ ] Regenerate Go code

**Estimated Effort**: 4-6 hours

---

## Phase 4: Testing & Validation

Ensure all changes work correctly.

### 4.1 Proto Validation

**Action Items**:
- [ ] Run buf lint on all proto files
- [ ] Run buf breaking change detection
- [ ] Verify all imports resolve
- [ ] Generate Go code successfully

**Estimated Effort**: 1-2 hours

---

### 4.2 SDK Wrapper Updates

**Action Items**:
- [ ] Update all SDK wrapper code for proto changes
- [ ] Fix compilation errors
- [ ] Update method signatures
- [ ] Add new wrapper functions for new protos

**Estimated Effort**: 8-12 hours

---

### 4.3 Integration Testing

**Action Items**:
- [ ] Test chat completions (already working)
- [ ] Test models API
- [ ] Test new endpoints (deferred, embed, etc.)
- [ ] Add integration tests
- [ ] Test with proxy

**Estimated Effort**: 4-6 hours

---

### 4.4 Documentation

**Action Items**:
- [ ] Update README with new features
- [ ] Update examples
- [ ] Add godoc comments
- [ ] Update CHANGELOG

**Estimated Effort**: 2-3 hours

---

## Execution Strategy

### Recommended Order

1. **Quick Wins** (1-2 days)
   - [ ] 2.5 shared.proto (shared types needed by others)
   - [ ] 2.7 usage.proto (needed by chat)
   - [ ] 1.6 models.proto verification
   - [ ] 1.7 tokenize.proto rename + update

2. **Critical Path** (2-3 days)
   - [ ] 1.2 chat.proto complete alignment
   - [ ] 2.3 embed.proto (embeddings support)
   - [ ] 2.1 deferred.proto (deferred completions)

3. **Medium Priority** (2-3 days)
   - [ ] 1.4 files.proto update
   - [ ] 1.1 auth.proto update
   - [ ] 2.6 types.proto

4. **Low Priority** (1-2 days)
   - [ ] 1.3 collections.proto major update
   - [ ] 1.5 image.proto rename + update
   - [ ] 2.2 documents.proto
   - [ ] 2.4 sample.proto

5. **Services & Testing** (3-4 days)
   - [ ] Phase 3: Service definitions
   - [ ] Phase 4: Testing & validation

### Total Estimated Effort

- **Phase 1**: 19-28 hours
- **Phase 2**: 16-23 hours
- **Phase 3**: 6-9 hours
- **Phase 4**: 15-23 hours

**Total**: 56-83 hours (7-10 working days)

---

## Success Criteria

- [ ] All 14 proto files present and aligned with Python SDK v1.4.0
- [ ] All messages have correct field numbers, types, and order
- [ ] All enums match Python SDK
- [ ] All service definitions match Python SDK
- [ ] Go code generates without errors
- [ ] All SDK wrappers compile
- [ ] Integration tests pass
- [ ] Proxy works with updated SDK
- [ ] Documentation updated

---

## Risks & Mitigation

### Risk 1: Breaking Changes
**Impact**: High  
**Mitigation**: 
- Version bump to v0.2.0 (breaking change)
- Maintain backward compatibility where possible
- Clear migration guide in CHANGELOG

### Risk 2: Service Definition Mismatches
**Impact**: Medium  
**Mitigation**:
- Extract service definitions from Python SDK gRPC stubs
- Test each service endpoint
- Use Python SDK as reference implementation

### Risk 3: Import Dependencies
**Impact**: Medium  
**Mitigation**:
- Map all proto imports carefully
- Use buf dependency management
- Test compilation at each step

### Risk 4: Time Estimation
**Impact**: Low  
**Mitigation**:
- Work in small incremental commits
- Test after each proto update
- Adjust plan as needed

---

## Progress Tracking

### Completed
- [x] Extract all proto definitions from Python SDK v1.4.0
- [x] Create comparison tools
- [x] Fix chat.proto Message field order (v0.1.5)
- [x] Create this alignment plan

### In Progress
- [ ] None currently

### Blocked
- [ ] None currently

---

## Notes

- All extracted proto definitions are in `proto/xai/v1/*.proto.extracted`
- Use `tools/compare_protos.sh` to check alignment status
- Use `tools/verify_protos.py` to re-extract if Python SDK updates
- Reference: xAI Python SDK v1.4.0 at `docs/xai-sdk-python/dist/xai_sdk-1.4.0/`

---

**Last Updated**: 2025-11-15  
**Next Review**: After Phase 1 completion
