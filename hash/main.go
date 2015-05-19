package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	for _, input := range(os.Args[1:]) {
		fmt.Println(hash(input))
	}
}

func hash(input string) string {
	h := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", h)
}
