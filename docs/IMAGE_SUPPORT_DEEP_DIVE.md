# Image Support Deep Dive: Python SDK vs Go SDK

## 🔍 How the Python SDK Handles Images in Chat

### Overview

The Python SDK supports images in chat messages through the `image()` helper function, which creates a `Content` object with an `ImageUrlContent` field.

---

## 📋 Python SDK Implementation

### 1. Content Structure (Proto)

```protobuf
// chat.proto
message Content {
  string text = 1;
  ImageUrlContent image_url = 2;  // ← Image support
  FileContent file = 3;
}

message ImageUrlContent {
  string image_url = 1;
  ImageDetail detail = 2;
}

enum ImageDetail {
  DETAIL_INVALID = 0;
  DETAIL_AUTO = 1;
  DETAIL_LOW = 2;
  DETAIL_HIGH = 3;
}
```

### 2. Python SDK `image()` Function

**Location**: `xai_sdk/chat.py:687-710`

```python
def image(image_url: str, *, detail: Optional[ImageDetail] = "auto") -> chat_pb2.Content:
    """Creates a new content object of type image for use in chat messages.

    Args:
        image_url: The URL or base64-encoded string of the image. Supported formats are PNG and JPG.
            If a URL is provided, the image is fetched for each API request without caching.
            Fetching uses the "XaiImageApiFetch/1.0" user agent with a 5-second timeout.
            The maximum image size is 10 MiB; larger images or failed fetches will cause the API request to fail.
        detail: Specifies the image resolution for model processing. One of:
        - `"auto"`: The system selects an appropriate resolution (default).
        - `"low"`: Uses a low-resolution image, reducing token usage and increasing speed.
        - `"high"`: Uses a high-resolution image, increasing token usage and processing time
            but capturing more detail.

    Returns:
        A `chat_pb2.Content` object representing the image content.
    """
    pb_detail = image_pb2.ImageDetail.DETAIL_AUTO
    if detail == "low":
        pb_detail = image_pb2.ImageDetail.DETAIL_LOW
    elif detail == "high":
        pb_detail = image_pb2.ImageDetail.DETAIL_HIGH

    return chat_pb2.Content(image_url=image_pb2.ImageUrlContent(image_url=image_url, detail=pb_detail))
```

### 3. Usage in Messages

**Message constructors** accept multiple `Content` objects:

```python
def user(*args: Content) -> chat_pb2.Message:
    """Creates a new message of role "user"."""
    return chat_pb2.Message(
        role=chat_pb2.MessageRole.ROLE_USER, 
        content=[_process_content(c) for c in args]
    )
```

**Content type** is a union:
```python
Content = Union[str, chat_pb2.Content]
```

**Auto-conversion** via `_process_content()`:
```python
def _process_content(content: Content) -> chat_pb2.Content:
    """Converts a `Content` type to a proto."""
    if isinstance(content, str):
        return text(content)  # Auto-wrap strings
    else:
        return content  # Already a Content proto
```

### 4. Example Usage

```python
from xai_sdk import Client
from xai_sdk.chat import user, image

client = Client()

# Create a message with text and image
chat = client.chat.create(
    model="grok-vision-beta",
    messages=[
        user(
            "What's in this image?",
            image("https://example.com/photo.jpg", detail="high")
        )
    ]
)

response = chat.sample()
print(response.content)
```

**Multiple images**:
```python
messages=[
    user(
        "Compare these two images:",
        image("https://example.com/photo1.jpg"),
        image("https://example.com/photo2.jpg")
    )
]
```

**Base64 images**:
```python
messages=[
    user(
        "Analyze this image:",
        image("data:image/jpeg;base64,/9j/4AAQSkZJRg...")
    )
]
```

---

## 🔴 Go SDK Current Status

### What's Missing

The Go SDK **does NOT support images in chat** currently:

1. **No `image()` helper function**
2. **Content proto is simplified** (only text)
3. **No ImageUrlContent support**
4. **No ImageDetail enum**

### Current Go SDK Content Structure

