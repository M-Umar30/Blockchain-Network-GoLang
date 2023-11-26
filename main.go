package main

import (
	"math/rand"
	"time"
)

func main() {
	// prevBlockHash, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	// merkleRoot, _ := hex.DecodeString("abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789")

	// // Create a block
	// block := create_block([32]byte(prevBlockHash), [32]byte{}, []string{"Transaction1", "Transaction2"}, [32]byte(merkleRoot))

	// // Mine the block
	// success := block.mine()

	// if success {
	// 	fmt.Println("Block mined successfully!")
	// } else {
	// 	fmt.Println("Mining failed.")
	// }

	rand.Seed(time.Now().UnixNano())

	// Create nodes
	node1 := NewNode(1, 3, 2) // ID, Transaction Threshold, Mining Difficulty
	node2 := NewNode(2, 3, 2)
	node3 := NewNode(3, 3, 2)

	// Connect nodes with neighbor channels
	node1.AddNeighborChannel(node2.transactionChannel)
	node1.AddNeighborChannel(node3.transactionChannel)
	node1.AddNeighbor(node2)
	node1.AddNeighbor(node3)

	node2.AddNeighborChannel(node1.transactionChannel)
	//node2.AddNeighborChannel(node3.transactionChannel)
	node2.AddNeighbor(node1)
	//node2.AddNeighbor(node3)

	node3.AddNeighborChannel(node1.transactionChannel)
	//node3.AddNeighborChannel(node2.transactionChannel)
	//node3.AddNeighbor(node2)
	node3.AddNeighbor(node1)

	// Start nodes
	go node1.Start()
	go node2.Start()
	go node3.Start()
	// Wait for the simulation to finish
	time.Sleep(10 * time.Second)
}
