==============================
 SW-P-08 Protocol Integration
==============================

This package can be embedded into your Go application to receive and respond to router control commands.

------------------------------
Step 1: Import the Package
------------------------------

	import "github.com/sevanjam/SW-P-08-Protocol/swp08"

Make sure your go.mod file includes the module path.

---------------------------------------------
Step 2: Implement the MatrixQuery Interface
---------------------------------------------

You must register a matrix backend that the protocol implementation can query and update.

	type MyMatrix struct{}

	func (m *MyMatrix) GetMatrixSize(matrix, level int) (sources, destinations int) {
		return 16, 16
	}

	func (m *MyMatrix) GetSourceForDestination(matrix, level, dest int) int {
		// Example logic: route source = dest
		return dest
	}

	func (m *MyMatrix) SetCrosspoint(matrix, level, dest, source int) {
		fmt.Printf("Set crosspoint: matrix=%d, level=%d, dest=%d, source=%d\n",
			matrix, level, dest, source)
	}

	func (m *MyMatrix) UseExtendedTallyDump(matrix, level int) bool {
		return false // or true if matrix > 191x191
	}

Register your matrix with:

	swp08.RegisterMatrixQuery(&MyMatrix{})

-------------------------------
Step 3: Start the TCP Server
-------------------------------

Call this from your main application to start listening for SW-P-08 control messages:

	func main() {
		// Optional: disable blocking for ACK/NAK
		swp08.ACKWaitMode = swp08.NonBlocking

		// Register matrix backend
		swp08.RegisterMatrixQuery(&MyMatrix{})

		// Start server
		go swp08.StartServer("0.0.0.0", 12345)

		// Keep the app running
		select {}
	}

-----------------------------
Logging & Debugging Notes
-----------------------------

All steps are printed to stdout using fmt.Println().
For production use, you can switch to Go’s log package or introduce a custom logger.
