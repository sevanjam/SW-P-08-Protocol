package swp08

import (
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	framer := Framer{}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("[SWP-08] Connection closed by", conn.RemoteAddr())
			} else {
				fmt.Println("[SWP-08] Read error:", err)
			}
			break
		}

		data := buffer[:n]

		fmt.Printf("[SWP-08] Received %d bytes:\n", n)
		printHex(data)

		messages := framer.Feed(data)

		for _, msg := range messages {
			fmt.Printf("[SWP-08] Extracted %d-byte message:\n", len(msg))
			printHex(msg)

			valid, reason := validateMessage(msg)
			if valid {
				fmt.Println("[SWP-08] ✅ Message is valid")
				sendAcknowledge(conn)
				handleMessage(conn, msg)
				// We’ll add handleMessage(msg) here next
			} else {
				fmt.Println("[SWP-08] ❌ Invalid message:", reason)
				sendNegativeAcknowledge(conn)
			}
		}
	}
}
