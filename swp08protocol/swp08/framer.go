package swp08

const (
	DLE = 0x10
	STX = 0x02
	ETX = 0x03
)

type Framer struct {
	buffer []byte
}

func (f *Framer) Feed(data []byte) [][]byte {
	f.buffer = append(f.buffer, data...)
	var messages [][]byte

	for {
		msg, remaining := extractNextMessage(f.buffer)
		if msg == nil {
			break
		}
		messages = append(messages, msg)
		f.buffer = remaining
	}

	return messages
}

func (f *Framer) Reset() {
	f.buffer = nil
}

func extractNextMessage(buffer []byte) ([]byte, []byte) {
	start := -1
	end := -1

	// Find start: DLE STX
	for i := 0; i < len(buffer)-1; i++ {
		if buffer[i] == DLE && buffer[i+1] == STX {
			start = i + 2
			break
		}
	}
	if start == -1 {
		return nil, buffer
	}

	// Find end: DLE ETX
	for i := start; i < len(buffer)-1; i++ {
		if buffer[i] == DLE && buffer[i+1] == ETX {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, buffer
	}

	payload := buffer[start:end]
	remaining := buffer[end+2:]

	// Unescape DLE DLE → DLE
	var unescaped []byte
	for i := 0; i < len(payload); i++ {
		if payload[i] == DLE && i+1 < len(payload) && payload[i+1] == DLE {
			unescaped = append(unescaped, DLE)
			i++
		} else {
			unescaped = append(unescaped, payload[i])
		}
	}

	return unescaped, remaining
}
