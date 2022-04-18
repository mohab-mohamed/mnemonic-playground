package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

func main() {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(rand.Read(b))
	// The slice should now contain random bytes instead of only zeroes.
	fmt.Println(b)
	fmt.Println(bytes.Equal(b, make([]byte, c)))

}