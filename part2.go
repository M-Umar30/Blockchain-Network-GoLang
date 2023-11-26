package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// type Block struct {
// 	Transactions []string
// 	Nonce        int
// 	PreviousHash string
// 	SelfHash     string
// }

type Node struct {
	ID                     int
	recentTransactions     map[string]bool
	transactionChannel     chan string
	neighborChannels       []chan string
	transactionMutex       sync.Mutex
	recentTransactionsMu   sync.Mutex
	blockChannel           chan []byte
	minedBlockChannel      chan Block
	miningInProgress       bool
	miningInProgressMu     sync.Mutex
	transactionThreshold   int
	miningDifficulty       int
	miningDifficultyString string
	block                  Block
	Neighbors              []*Node
	NeighborChans          []chan []byte // Each neighbor has a dedicated channel
	neighborChannel        chan []byte   // Channel for general communication with neighbors
	recentlySentBlocks     map[int]bool
	lastMinedSenderID      int
}

func NewNode(id, transactionThreshold, miningDifficulty int) *Node {
	return &Node{
		ID:                   id,
		recentTransactions:   make(map[string]bool),
		transactionChannel:   make(chan string),
		blockChannel:         make(chan []byte),
		minedBlockChannel:    make(chan Block),
		transactionThreshold: transactionThreshold,
		miningDifficulty:     miningDifficulty,
		recentlySentBlocks:   make(map[int]bool),
	}
}

func (node *Node) AddNeighborChannel(neighborChannel chan string) {
	node.neighborChannels = append(node.neighborChannels, neighborChannel)
	//node.blockChannel = append(node.blockChannel, neighborChannel)
}
func (node *Node) AddNeighbor(neighbor *Node) {
	// Inside NewNode function
	node.NeighborChans = append(node.NeighborChans, make(chan []byte))

}

// func (node *Node) AddNeighborChannel2(neighborChannel chan []byte) {
// 	node.neighborChannels = append(node.neighborChannels, neighborChannel)
// }

func (node *Node) ReceiveTransaction(transaction string) {
	node.transactionMutex.Lock()
	defer node.transactionMutex.Unlock()

	// Check if the transaction is already in the recent transactions list
	if _, exists := node.recentTransactions[transaction]; !exists {
		// Add the transaction to the recent transactions list
		node.recentTransactions[transaction] = true

		// Broadcast the transaction to neighbors
		go node.BroadcastTransaction(transaction)
	}
}

func (node *Node) BroadcastTransaction(transaction string) {
	for _, neighborChannel := range node.neighborChannels {
		// Send the transaction to all neighbors
		neighborChannel <- transaction
	}
}

func (node *Node) MineBlock(transactions []string, previousHash string) Block {
	prevhash, err := hex.DecodeString(previousHash)
	if err != nil {
		fmt.Println("Error decoding prevhash:", err)
	}

	root := CreateMerkleTree(transactions)
	merkleroot, err := hex.DecodeString(root.Self_Hash)
	if err != nil {
		fmt.Println("Error decoding merkleroot:", err)
	}

	// Create [32]byte from the byte slices
	var prevHashArray [32]byte
	copy(prevHashArray[:], prevhash)

	var merkleRootArray [32]byte
	copy(merkleRootArray[:], merkleroot)

	node.block = create_block(prevHashArray, [32]byte{}, transactions, merkleRootArray)

	// Mine the block
	node.block.mine()

	return node.block

}

// func calculateBlockHash(block Block) string {
// 	// This is a simplified hash calculation for demonstration purposes.
// 	return fmt.Sprintf("%d-%s-%d", block.Nonce, block.Prev_Block_Hash, len(block.Transactions))
// }

// func isHashValid(hash string, difficulty int) bool {
// 	// Check if the hash starts with a certain number of zeros based on difficulty
// 	return hash[:difficulty] == node.miningDifficultyString
// }

func (node *Node) Start() {
	// Create a WaitGroup
	var wg sync.WaitGroup

	// Increment the WaitGroup counter for the handleIncomingBlocks goroutine
	wg.Add(1)

	// Start handling incoming blocks concurrently
	go node.handleIncomingBlocks(&wg)

	// Simulate receiving transactions over time
	for i := 1; i <= 8; i++ {
		transaction := fmt.Sprintf("Transaction%d", i)
		time.Sleep(time.Second)
		fmt.Printf("Node %d received transaction: %s\n", node.ID, transaction)
		node.ReceiveTransaction(transaction)

		// If the transaction threshold is reached, start mining a block
		if i%node.transactionThreshold == 0 {
			go node.MineAndBroadcastBlock()
		}
	}

	// Wait for a short duration for the handleIncomingBlocks goroutine to execute
	time.Sleep(5 * time.Second)

	// Signal that the program can exit by decrementing the WaitGroup counter
	wg.Done()

	// Wait for the handleIncomingBlocks goroutine to finish
	wg.Wait()
}

func (node *Node) MineAndBroadcastBlock() {
	// Ensure only one mining process is in progress at a time
	node.miningInProgressMu.Lock()
	defer node.miningInProgressMu.Unlock()

	// Check if mining is already in progress
	if node.miningInProgress {
		return
	}

	// Set mining in progress flag
	node.miningInProgress = true

	// Get recent transactions for mining
	node.recentTransactionsMu.Lock()
	transactions := make([]string, 0, len(node.recentTransactions))
	for transaction := range node.recentTransactions {
		transactions = append(transactions, transaction)
	}
	node.recentTransactionsMu.Unlock()

	// Simulate mining process
	minedBlock := node.MineBlock(transactions, "PreviousHash")
	node.lastMinedSenderID = node.ID

	// Reset recent transactions after including them in a block
	node.recentTransactionsMu.Lock()
	node.recentTransactions = make(map[string]bool)
	node.recentTransactionsMu.Unlock()

	// Broadcast the mined block to neighbors
	fmt.Printf("Node %d sending mined block: %s\n", node.ID, hex.EncodeToString(minedBlock.Self_Hash[:]))
	go node.BroadcastMinedBlock(minedBlock)

	// Reset mining in progress flag
	node.miningInProgress = false
}

