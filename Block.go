// TODO: Figure out trailing zeroes and how to implement with block hash

package main 

type Block struct {
	Prev_Block_Hash [32]byte
	Nonce [32]byte
	Transactions []string
	Merkel_Root [32]byte
}

func create_block(prev_block_hash [32]byte, nonce [32]byte, transactions []string, merkle_root [32]byte) Block {
	return Block{prev_block_hash, nonce, transactions, merkle_root}
}