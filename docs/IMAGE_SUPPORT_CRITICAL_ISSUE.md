# 🔴 CRITICAL: Image Support Missing in Go SDK

## The Problem

**YES - Our implementation (or lack thereof) IS hindering grok-4 from getting images!**

---

## 🚨 The Critical Bug

### What the Proto Supports

```protobuf
// proto/xai/v1/chat.proto
message Content {
  string text = 1;
  ImageUrlContent image_url = 2;  // ✅ Proto SUPPORTS images
  FileContent file = 3;            // ✅ Proto SUPPORTS files
}
```

### What Our Go SDK Does

```go
// xai/chat/message.go:18-25
func NewMessage(role string, parts ...Part) *Message {
    contents := make([]*xaiv1.Content, 0, len(parts))
    for _, p := range parts {
        contents = append(contents, &xaiv1.Content{
            Text: p.Content(),  // ❌ ONLY sets Text field!
            // ImageUrl: NOT SET
            // File: NOT SET
        })
    }
    // ...
}
```

**The bug**: We're creating `Content` proto messages but **ONLY populating the `Text` field**, completely ignoring `ImageUrl` and `File` fields!

---

## 💥 Impact

### What Happens When OpenWebUI Sends Images

1. **OpenWebUI** sends a request with image data (probably in OpenAI format)
2. **Some middleware** (if any) tries to convert to xAI format
3. **Go SDK** receives the request
4. **Our code** creates `Content` messages with **ONLY the `Text` field**
5. **Images are silently dropped** - never sent to the API
6. **grok-4 never sees the images**

---

## 🔍 Evidence

### Proto Definition (CORRECT)
```protobuf
message Content {
  string text = 1;
  ImageUrlContent image_url = 2;  // ← Field exists!
  FileContent file = 3;
}

message ImageUrlContent {
  string image_url = 1;
  ImageDetail detail = 2;
}
```

### Our Implementation (BROKEN)
```go
// We create Content but only set Text:
contents = append(contents, &xaiv1.Content{
    Text: p.Content(),  // ❌ Only this
    // ImageUrl: nil,   // ❌ Never set
    // File: nil,       // ❌ Never set
})
```

### Python SDK (CORRECT)
```python
def image(image_url: str, *, detail: Optional[ImageDetail] = "auto") -> chat_pb2.Content:
    return chat_pb2.Content(
        image_url=image_pb2.ImageUrlContent(
            image_url=image_url,
            detail=pb_detail
        )
    )
```

---

## 🎯 The Fix

### Step 1: Extend the Part Interface

```go
// xai/chat/content.go

type Part interface {
    Content() string
    Type() PartType  // NEW: Identify part type
}

type PartType int

const (
    PartTypeText PartType = iota
    PartTypeImage
    PartTypeFile
)
```

### Step 2: Create ImagePart

```go
// xai/chat/content.go

type ImagePart struct {
    url    string
    detail xaiv1.ImageDetail
}

func (i *ImagePart) Content() string {
    return i.url
}

func (i *ImagePart) Type() PartType {
    return PartTypeImage
}

// Image creates an image content part
func Image(imageURL string, detail ...string) Part {
    d := xaiv1.ImageDetail_DETAIL_AUTO
    if len(detail) > 0 {
        switch detail[0] {
        case "low":
            d = xaiv1.ImageDetail_DETAIL_LOW
        case "high":
            d = xaiv1.ImageDetail_DETAIL_HIGH
        }
    }
    
    return &ImagePart{
        url:    imageURL,
        detail: d,
    }
}
```

### Step 3: Fix NewMessage to Handle All Content Types

