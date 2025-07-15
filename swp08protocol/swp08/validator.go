package swp08

import "fmt"

// Validate a complete decoded message
// Returns: isValid, errorMessage
func validateMessage(data []byte) (bool, string) {
	if len(data) < 3 {
		return false, "Message too short to contain command, byte count, and checksum"
	}

	// Split fields
	command := data[0]
	byteCount := data[len(data)-2]
	checksum := data[len(data)-1]
	message := data[1 : len(data)-2]

	// Byte count check
	expectedCount := 1 + len(message) // command + message
	if int(byteCount) != expectedCount {
		return false, fmt.Sprintf("Invalid byte count: got %d, expected %d", byteCount, expectedCount)
	}

	// Checksum check
	calculated := calculateChecksum(command, message, byteCount)
	if checksum != calculated {
		return false, fmt.Sprintf("Invalid checksum: got 0x%02X, expected 0x%02X", checksum, calculated)
	}

	return true, ""
}

// Calculate SW-P-08 checksum
// Formula: 0xFF - (sum of command + message + byteCount) + 1
func calculateChecksum(command byte, message []byte, byteCount byte) byte {
	sum := int(command) + int(byteCount)
	for _, b := range message {
		sum += int(b)
	}
	lsb := byte(sum & 0xFF)
	return byte(0xFF - lsb + 1)
}