```go
// xai/chat/content.go
type Part interface {
    Content() string  // Only returns string!
}

type TextPart struct {
    text string
}

func Text(text string) Part {
    return &TextPart{text: text}
}
```

**Proto conversion** (in `message.go`):
```go
func NewMessage(role string, parts ...Part) *Message {
    contents := make([]*xaiv1.Content, 0, len(parts))
    for _, p := range parts {
        contents = append(contents, &xaiv1.Content{
            Text: p.Content(),  // ❌ Only sets Text field!
        })
    }
    // ...
}
```

---

## 🚨 Why Images Might Not Work in OpenWebUI

### Potential Issues

1. **OpenWebUI sends images in OpenAI format**
   - OpenWebUI likely uses OpenAI's API format
   - OpenAI format: `{"type": "image_url", "image_url": {"url": "..."}}`
   - xAI format: `Content{image_url: ImageUrlContent{image_url: "...", detail: ...}}`

2. **Format mismatch**
   - If OpenWebUI is using the OpenAI-compatible endpoint, it might not be translating the image format correctly
   - The xAI API expects proto format, not OpenAI JSON format

3. **Model support**
   - Not all models support images
   - Only vision models like `grok-vision-beta` or `grok-2-vision-1212` support images
   - If OpenWebUI is using `grok-beta` or `grok-2-1212`, images won't work

4. **Missing image_url field**
   - If the request doesn't include the `image_url` field in the `Content` proto, the API will ignore the image

---

## 🔧 How to Fix Image Support in Go SDK

### Implementation Plan

#### 1. Add ImagePart Type

```go
// xai/chat/content.go

type ImagePart struct {
    url    string
    detail xaiv1.ImageDetail
}

func (i *ImagePart) Content() string {
    return i.url
}

func (i *ImagePart) ImageURL() string {
    return i.url
}

func (i *ImagePart) Detail() xaiv1.ImageDetail {
    return i.detail
}

// IsImage returns true if this is an image part
func (i *ImagePart) IsImage() bool {
    return true
}
```

#### 2. Add Image() Helper Function

```go
// xai/chat/content.go

type ImageDetail string

const (
    ImageDetailAuto ImageDetail = "auto"
    ImageDetailLow  ImageDetail = "low"
    ImageDetailHigh ImageDetail = "high"
)

// Image creates an image content part.
// The image_url can be:
// - A URL to an image (PNG or JPG)
// - A base64-encoded data URI (e.g., "data:image/jpeg;base64,...")
//
// The detail parameter controls image resolution:
// - "auto": System selects appropriate resolution (default)
// - "low": Low resolution, faster, uses fewer tokens
// - "high": High resolution, slower, uses more tokens but captures more detail
func Image(imageURL string, detail ...ImageDetail) Part {
    d := ImageDetailAuto
    if len(detail) > 0 {
        d = detail[0]
    }

    var pbDetail xaiv1.ImageDetail
    switch d {
    case ImageDetailLow:
        pbDetail = xaiv1.ImageDetail_DETAIL_LOW
    case ImageDetailHigh:
        pbDetail = xaiv1.ImageDetail_DETAIL_HIGH
    default:
        pbDetail = xaiv1.ImageDetail_DETAIL_AUTO
    }

    return &ImagePart{
        url:    imageURL,
        detail: pbDetail,
    }
}
```

#### 3. Update Part Interface

```go
// xai/chat/content.go

type Part interface {
    Content() string
    IsImage() bool  // Add this
}

// Update TextPart
func (t *TextPart) IsImage() bool {
    return false
}
```

#### 4. Update Message Creation

