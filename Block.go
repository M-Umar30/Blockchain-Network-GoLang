package main

import "fmt"

//"crypto/rand"

var trailing_zeros = 8
type Block struct {
	Prev_Block_Hash [32]byte
	Nonce [32]byte
	Transactions []string
	Merkel_Root *InternalNode
}

func create_block(prev_block_hash [32]byte, nonce [32]byte, transactions []string) Block {
	//calling the create merkle tree function
	merkle_root := CreateMerkleTree(transactions)
	return Block{prev_block_hash, nonce, transactions, merkle_root}
}

func create_merkle_tree(transactions []string) {
	panic("unimplemented")
}

func (block *Block) mine() bool {
	fmt.Println("Mining...")
	for {
		fmt.Println("Trying...")
		block.Nonce = generate_nonce()
		
		// Access the hash value within the InternalNode and convert it to a byte slice
        merkleRootHash := []byte(block.Merkel_Root.Self_Hash)


		//output := concatenate_hashes(block.Prev_Block_Hash[:], block.Nonce[:], block.Merkel_Root[:])
		output := concatenate_hashes(block.Prev_Block_Hash[:], block.Nonce[:], merkleRootHash[:])		
		// check trailing zeros
		if count_trailing_zeros(output) >= trailing_zeros {
			fmt.Println("Success")
			// print nonce
			fmt.Println("Nonce: ", hash_to_string(block.Nonce))
			return true
		}
	}
}