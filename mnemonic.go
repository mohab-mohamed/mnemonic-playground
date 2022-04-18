package main

import (
	"strconv"
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

	output, err := strconv.ParseInt(entropy[:11], 2, 64)  
	if err != nil {  
	fmt.Println(err)  
	return  
	}  
	
	fmt.Println()
	fmt.Println(entropy[:11])
	fmt.Printf("Output %d\n", output)  

}