```go
// xai/chat/message.go

func NewMessage(role string, parts ...Part) *Message {
    contents := make([]*xaiv1.Content, 0, len(parts))
    for _, p := range parts {
        if p.IsImage() {
            // Handle image parts
            if imgPart, ok := p.(*ImagePart); ok {
                contents = append(contents, &xaiv1.Content{
                    ImageUrl: &xaiv1.ImageUrlContent{
                        ImageUrl: imgPart.ImageURL(),
                        Detail:   imgPart.Detail(),
                    },
                })
            }
        } else {
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

#### 5. Usage Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/ZaguanLabs/xai-sdk-go/xai"
    "github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func main() {
    client := xai.NewClient("your-api-key")
    defer client.Close()

    req := chat.NewRequest("grok-vision-beta",
        chat.WithMessage(
            chat.User(
                chat.Text("What's in this image?"),
                chat.Image("https://example.com/photo.jpg", chat.ImageDetailHigh),
            ),
        ),
    )

    response, err := client.Chat().Sample(context.Background(), req)
    if err != nil {
        panic(err)
    }

    fmt.Println(response.Content())
}
```

---

## 🔍 Debugging OpenWebUI Image Issues

### Steps to Debug

1. **Check the model being used**
   ```
   - Is it a vision model? (grok-vision-beta, grok-2-vision-1212)
   - Non-vision models will ignore images
   ```

2. **Inspect the request format**
   ```
   - Enable debug logging in OpenWebUI
   - Check if images are being sent in the request
   - Verify the format matches xAI's proto format
   ```

3. **Check image URL accessibility**
   ```
   - Can the xAI API fetch the image?
   - Is it within the 10 MiB limit?
   - Is it PNG or JPG format?
   - Does it return within 5 seconds?
   ```

4. **Test with xAI Python SDK directly**
   ```python
   from xai_sdk import Client
   from xai_sdk.chat import user, image

   client = Client(api_key="your-key")
   
   chat = client.chat.create(
       model="grok-vision-beta",
       messages=[
           user(
               "What's in this image?",
               image("https://your-image-url.jpg")
           )
       ]
   )
   
   response = chat.sample()
   print(response.content)
   ```

5. **Check OpenWebUI configuration**
   ```
   - Is it using the correct xAI endpoint?
   - Is it translating OpenAI format to xAI format?
   - Check OpenWebUI logs for errors
   ```

---

## 📊 Feature Comparison

| Feature | Python SDK | Go SDK v0.5.2 | Status |
|---------|-----------|---------------|--------|
| Text content | ✅ `text()` | ✅ `Text()` | ✅ Implemented |
| Image content | ✅ `image()` | ❌ Missing | ❌ Not implemented |
| File content | ✅ `file()` | ❌ Missing | ❌ Not implemented |
| ImageDetail enum | ✅ auto/low/high | ❌ Missing | ❌ Not implemented |
| Multiple content types | ✅ Yes | ⚠️ Text only | ⚠️ Partial |
| Base64 images | ✅ Yes | ❌ Missing | ❌ Not implemented |
| Image URLs | ✅ Yes | ❌ Missing | ❌ Not implemented |

---

## 🎯 Recommendations

### For OpenWebUI Issue

1. **Verify model**: Ensure you're using a vision model (e.g., `grok-vision-beta`)
2. **Check format**: OpenWebUI might need to translate OpenAI format to xAI proto format
3. **Test directly**: Use xAI Python SDK to verify images work with your image URLs
4. **Check logs**: Enable debug logging in OpenWebUI to see the actual requests

### For Go SDK

1. **Implement image support** following the plan above
2. **Add file support** similarly (using `FileContent` proto)
3. **Add comprehensive tests** for image and file content
4. **Update documentation** with examples

---

## 📝 Key Takeaways

1. **Python SDK uses proto `Content` message** with three fields: `text`, `image_url`, `file`
2. **Images require `ImageUrlContent`** with `image_url` and `detail` fields
3. **Go SDK currently only supports text** - images are not implemented
4. **OpenWebUI likely has format translation issues** between OpenAI and xAI formats
5. **Vision models are required** for image support (grok-vision-beta, grok-2-vision-1212)

---

## 🔗 Related Proto Definitions

- `proto/xai/v1/chat.proto` - Content, ImageUrlContent, FileContent
- `proto/xai/v1/image.proto` - ImageDetail enum, ImageUrlContent
- `xai_sdk/chat.py:687-710` - Python SDK image() function
- `xai_sdk/chat.py:557-569` - Python SDK message constructors
