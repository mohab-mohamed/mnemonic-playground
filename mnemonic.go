package main

import (
	"strconv"
	"strings"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/hmac"
	"hash"
	"fmt"
	"os"
	"bufio"
	"encoding/hex"
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

// helper function used to split bit string into 11 unit chunks for mnemonic indexing
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
func BitStringToNumber(bitString string) (int64, error) {
	output, err := strconv.ParseInt(bitString, 2, 64)  
	if err != nil {  
		fmt.Println(err)  
		return -1, err 
	} 
	return output, err 
}

func GetMnemonic(chunks []string, wordList []string) []string {
	var sentence []string
	for _, chunk := range(chunks) {
		num, _ := BitStringToNumber(chunk)
		sentence = append(sentence, wordList[num])
	}
	return sentence
}

// Key derives a key from the password, salt and iteration count, returning a
// []byte of length keylen that can be used as cryptographic key. The key is
// derived based on the method described as PBKDF2 with the HMAC variant using
// the supplied hash function.
//
// For example, to use a HMAC-SHA-1 based PBKDF2 key derivation function, you
// can get a derived key for e.g. AES-256 (which needs a 32-byte key) by
// doing:
//
//	dk := pbkdf2.Key([]byte("some password"), salt, 4096, 32, sha1.New)
//
// Remember to get a good random salt. At least 8 bytes is recommended by the
// RFC.
//
// Using a higher iteration count will increase the cost of an exhaustive
// search but will also make derivation proportionally slower.
func Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
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
	
	// fmt.Println()
	// fmt.Println(entropy[:11])
	// fmt.Printf("Output %d\n", output)  

	hash := sha256.Sum256(b)
	hashString := ByteArrayToBitString(hash[:])
	fmt.Println("sha256 hash: ", hashString, len(hashString))
	withChecksum := AddChecksum(entropy, hashString)
	fmt.Println("entropy + checksum: ", withChecksum, len(withChecksum))
	chunks := Chunks(withChecksum, 11)
	fmt.Println("chunks: ", chunks)
	wordList, err := ReadLines("wordlist.txt")
	fmt.Println("word array: ", len(wordList))
	mnemonicList := GetMnemonic(chunks, wordList)
	fmt.Println("mnemonic list: ", mnemonicList, "word count: ", len(mnemonicList))
	mnemonicSentence := strings.Join(mnemonicList, " ")
	fmt.Println("mnemonic sentence: ", mnemonicSentence)
	salt := []byte("mnemonic")
	seed := Key([]byte(mnemonicSentence), salt, 4096, 64, sha512.New)
	fmt.Println("seed: ", seed, "# bytes: ", len(seed))
	hexSeed := hex.EncodeToString(seed)
	fmt.Println("hex seed: ", hexSeed, "length: ", len(hexSeed))
}