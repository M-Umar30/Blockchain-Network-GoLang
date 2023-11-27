package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Node represents a peer in the P2P network.
type Node struct {
	ID        int
	Port      int
	IP        string
	Neighbors []*Node
	mu        sync.RWMutex
}

// BootstrapNode represents the bootstrap node that keeps track of all nodes in the network.
type BootstrapNode struct {
	Nodes map[int]*Node
	mu    sync.RWMutex
}

// P2PNetwork represents the P2P network.
type P2PNetwork struct {
	BootstrapNode *BootstrapNode
}

func main() {
	// Create a P2P network
	p2pNetwork := NewP2PNetwork()

	// Create initial nodes in the P2P network
	node1 := NewNode(1, 8081)
	node2 := NewNode(2, 8082)
	node3 := NewNode(3, 8083)

	// Register initial nodes with the bootstrap node
	p2pNetwork.BootstrapNode.RegisterNode(node1)
	p2pNetwork.BootstrapNode.RegisterNode(node2)
	p2pNetwork.BootstrapNode.RegisterNode(node3)

	// Start nodes
	go node1.Start(p2pNetwork)
	go node2.Start(p2pNetwork)
	go node3.Start(p2pNetwork)

	// Allow some time for nodes to establish connections
	time.Sleep(2 * time.Second)

	// Create a new node and join the network
	newNode := NewNode(4, 8084)
	newNode.JoinNetwork(p2pNetwork)

	// Display the P2P network
	p2pNetwork.DisplayNetwork()

	// Wait for the simulation to finish
	time.Sleep(20 * time.Second)

	// Display connected ports for each node
	fmt.Printf("\nConnected Ports:\n")
	for _, node := range p2pNetwork.BootstrapNode.Nodes {
		connectedPorts := node.GetConnectedPorts()
		fmt.Printf("Node %d connected to ports: %v\n", node.ID, connectedPorts)
	}
}

// NewNode creates a new Node instance.
func NewNode(id, port int) *Node {
	return &Node{
		ID:   id,
		Port: port,
		IP:   "localhost",
	}
}

// Start starts the server and client functionalities of the node.
func (n *Node) Start(p2pNetwork *P2PNetwork) {
	// Start server
	go n.StartServer()

	// Start client to communicate with bootstrap node
	go n.StartClient(p2pNetwork)
}

// StartServer starts the server functionality of the node.
func (n *Node) StartServer() {
	listenAddr := fmt.Sprintf("%s:%d", n.IP, n.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Node %d: Error starting server: %v\n", n.ID, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Node %d: Server started at %s\n", n.ID, listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Node %d: Error accepting connection: %v\n", n.ID, err)
			continue
		}

		go n.HandleClient(conn)
	}
}

// StartClient starts the client functionality of the node to communicate with the bootstrap node.
func (n *Node) StartClient(p2pNetwork *P2PNetwork) {
	for {
		// Simulate periodic request to bootstrap node for neighbors
		time.Sleep(5 * time.Second)

		// Contact bootstrap node to get neighbors
		neighbors := p2pNetwork.BootstrapNode.GetNeighbors(n.ID)
		if len(neighbors) > 0 {
			fmt.Printf("Node %d: Received neighbors from bootstrap: %v\n", n.ID, neighbors)
			// Establish connections with neighbors
			n.EstablishConnections(neighbors)
		}
	}
}

// JoinNetwork joins the P2P network by contacting the bootstrap node.
func (n *Node) JoinNetwork(p2pNetwork *P2PNetwork) {
	// Start server to accept incoming connections
	go n.StartServer()

	// Contact bootstrap node to get initial neighbors
	neighbors := p2pNetwork.BootstrapNode.GetNeighbors(n.ID)
	if len(neighbors) > 0 {
		fmt.Printf("Node %d: Joined the network with neighbors: %v\n", n.ID, neighbors)
		// Establish connections with neighbors
		n.EstablishConnections(neighbors)
		// Register with bootstrap node
		p2pNetwork.BootstrapNode.RegisterNode(n)
	}
}

// EstablishConnections establishes connections with the specified neighbors.
func (n *Node) EstablishConnections(neighbors []*Node) {
	for _, neighbor := range neighbors {
		if n.ID != neighbor.ID {
			go n.ConnectToNeighbor(neighbor)
		}
	}
}

// ConnectToNeighbor establishes a connection with a neighbor.
func (n *Node) ConnectToNeighbor(neighbor *Node) {
	addr := fmt.Sprintf("%s:%d", neighbor.IP, neighbor.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Node %d: Error connecting to neighbor %d: %v\n", n.ID, neighbor.ID, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Node %d: Connected to neighbor %d at %s\n", n.ID, neighbor.ID, addr)

	// Perform additional communication with the neighbor as needed
	// For simplicity, we're not sending any data in this example.
}

// HandleClient handles incoming connections from other nodes.
func (n *Node) HandleClient(conn net.Conn) {
	defer conn.Close()

	// Perform communication with the connected node as needed
	// For simplicity, we're not receiving any data in this example.
}

// NewBootstrapNode creates a new BootstrapNode instance.
func NewBootstrapNode() *BootstrapNode {
	return &BootstrapNode{
		Nodes: make(map[int]*Node),
	}
}

// RegisterNode registers a node with the bootstrap node.
func (b *BootstrapNode) RegisterNode(node *Node) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Nodes[node.ID] = node
	fmt.Printf("Node %d registered with BootstrapNode\n", node.ID)
}

// GetNeighbors returns the neighbors of a node based on its ID.
func (b *BootstrapNode) GetNeighbors(nodeID int) []*Node {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var neighbors []*Node

	for id, node := range b.Nodes {
		if id != nodeID {
			neighbors = append(neighbors, node)
		}
	}

	return neighbors
}

// NewP2PNetwork creates a new P2PNetwork instance.
func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		BootstrapNode: NewBootstrapNode(),
	}
}

// DisplayNetwork displays all the nodes in the P2P network and their connections.
func (p *P2PNetwork) DisplayNetwork() {
	p.BootstrapNode.mu.RLock()
	defer p.BootstrapNode.mu.RUnlock()

	fmt.Println("P2P Network:")
	for _, node := range p.BootstrapNode.Nodes {
		fmt.Printf("Node %d connected to: ", node.ID)
		for _, neighbor := range node.Neighbors {
			fmt.Printf("%d ", neighbor.ID)
		}
		fmt.Println()
	}
}

// GetConnectedPorts returns the port numbers of all nodes connected to the server.
func (n *Node) GetConnectedPorts() []int {
	// Lock neighbors for reading
	n.mu.RLock()
	defer n.mu.RUnlock()

	var ports []int

	for _, neighbor := range n.Neighbors {
		ports = append(ports, neighbor.Port)
	}

	return ports
}
