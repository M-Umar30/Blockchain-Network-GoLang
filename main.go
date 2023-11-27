package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// test_propagation()
	test_network()
}

func test_propagation() {
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
	node1 := NewNode(1, 3) // ID, Transaction Threshold, Mining Difficulty
	node2 := NewNode(2, 3)
	node3 := NewNode(3, 3)

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
	time.Sleep(15 * time.Second)
}

func test_network() {
	// Create a P2P network
	p2pNetwork := NewP2PNetwork()

	// Create initial Peers in the P2P network
	Peer1 := NewPeer(1, 8081)
	Peer2 := NewPeer(2, 8082)
	Peer3 := NewPeer(3, 8083)

	// Register initial Peers with the bootstrap Peer
	p2pNetwork.BootstrapPeer.RegisterPeer(Peer1)
	p2pNetwork.BootstrapPeer.RegisterPeer(Peer2)
	p2pNetwork.BootstrapPeer.RegisterPeer(Peer3)

	// Start Peers
	go Peer1.Start(p2pNetwork)
	go Peer2.Start(p2pNetwork)
	go Peer3.Start(p2pNetwork)

	// Allow some time for Peers to establish connections
	time.Sleep(2 * time.Second)

	// Create a new Peer and join the network
	newPeer := NewPeer(4, 8084)
	newPeer.JoinNetwork(p2pNetwork)

	// Display the P2P network
	p2pNetwork.DisplayNetwork()

	// Wait for the simulation to finish
	time.Sleep(20 * time.Second)

	// Display connected ports for each Peer
	fmt.Printf("\nConnected Ports:\n")
	for _, Peer := range p2pNetwork.BootstrapPeer.Peers {
		connectedPorts := Peer.GetConnectedPorts()
		fmt.Printf("Peer %d connected to ports: %v\n", Peer.ID, connectedPorts)
	}
}