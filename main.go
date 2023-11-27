package main

import (
	"math/rand"
)

func main() {
	// test_propagation()
	test_propagation()
}

func test_propagation() {
	// Create a list of nodes
	nodes := make([]*Node, 0, 10)
	bootstrap_address := NewAddress(5000)

	for i := 0; i < 10; i++ {
		if i == 0 {
			node := NewNode(i+1, 4, 5000, Blockchain{}, bootstrap_address, true)
			nodes = append(nodes, node)

		} else {
			node := NewNode(i+1, 4, 5000+i, Blockchain{}, bootstrap_address, false)
			nodes = append(nodes, node)
		}
	}
	for _, node := range nodes {
		// Number of random connections per node (adjust as needed)
		numConnections := rand.Intn(3) + 1 // Generates a random number between 1 and 3

		// Select random nodes to connect
		for i := 0; i < numConnections; i++ {
			randomNode := nodes[rand.Intn(len(nodes))] // Pick a random node from the list
			if randomNode != node {
				node.AddNeighbor(randomNode)
			}
		}
	}

	// Start the servers for each node
	for _, node := range nodes {
		node.Start_Server(node.address.Port)
	}
}
