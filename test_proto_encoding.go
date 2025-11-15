package main

import (
	"encoding/hex"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Create a simple message like the SDK would
	msg := &xaiv1.Message{
		Role: xaiv1.MessageRole_ROLE_USER,
		Content: []*xaiv1.Content{
			{Text: "Hello and greetings to you!"},
		},
	}

	// Encode it
	data, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encoded message (%d bytes):\n", len(data))
	fmt.Printf("Hex: %s\n", hex.EncodeToString(data))
	fmt.Printf("Bytes: %v\n", data)

	// Decode field by field
	fmt.Println("\nField analysis:")
	for i := 0; i < len(data); i++ {
		b := data[i]
		fieldNum := b >> 3
		wireType := b & 0x07
		wireTypeStr := ""
		switch wireType {
		case 0:
			wireTypeStr = "Varint"
		case 1:
			wireTypeStr = "64-bit"
		case 2:
			wireTypeStr = "LengthDelimited"
		case 3:
			wireTypeStr = "StartGroup"
		case 4:
			wireTypeStr = "EndGroup"
		case 5:
			wireTypeStr = "32-bit"
		}
		fmt.Printf("Byte %d: 0x%02x -> Field %d, WireType %d (%s)\n", i, b, fieldNum, wireType, wireTypeStr)
	}
}
