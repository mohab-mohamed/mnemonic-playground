package main

import (
	"strings"
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
	var sb strings.Builder
	for i, n := range(b) {
		fmt.Println(i)
		binary := fmt.Sprintf("%b", n) // prints 00000000 11111101
		fmt.Println(binary)
		fmt.Printf("%d\n\n", n)
		sb.WriteString(binary)
    	}
	entropy := sb.String()
	fmt.Println()
	fmt.Println(entropy)

}