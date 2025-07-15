package main

import (
	// Replace your-app/swp08 with the correct import path after you run go mod init.
	"swp08protocol/swp08"
)

func main() {
	// Replace with any IP and port you want
	swp08.RegisterMatrixQuery(swp08.NewMockMatrix(300, 300))
	swp08.StartServer("0.0.0.0", 12345)
}
