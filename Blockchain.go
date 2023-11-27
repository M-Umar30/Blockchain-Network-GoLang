package main

type Blockchain struct {
	Blocks []Block
}

func create() Blockchain {
	return Blockchain{[]Block{}}
}

func (blockchain *Blockchain) insert_block(block Block) bool{
	if block.Prev_Block_Hash == blockchain.get_last_block().Self_Hash {
		blockchain.Blocks = append(blockchain.Blocks, block)
		return true
	}
	return false
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

func (blockchain *Blockchain) search_block(block Block) bool {
	for i := 0; i < blockchain.get_length(); i++ {
		if blockchain.get_block(i).Self_Hash == block.Self_Hash {
			return true
		}
	}
	return false
}