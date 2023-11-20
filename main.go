/*
	20I-0518 - Muhammad Umar
	20I-0450 - Fatima Zubeda
	20L-
	20I-
*/

package main

import (
	"fmt"
	"encoding/hex"
)

// TODO: Verifying block and chain (shouldn't it be like a given that its verified cuz implementation?) just make a function which iterates through everything maybe?
func main() {
	// Example data
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