```go
// xai/chat/message.go

func NewMessage(role string, parts ...Part) *Message {
    contents := make([]*xaiv1.Content, 0, len(parts))
    for _, p := range parts {
        switch p.Type() {
        case PartTypeImage:
            // Handle image parts
            if img, ok := p.(*ImagePart); ok {
                contents = append(contents, &xaiv1.Content{
                    ImageUrl: &xaiv1.ImageUrlContent{
                        ImageUrl: img.url,
                        Detail:   img.detail,
                    },
                })
            }
        case PartTypeFile:
            // Handle file parts
            if file, ok := p.(*FilePart); ok {
                contents = append(contents, &xaiv1.Content{
                    File: &xaiv1.FileContent{
                        FileId: file.fileID,
                    },
                })
            }
        default:
            // Handle text parts
            contents = append(contents, &xaiv1.Content{
                Text: p.Content(),
            })
        }
    }

    return &Message{
        proto: &xaiv1.Message{
            Role:    roleToProto(role),
            Content: contents,
        },
        parts: parts,
    }
}
```

### Step 4: Update TextPart

```go
// xai/chat/content.go

func (t *TextPart) Type() PartType {
    return PartTypeText
}
```

---

## 📝 Usage After Fix

```go
import (
    "github.com/ZaguanLabs/xai-sdk-go/xai"
    "github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

// Single image
req := chat.NewRequest("grok-2-vision",
    chat.WithMessage(
        chat.User(
            chat.Text("What's in this image?"),
            chat.Image("https://example.com/photo.jpg", "high"),
        ),
    ),
)

// Multiple images
req := chat.NewRequest("grok-2-vision",
    chat.WithMessage(
        chat.User(
            chat.Text("Compare these images:"),
            chat.Image("https://example.com/photo1.jpg"),
            chat.Image("https://example.com/photo2.jpg"),
        ),
    ),
)

// Base64 image
req := chat.NewRequest("grok-2-vision",
    chat.WithMessage(
        chat.User(
            chat.Text("Analyze this:"),
            chat.Image("data:image/jpeg;base64,/9j/4AAQSkZJRg..."),
        ),
    ),
)
```

---

## ⚠️ Important Notes

### Model Support

**Vision models** (support images):
- ✅ `grok-2-vision`
- ✅ `grok-2-vision-1212`
- ✅ Possibly `grok-4` (if it's a vision variant)

**Non-vision models** (don't support images):
- ❌ `grok-beta`
- ❌ `grok-2-1212`
- ❌ `grok-4-fast` (probably text-only)

### Image Requirements

- **Formats**: PNG or JPG only
- **Size**: Max 10 MiB
- **Fetch timeout**: 5 seconds
- **URL or base64**: Both supported
- **Detail levels**: `auto` (default), `low`, `high`

---

## 🔥 Why This Is Critical

1. **Silent failure**: Images are dropped without error
2. **User confusion**: Users think images are being sent but they're not
3. **Feature gap**: Python SDK works, Go SDK doesn't
4. **OpenWebUI broken**: Can't use vision features through Go SDK

---

## ✅ Action Items

### Immediate (v0.5.3 or v0.6.0)

1. ✅ Implement `ImagePart` type
2. ✅ Implement `Image()` helper function
3. ✅ Fix `NewMessage()` to handle image content
4. ✅ Add `FilePart` and `File()` for completeness
5. ✅ Add comprehensive tests
6. ✅ Add examples
7. ✅ Update documentation

### Testing

1. Test with actual image URLs
2. Test with base64 images
3. Test with multiple images
4. Test with mixed text and images
5. Test with different detail levels
6. Test error cases (invalid URLs, too large, etc.)

---

## 📊 Current Status

| Feature | Proto | Python SDK | Go SDK v0.5.2 | Status |
|---------|-------|-----------|---------------|--------|
| Text content | ✅ | ✅ | ✅ | Working |
| Image content | ✅ | ✅ | ❌ | **BROKEN** |
| File content | ✅ | ✅ | ❌ | **BROKEN** |
| ImageDetail enum | ✅ | ✅ | ❌ | Missing |
| Multiple content | ✅ | ✅ | ⚠️ | Text only |

---

## 🎯 Bottom Line

**YES - Our lack of image support IS preventing grok-4 from receiving images.**

The proto supports it, the Python SDK implements it, but our Go SDK **silently drops all non-text content**.

This needs to be fixed ASAP for vision model support.
