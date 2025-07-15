package swp08

import (
	"fmt"
	"net"
)

// Dual Controller Status Responce (command 0x09)
func handleDualControllerStatus(conn net.Conn) {
	fmt.Println("[SWP-08] 📡 Dual Controller Status Request (command 0x08)")

	responseMessage := []byte{0x00, 0x00} // Master active, idle OK
	response := constructResponse(0x09, responseMessage)

	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Failed to send Dual Controller Status Response:", err)
		return
	}

	fmt.Println("[SWP-08] ✅ Sent Dual Controller Status Response (command 0x09)")
	printHex(response)

	ack, err := waitForAck(conn)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Error waiting for ACK/NAK:", err)
		return
	}
	if !ack {
		fmt.Println("[SWP-08] ❌ Remote NAK received")
	}
}

// handleTallyDumpRequest handles command 0x15 — Crosspoint Tally Dump Request
func handleTallyDumpRequest(conn net.Conn, matrixByte byte) {
	fmt.Println("[SWP-08] 📥 Crosspoint Tally Dump Request (command 0x15)")

	if matrix == nil {
		fmt.Println("[SWP-08] ❌ MatrixQuery not registered")
		return
	}

	matrixID := (matrixByte & 0xF0) >> 4
	levelID := matrixByte & 0x0F

	fmt.Printf("[SWP-08] 📥 Tally Dump Request for Matrix %d, Level %d\n", int(matrixID), int(levelID))

	numSources, numDestinations := matrix.GetMatrixSize(int(matrixID), int(levelID))
	useExtended := false

	if m, ok := matrix.(interface {
		UseExtendedTallyDump(matrix, level int) bool
	}); ok {
		useExtended = m.UseExtendedTallyDump(int(matrixID), int(levelID))
	} else {
		useExtended = numSources > 191 || numDestinations > 191
	}

	var messages [][]byte

	if useExtended {
		messages = buildTallyDumpWordMessages(matrixByte, int(matrixID), int(levelID), numDestinations)
	} else {
		messages = buildTallyDumpByteMessages(matrixByte, int(matrixID), int(levelID), numDestinations)
	}

	for i, msg := range messages {
		cmd := byte(0x16)
		if useExtended {
			cmd = 0x17
		}

		response := constructResponse(cmd, msg)
		_, err := conn.Write(response)
		if err != nil {
			fmt.Printf("[SWP-08] ❌ Failed to send chunk %d: %v\n", i, err)
			return
		}

		fmt.Printf("[SWP-08] ✅ Sent Crosspoint Tally Dump chunk %d (command 0x%02X)\n", i, cmd)
		printHex(response)

		ack, err := waitForAck(conn)
		if err != nil {
			fmt.Println("[SWP-08] ❌ Error waiting for ACK/NAK:", err)
			return
		}
		if !ack {
			fmt.Println("[SWP-08] ❌ Remote NAK received")
			return
		}
	}
}
func buildTallyDumpByteMessages(matrixByte byte, matrixID, level, numDest int) [][]byte {
	const maxSourcesPerMessage = 125
	var messages [][]byte

	for i := 0; i < numDest; i += maxSourcesPerMessage {
		count := min(maxSourcesPerMessage, numDest-i)
		msg := []byte{matrixByte, byte(count), byte(i)}

		for j := 0; j < count; j++ {
			source := matrix.GetSourceForDestination(matrixID, level, i+j)
			msg = append(msg, byte(source&0xFF))
		}
		messages = append(messages, msg)
	}

	return messages
}
func buildTallyDumpWordMessages(matrixByte byte, matrixID, level, numDest int) [][]byte {
	const maxTallies = 64
	var messages [][]byte

	for i := 0; i < numDest; i += maxTallies {
		count := min(maxTallies, numDest-i)
		destHi := byte((i >> 8) & 0xFF)
		destLo := byte(i & 0xFF)

		msg := []byte{matrixByte, byte(count), destHi, destLo}

		for j := 0; j < count; j++ {
			source := matrix.GetSourceForDestination(matrixID, level, i+j)
			srcHi := byte((source >> 8) & 0xFF)
			srcLo := byte(source & 0xFF)
			msg = append(msg, srcHi, srcLo)
		}
		messages = append(messages, msg)
	}

	return messages
}

// Set Crosspoint request
func handleSetCrosspoint(conn net.Conn, message []byte) {
	if len(message) < 4 {
		fmt.Println("[SWP-08] ❌ Invalid Crosspoint Connect message length")
		sendNegativeAcknowledge(conn)
		return
	}

	matrixByte := message[0]
	matrixNumber := (matrixByte & 0xF0) >> 4
	levelNumber := matrixByte & 0x0F

	multiplierByte := message[1]
	destMultiplier := (multiplierByte & 0x70) >> 4 // bits 4–6
	srcMultiplier := multiplierByte & 0x07         // bits 0–2

	dest := int(message[2]) + int(destMultiplier)*128
	src := int(message[3]) + int(srcMultiplier)*128

	fmt.Printf("[SWP-08] 🔀 Crosspoint Connect Request: Matrix %d Level %d | Destination %d → Source %d\n",
		matrixNumber, levelNumber, dest, src)

	if matrix == nil {
		fmt.Println("[SWP-08] ❌ No matrix registered")
		sendNegativeAcknowledge(conn)
		return
	}

	ok := false

	if m, okCast := matrix.(interface {
		SetSourceForDestination(matrix, level, destination, source int) bool
	}); okCast {
		ok = m.SetSourceForDestination(int(matrixNumber), int(levelNumber), dest, src)
	} else {
		fmt.Println("[SWP-08] ❌ Matrix does not support SetSourceForDestination")
	}

	if ok {
		// sendAcknowledge(conn)
		sendCrosspointConnected(conn, matrixByte, multiplierByte, message[2], message[3])
	} else {
		sendNegativeAcknowledge(conn)
	}
}
func sendCrosspointConnected(conn net.Conn, matrixByte, multiplierByte, destByte, sourceByte byte) {
	response := []byte{matrixByte, multiplierByte, destByte, sourceByte}
	full := constructResponse(0x04, response)

	_, err := conn.Write(full)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Failed to send Crosspoint Connected response:", err)
		return
	}

	fmt.Println("[SWP-08] ✅ Sent Crosspoint Connected response (command 0x04)")
	printHex(full)

	ack, err := waitForAck(conn)
	if err != nil {
		fmt.Println("[SWP-08] ❌ Error waiting for ACK/NAK after 0x04:", err)
		return
	}
	if !ack {
		fmt.Println("[SWP-08] ❌ Remote NAK received for Crosspoint Connected")
	}
}
