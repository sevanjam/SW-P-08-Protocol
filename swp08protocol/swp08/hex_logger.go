package swp08

import (
	"fmt"
)

func printHex(data []byte) {
	for i, b := range data {
		fmt.Printf("%02X ", b)
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	if len(data)%16 != 0 {
		fmt.Println()
	}
}
