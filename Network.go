package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Peer represents a peer in the P2P network.
type Peer struct {
	ID        int
	Port      int
	IP        string
	Neighbors []*Peer
	mu        sync.RWMutex
}

// BootstrapPeer represents the bootstrap Peer that keeps track of all Peers in the network.
type BootstrapPeer struct {
	Peers map[int]*Peer
	mu    sync.RWMutex
}

// P2PNetwork represents the P2P network.
type P2PNetwork struct {
	BootstrapPeer *BootstrapPeer
}



// NewPeer creates a new Peer instance.
func NewPeer(id, port int) *Peer {
	return &Peer{
		ID:   id,
		Port: port,
		IP:   "localhost",
	}
}

// Start starts the server and client functionalities of the Peer.
func (n *Peer) Start(p2pNetwork *P2PNetwork) {
	// Start server
	go n.StartServer()

	// Start client to communicate with bootstrap Peer
	go n.StartClient(p2pNetwork)
}

// StartServer starts the server functionality of the Peer.
func (n *Peer) StartServer() {
	listenAddr := fmt.Sprintf("%s:%d", n.IP, n.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Peer %d: Error starting server: %v\n", n.ID, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Peer %d: Server started at %s\n", n.ID, listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Peer %d: Error accepting connection: %v\n", n.ID, err)
			continue
		}

		go n.HandleClient(conn)
	}
}

// StartClient starts the client functionality of the Peer to communicate with the bootstrap Peer.
func (n *Peer) StartClient(p2pNetwork *P2PNetwork) {
	for {
		// Simulate periodic request to bootstrap Peer for neighbors
		time.Sleep(5 * time.Second)

		// Contact bootstrap Peer to get neighbors
		neighbors := p2pNetwork.BootstrapPeer.GetNeighbors(n.ID)
		if len(neighbors) > 0 {
			fmt.Printf("Peer %d: Received neighbors from bootstrap: %v\n", n.ID, neighbors)
			// Establish connections with neighbors
			n.EstablishConnections(neighbors)
		}
	}
}

// JoinNetwork joins the P2P network by contacting the bootstrap Peer.
func (n *Peer) JoinNetwork(p2pNetwork *P2PNetwork) {
	// Start server to accept incoming connections
	go n.StartServer()

	// Contact bootstrap Peer to get initial neighbors
	neighbors := p2pNetwork.BootstrapPeer.GetNeighbors(n.ID)
	if len(neighbors) > 0 {
		fmt.Printf("Peer %d: Joined the network with neighbors: %v\n", n.ID, neighbors)
		// Establish connections with neighbors
		n.EstablishConnections(neighbors)
		// Register with bootstrap Peer
		p2pNetwork.BootstrapPeer.RegisterPeer(n)
	}
}

// EstablishConnections establishes connections with the specified neighbors.
func (n *Peer) EstablishConnections(neighbors []*Peer) {
	for _, neighbor := range neighbors {
		if n.ID != neighbor.ID {
			go n.ConnectToNeighbor(neighbor)
		}
	}
}

// ConnectToNeighbor establishes a connection with a neighbor.
func (n *Peer) ConnectToNeighbor(neighbor *Peer) {
	addr := fmt.Sprintf("%s:%d", neighbor.IP, neighbor.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Peer %d: Error connecting to neighbor %d: %v\n", n.ID, neighbor.ID, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Peer %d: Connected to neighbor %d at %s\n", n.ID, neighbor.ID, addr)

	// Perform additional communication with the neighbor as needed
	// For simplicity, we're not sending any data in this example.
}

// HandleClient handles incoming connections from other Peers.
func (n *Peer) HandleClient(conn net.Conn) {
	defer conn.Close()

	// Perform communication with the connected Peer as needed
	// For simplicity, we're not receiving any data in this example.
}

// NewBootstrapPeer creates a new BootstrapPeer instance.
func NewBootstrapPeer() *BootstrapPeer {
	return &BootstrapPeer{
		Peers: make(map[int]*Peer),
	}
}

// RegisterPeer registers a Peer with the bootstrap Peer.
func (b *BootstrapPeer) RegisterPeer(Peer *Peer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Peers[Peer.ID] = Peer
	fmt.Printf("Peer %d registered with BootstrapPeer\n", Peer.ID)
}

// GetNeighbors returns the neighbors of a Peer based on its ID.
func (b *BootstrapPeer) GetNeighbors(PeerID int) []*Peer {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var neighbors []*Peer

	for id, Peer := range b.Peers {
		if id != PeerID {
			neighbors = append(neighbors, Peer)
		}
	}

	return neighbors
}

// NewP2PNetwork creates a new P2PNetwork instance.
func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		BootstrapPeer: NewBootstrapPeer(),
	}
}

// DisplayNetwork displays all the Peers in the P2P network and their connections.
func (p *P2PNetwork) DisplayNetwork() {
	p.BootstrapPeer.mu.RLock()
	defer p.BootstrapPeer.mu.RUnlock()

	fmt.Println("P2P Network:")
	for _, Peer := range p.BootstrapPeer.Peers {
		fmt.Printf("Peer %d connected to: ", Peer.ID)
		for _, neighbor := range Peer.Neighbors {
			fmt.Printf("%d ", neighbor.ID)
		}
		fmt.Println()
	}
}

// GetConnectedPorts returns the port numbers of all Peers connected to the server.
func (n *Peer) GetConnectedPorts() []int {
	// Lock neighbors for reading
	n.mu.RLock()
	defer n.mu.RUnlock()

	var ports []int

	for _, neighbor := range n.Neighbors {
		ports = append(ports, neighbor.Port)
	}

	return ports
}