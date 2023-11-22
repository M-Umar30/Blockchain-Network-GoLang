package main

import "fmt"

//"crypto/rand"

var trailing_zeros = 8
type Block struct {
	Prev_Block_Hash [32]byte
	Nonce [32]byte
	Transactions []string
	Merkel_Root [32]byte
}

func create_block(prev_block_hash [32]byte, nonce [32]byte, transactions []string, merkle_root [32]byte) Block {
	return Block{prev_block_hash, nonce, transactions, merkle_root}
}

func (block *Block) mine() bool {
	fmt.Println("Mining...")
	for {
		fmt.Println("Trying...")
		block.Nonce = generate_nonce()
		output := concatenate_hashes(block.Prev_Block_Hash[:], block.Nonce[:], block.Merkel_Root[:])		
		// check trailing zeros
		if count_trailing_zeros(output) >= trailing_zeros {
			fmt.Println("Success")
			// print nonce
			fmt.Println("Nonce: ", hash_to_string(block.Nonce))
			return true
		}
	}
}