package main

type Blockchain struct {
	Blocks []Block
}

func create() Blockchain {
	return Blockchain{[]Block{}}
}

func (blockchain *Blockchain) insert_block(block Block) {
	blockchain.Blocks = append(blockchain.Blocks, block)
}

func (blockchain *Blockchain) get_length() int {
	return len(blockchain.Blocks)
}

func (blockchain *Blockchain) get_block(index int) Block {
	return blockchain.Blocks[index]
}

func (blockchain *Blockchain) get_last_block() Block {
	return blockchain.get_block(blockchain.get_length())
}