package main

type MessageType int

const (
	// Message types
	TRANSACTION MessageType = iota
	BLOCK
	REQUEST
	BLOCKCHAIN
)

type Message struct {
	Type MessageType
	sender Address
	Payload interface{}
}

type TransactionPayload struct {
	Transaction string
}

type BlockPayload struct {
	Block Block
}

type RequestPayload struct {
	blockchain Blockchain
}

type BlockchainPayload struct {
	Blockchain Blockchain
}