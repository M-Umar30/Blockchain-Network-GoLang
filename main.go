package main

import (
	"fmt"
	"encoding/hex"
)

func main() {
	prevBlockHash, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	merkleRoot, _ := hex.DecodeString("abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789")

	// Create a block
	block := create_block([32]byte(prevBlockHash), [32]byte{}, []string{"Transaction1", "Transaction2"}, [32]byte(merkleRoot))

	// Mine the block
	success := block.mine()

	if success {
		fmt.Println("Block mined successfully!")
	} else {
		fmt.Println("Mining failed.")
	}
}
