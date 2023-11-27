package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

type Node struct {
	ID                     int
	address 			   Address
	bootstrap_address      Address
	is_bootstrap           bool
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
	block                  Block
	Neighbors              []*Address
	NeighborChans          []chan []byte // Each neighbor has a dedicated channel
	neighborChannel        chan []byte   // Channel for general communication with neighbors
	recentlySentBlocks     map[int]bool
	lastMinedSenderID      int
	listener 			 net.Listener
	connections 		 map[int]net.Conn
	blockchain 			 Blockchain
}

func NewNode(id, transactionThreshold, port int, blockchain Blockchain, bootstrap Address, is_bootstrap bool) *Node {
	return &Node{
		address: 			Address{IP: net.IPv4(127, 0, 0, 1), Port: port},
		ID:                   id,
		recentTransactions:   make(map[string]bool),
		transactionChannel:   make(chan string),
		blockChannel:         make(chan []byte),
		minedBlockChannel:    make(chan Block),
		transactionThreshold: transactionThreshold,
		recentlySentBlocks:   make(map[int]bool),
		connections:          make(map[int]net.Conn),
		blockchain:			  blockchain,
		bootstrap_address:    bootstrap,
		is_bootstrap:		  is_bootstrap,
	}
	
}

// TODO: function that constantly reads transactions from channel and adds them to block

func (node *Node) Start_Server(port int) {
	// Start server
	if node.is_bootstrap {
		fmt.Printf("Node %d assigned as bootstrap node\n", node.ID)
	}
	listenAddr := fmt.Sprintf("%s:%d", node.address.IP, node.address.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Peer %d: Error starting server: %v\n", node.ID, err)
		return
	}
	node.listener = listener
	fmt.Printf("Peer %d: Server started at %s\n", node.ID, listenAddr)

	// Accept incoming connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Peer %d: Error accepting connection: %v\n", node.ID, err)
				continue
			}
			fmt.Printf("Peer %d: Accepted connection from %s\n", node.ID, conn.RemoteAddr().String())
			// Add the connection to the list of connections
			node.connections[node.ID] = conn
			go node.HandleConnection(conn)
		}
	}()
}

func (node *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		var msg Message
		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Peer %d: Error decoding message: %v\n", node.ID, err)
			return
		}
		switch msg.Type {
		case TRANSACTION:
			// Process the received transaction
			fmt.Printf("Peer %d: Received transaction: %s\n", node.ID, msg.Payload.(string))
			node.ReceiveTransaction(msg.Payload.(string))
		case BLOCK:
			// Process the received mined block
			fmt.Printf("Peer %d: Received mined block: %x\n", node.ID, msg.Payload.(Block).Self_Hash)
			node.ReceiveMinedBlock(msg.Payload.(Block))
		
		case REQUEST:
			// Process the received request
			fmt.Println("Peer %d: Received missing blocks request: %s\n", node.ID, msg.Payload.(string))
			// send blockchain to requestee 
			node.SendBlockchain(msg)
		case BLOCKCHAIN:
			// Process the received blockchain
			fmt.Println("Peer %d: Received blockchain: %s\n", node.ID, msg.Payload.(string))
			// update blockchain
			node.blockchain = msg.Payload.(Blockchain)
		default:
			fmt.Printf("Peer %d: Unknown message type received\n", node.ID)
		}
		
	}
}

func (node *Node) SendBlockchain(msg Message) {
	blockchain := msg.Payload.(Blockchain)
	if blockchain.get_length() < node.blockchain.get_length() {
		// send blockchain to requestee
		node.SendMessage(msg.sender, BLOCKCHAIN, BlockchainPayload{node.blockchain})
	} else if blockchain.get_length() > node.blockchain.get_length() {
		// update blockchain
		node.blockchain = blockchain
	}

}

func (node *Node) SendMessage(address Address, msgType MessageType, payload interface{}) {
	// Encode the message
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(Message{msgType, node.address, payload}); err != nil {
		fmt.Printf("Peer %d: Error encoding message: %v\n", node.ID, err)
		return
	}

	// Send the encoded message to the address
	conn, err := net.Dial("tcp", address.String())
	if err != nil {
		fmt.Printf("Peer %d: Error sending message to %s: %v\n", node.ID, address.String(), err)
		return
	}
	defer conn.Close()

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Printf("Peer %d: Error sending message to %s: %v\n", node.ID, address.String(), err)
		return
	}
}

func (node *Node) AddNeighborChannel(neighborChannel chan string) {
	node.neighborChannels = append(node.neighborChannels, neighborChannel)
	//node.blockChannel = append(node.blockChannel, neighborChannel)
}
func (node *Node) AddNeighbor(neighbor *Node) {
	node.Neighbors = append(node.Neighbors, &neighbor.address)
	fmt.Printf("Node %d connected to neighbor %d at address %s\n", node.ID, neighbor.ID, neighbor.address.String())
	//node.NeighborChans = append(node.NeighborChans, neighborChannel)

}



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
	/* for _, neighborChannel := range node.neighborChannels {
		// Send the transaction to all neighbors
		neighborChannel <- transaction
	} */

	// broadcast to all neighbors through tcp
	for _, neighbor := range node.Neighbors {
		node.SendMessage(*neighbor, TRANSACTION, transaction)
	}
}

func (node *Node) MineBlock(transactions []string, previousHash string) Block {
	prevhash, err := hex.DecodeString(hash_to_string(hash_string(previousHash)))
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
	// TODO: make it so transactions are constantly being read, and block is made and mined as soon as threshold is reached

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
	previous_hash := hash_to_string(node.blockchain.get_last_block().Self_Hash)
	minedBlock := node.MineBlock(transactions, previous_hash)
	fmt.Println("Node ", node.ID, " Mined Block", hash_to_string(minedBlock.Self_Hash))
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

	/* for i, neighborChan := range node.NeighborChans {
		//if i != node.lastMinedSenderID {
		neighborChan <- buffer.Bytes()
		fmt.Printf("Node %d sent block to neighbor %d\n", node.ID, i)
		//}
	} */

	// broadcast to all neighbors through tcp
	for _, neighbor := range node.Neighbors {
		node.SendMessage(*neighbor, BLOCK, minedBlock)
	}

	// Broadcast the encoded block to the main block channel
	node.blockChannel <- buffer.Bytes()
	fmt.Printf("Node %d sent block to main channel\n", node.ID)
}



func (node *Node) ReceiveMinedBlock(minedBlock Block) {
	// Process the received mined block
	fmt.Printf("Node %d received mined block: %s\n", node.ID, minedBlock.Self_Hash)
	// check if block is already in blockchain
	if node.blockchain.search_block(minedBlock) {
		return
	}
	// check if block is valid
	if !node.blockchain.insert_block(minedBlock) {
		//TODO: request missing blocks
		

	}
	// add block to blockchain
	node.blockchain.insert_block(minedBlock)
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

	}
}


