package swp08

// Constructs a full SW-P-08 response message
func constructResponse(command byte, message []byte) []byte {
	byteCount := byte(1 + len(message)) // 1 for command + N message bytes
	checksum := calculateChecksum(command, message, byteCount)

	// Build payload: [Command][Message][ByteCount][Checksum]
	payload := append([]byte{command}, message...)
	payload = append(payload, byteCount)
	payload = append(payload, checksum)

	// Escape DLE bytes (0x10) by doubling them
	escaped := escapeDLE(payload)

	// Wrap with DLE STX and DLE ETX
	framed := []byte{0x10, 0x02}        // DLE STX
	framed = append(framed, escaped...) // Payload
	framed = append(framed, 0x10, 0x03) // DLE ETX

	return framed
}

// Escapes any 0x10 byte by doubling it: 0x10 → 0x10 0x10
func escapeDLE(data []byte) []byte {
	var escaped []byte
	for _, b := range data {
		escaped = append(escaped, b)
		if b == 0x10 {
			escaped = append(escaped, 0x10)
		}
	}
	return escaped
}
