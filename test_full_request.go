package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Create a request exactly like your proxy would
	req := chat.NewRequest("grok-4-0709",
		chat.WithMessages(
			chat.User(chat.Text("Hello and greetings to you!")),
		),
	)

	// Get the proto
	protoReq := req.Proto()

	fmt.Println("=== Request Proto ===")
	fmt.Printf("Model: %s\n", protoReq.Model)
	fmt.Printf("Messages count: %d\n", len(protoReq.Messages))

	if len(protoReq.Messages) > 0 {
		msg := protoReq.Messages[0]
		fmt.Printf("\nMessage[0]:\n")
		fmt.Printf("  Role: %v (%d)\n", msg.Role, msg.Role)
		fmt.Printf("  Content: %q\n", msg.Content)
		fmt.Printf("  Content type: %T\n", msg.Content)

		// Encode just the message
		msgData, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\n=== Encoded Message (%d bytes) ===\n", len(msgData))
		fmt.Printf("Hex: %s\n", hex.EncodeToString(msgData))
		fmt.Printf("Bytes: %v\n", msgData)

		// Analyze wire format
		fmt.Println("\n=== Wire Format Analysis ===")
		for i := 0; i < len(msgData) && i < 10; i++ {
			b := msgData[i]
			fieldNum := b >> 3
			wireType := b & 0x07
			wireTypeStr := ""
			switch wireType {
			case 0:
				wireTypeStr = "Varint"
			case 2:
				wireTypeStr = "LengthDelimited"
			default:
				wireTypeStr = fmt.Sprintf("Type%d", wireType)
			}
			fmt.Printf("Byte %d: 0x%02x -> Field %d, WireType %d (%s)\n", i, b, fieldNum, wireType, wireTypeStr)
		}
	}

	// Encode the full request
	reqData, err := proto.Marshal(protoReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n=== Full Request (%d bytes) ===\n", len(reqData))
	fmt.Printf("First 50 bytes hex: %s\n", hex.EncodeToString(reqData[:min(50, len(reqData))]))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
