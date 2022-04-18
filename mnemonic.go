package main

import (
	// "strconv"
	"strings"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"bufio"
)

func ByteArrayToBitString(b []byte) string {
	var sb strings.Builder
	for _, n := range(b) {
		binary := fmt.Sprintf("%08b", n) // prints 00000000 11111101
		sb.WriteString(binary)
    	}
	return sb.String()
}

func AddChecksum(entropy string, sha256 string) string {
	var sb strings.Builder
	sb.WriteString(entropy)
	for i := 0; i < len(entropy)/32; i++ {
		fmt.Println(string(entropy[i]))
		sb.WriteString(string(entropy[i]))
    	}
	fmt.Println("checksum: ", sb.String())
	return sb.String()
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
	    return nil
	}
	if chunkSize >= len(s) {
	    return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
	    if currentLen == chunkSize {
		chunks = append(chunks, s[currentStart:i])
		currentLen = 0
		currentStart = i
	    }
	    currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
	    return nil, err
	}
	defer file.Close()
    
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func BitStringToNumber(bitString string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
	    return nil, err
	}
	defer file.Close()
    
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}



func main() {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(b)
	entropy := ByteArrayToBitString(b)
	fmt.Println()
	fmt.Println("entropy: ", entropy, len(entropy))

	// output, err := strconv.ParseInt(entropy[:11], 2, 64)  
	// if err != nil {  
	// fmt.Println(err)  
	// return  
	// }  
	
	// fmt.Println()
	// fmt.Println(entropy[:11])
	// fmt.Printf("Output %d\n", output)  

	hash := sha256.Sum256(b)
	hashString := ByteArrayToBitString(hash[:])
	fmt.Println("sha256 hash: ", hashString, len(hashString))
	withChecksum := AddChecksum(entropy, hashString)
	fmt.Println("entropy + checksum: ", withChecksum, len(withChecksum))
	fmt.Println("chunks: ", Chunks(withChecksum, 11))
	wordArray, err := ReadLines("wordlist.txt")
	fmt.Println("word array: ", len(wordArray))
}