// // Function to broadcast a mined block to all neighbors
// func (node *Node) BroadcastMinedBlock(minedBlock Block) {
// 	for _, neighborChan := range node.NeighborChans {
// 		// Send the mined block to each neighbor's channel
// 		neighborChan <- minedBlock
// 	}

// 	// You can also broadcast to the main block channel if needed
// 	node.blockChannel <- minedBlock
// }

func (node *Node) BroadcastMinedBlock(minedBlock Block) {
	// Encode the mined block
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(minedBlock); err != nil {
		fmt.Printf("Error encoding mined block: %v\n", err)
		return
	}
	fmt.Println("Entering loop")

	// Broadcast the encoded block to neighbors
	fmt.Printf("Last mined sender ID: %d\n", node.lastMinedSenderID)
	fmt.Printf("Number of NeighborChans for node %d: %d\n", node.ID, len(node.NeighborChans))

	for i, neighborChan := range node.NeighborChans {
		//if i != node.lastMinedSenderID {
		neighborChan <- buffer.Bytes()
		fmt.Printf("Node %d sent block to neighbor %d\n", node.ID, i)
		//}
	}

	// Broadcast the encoded block to the main block channel
	node.blockChannel <- buffer.Bytes()
	fmt.Printf("Node %d sent block to main channel\n", node.ID)
}

// func (node *Node) BroadcastMinedBlock(minedBlock Block) {
// 	for _, neighborChannel := range node.NeighborChans {
// 		// Create a buffer to encode the mined block
// 		var buffer bytes.Buffer
// 		encoder := gob.NewEncoder(&buffer)
// 		if err := encoder.Encode(minedBlock); err != nil {
// 			fmt.Printf("Error encoding mined block: %v\n", err)
// 			continue
// 		}

// 		// Send the encoded block to all neighbors
// 		neighborChannel <- buffer.Bytes()
// 	}

// 	// Create a buffer to encode the mined block
// 	var mainBuffer bytes.Buffer
// 	mainEncoder := gob.NewEncoder(&mainBuffer)
// 	if err := mainEncoder.Encode(minedBlock); err != nil {
// 		fmt.Printf("Error encoding mined block: %v\n", err)
// 		return
// 	}

// 	// Send the encoded block to the main block channel
// 	node.blockChannel <- mainBuffer.Bytes()
// }

func (node *Node) ReceiveMinedBlock(minedBlock Block) {
	// Process the received mined block
	fmt.Printf("Node %d received mined block: %s\n", node.ID, minedBlock.Self_Hash)

	// TODO: Further processing of the mined block
}

func (node *Node) handleIncomingBlocks(wg *sync.WaitGroup) {
	fmt.Printf("Node %d started handling incoming blocks\n", node.ID)
	defer wg.Done()

	for {
		// Read from each neighbor's channel
		for i, neighborChan := range node.NeighborChans {
			select {
			case encodedBlock := <-neighborChan:
				// Decode the received mined block from a neighbor's channel
				var receivedBlock Block
				decoder := gob.NewDecoder(bytes.NewReader(encodedBlock))
				if err := decoder.Decode(&receivedBlock); err != nil {
					fmt.Printf("Error decoding mined block: %v\n", err)
					continue
				}

				// Process the received mined block from a neighbor's channel
				fmt.Printf("Node %d received mined block from neighbor %d: %x\n", node.ID, i, receivedBlock.Self_Hash)
				// Add more cases if you have other types of messages or channels
			}
		}

		// Read from the main block channel
		// select {
		// case encodedBlock := <-node.blockChannel:
		// 	// Decode the received mined block from the main channel
		// 	var receivedBlock Block
		// 	mainDecoder := gob.NewDecoder(bytes.NewReader(encodedBlock))
		// 	if err := mainDecoder.Decode(&receivedBlock); err != nil {
		// 		fmt.Printf("Error decoding mined block: %v\n", err)
		// 		continue
		// 	}

		// 	// Process the received mined block from the main channel
		// 	fmt.Printf("Node %d received mined block from main channel: %x\n", node.ID, receivedBlock.Self_Hash)
		// 	// Add more cases if you have other types of messages or channels
		// }
	}
}

// func (node *Node) handleIncomingBlocks(wg *sync.WaitGroup) {
// 	fmt.Printf("Node %d started handling incoming blocks\n", node.ID)
// 	defer wg.Done()

// 	for {
// 		// Read from each neighbor's channel
// 		for i, neighborChan := range node.NeighborChans {
// 			select {
// 			case minedBlock := <-neighborChan:
// 				// Process the received mined block from a neighbor's channel
// 				fmt.Printf("Node %d received mined block from neighbor: %s\n", node.ID, minedBlock.Self_Hash)
// 				// Add more cases if you have other types of messages or channels
// 			}
// 		}
// 	}
// }

// func (node *Node) handleIncomingBlocks() {
// 	for minedBlock := range node.blockChannel {
// 		// Process the received mined block
// 		node.ReceiveMinedBlock(minedBlock)
// 	}
// }

// func (node *Node) ReceiveMinedBlock(minedBlock Block) {
// 	// Process the received mined block
// 	fmt.Printf("Node %d received mined block: %s\n", node.ID, minedBlock.Self_Hash)

// 	// TODO: Further processing of the mined block
// }