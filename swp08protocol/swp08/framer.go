package swp08

import (
	"bytes"
)

type Framer struct {
	buffer []byte
}

// Add incoming data to buffer and extract complete messages
func (f *Framer) Feed(data []byte) [][]byte {
	f.buffer = append(f.buffer, data...)

	var messages [][]byte

	for {
		start := bytes.Index(f.buffer, []byte{0x10, 0x02}) // DLE STX
		if start == -1 {
			// No start marker, discard garbage
			f.buffer = nil
			break
		}

		// Remove data before start
		if start > 0 {
			f.buffer = f.buffer[start:]
		}

		end := findFrameEnd(f.buffer)
		if end == -1 {
			// Incomplete message, wait for more
			break
		}

		// Extract raw message content (including DLE STX/ETX)
		frame := f.buffer[:end+2]
		f.buffer = f.buffer[end+2:]

		// Strip frame markers, unescape DLE DLE
		decoded := decodeMessage(frame[2 : len(frame)-2])
		messages = append(messages, decoded)
	}

	return messages
}

// Look for DLE ETX that is not part of DLE DLE
func findFrameEnd(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x10 && data[i+1] == 0x03 {
			return i
		}
	}
	return -1
}

// Replace DLE DLE with single DLE
func decodeMessage(data []byte) []byte {
	var result []byte
	skip := false
	for i := 0; i < len(data); i++ {
		if skip {
			skip = false
			continue
		}
		if data[i] == 0x10 {
			if i+1 < len(data) && data[i+1] == 0x10 {
				result = append(result, 0x10)
				skip = true
			} else {
				result = append(result, data[i])
			}
		} else {
			result = append(result, data[i])
		}
	}
	return result
}
