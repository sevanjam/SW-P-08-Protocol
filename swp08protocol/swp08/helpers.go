package swp08

var ACKWaitMode = Blocking

type ackMode bool

const (
	Blocking    ackMode = true
	NonBlocking ackMode = false
)

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
