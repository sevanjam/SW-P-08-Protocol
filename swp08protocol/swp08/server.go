package swp08

import (
	"fmt"
	"net"
)

func StartServer(ip string, port int) {
	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println("[SWP-08] Starting TCP server on", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("[SWP-08] Failed to start server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[SWP-08] Failed to accept connection:", err)
			continue
		}
		fmt.Println("[SWP-08] Accepted new connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
