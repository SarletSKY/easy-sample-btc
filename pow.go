package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// pow算力问题
func ProofOfWork() {
	data := "luotuo"
	randomSum := 0
	level := 2
	for {
		bytes := sha256.Sum256([]byte(data + string(randomSum)))
		hash := hex.EncodeToString(bytes[:])
		if hash[:level] == "00" {
			fmt.Println(hash)
			fmt.Println(randomSum)
			break
		}
		randomSum++
	}
}
