package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// TODO: Creating a merkle tree base node (done)
// TODO: creating internal nodes to store pairs of transactions (done)
// TODO: creating Leaf nodes to store transactions	(done)
// TODO: Create function to store/update node data
// TODO: Create function to calculate hash of internal node
// TODO: Create function to traverse the tree and print the hashes/data of all nodes

//------Dummy Hashing Function------//
func calculateHash(data string) string {
	// Create a new SHA-256 hash
	hash := sha256.New()

	// Write the data to the hash
	hash.Write([]byte(data))

	// Get the final hash and convert it to a hex string
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString
}

// base node for merkel tree
type TreeNode interface {
	GetHash() string
}

// InternalNode represents internal nodes with left and right child nodes.
type InternalNode struct {
	Data        string
	Self_Hash   string
	Left_child  TreeNode // Can be either InternalNode or LeafNode
	Right_child TreeNode // Can be either InternalNode or LeafNode
}

// LeafNode represents to store transactions.
type LeafNode struct {
	Self_Hash   string
	Transaction string
}
// Implementing GetHash() method for LeafNode
func (leaf LeafNode) GetHash() string {
	return string(leaf.Self_Hash)
}
func (Internal InternalNode) GetHash() string {
	return string(Internal.Self_Hash)
}


//////////------------dummy functionality to create a small merkel tree------------//////////
func Create_Merkel(){

	//Root node of merkel tree
	var root InternalNode = InternalNode{
		Data:        "root",
		Self_Hash:   "",
		Left_child:  nil,
		Right_child: nil,
	}

	//Dummy transactions
	transaction1:="Transaction 1"
	transaction2:="Transaction 2"
	transaction3:="Transaction 3"
	transaction4:="Transaction 4"

	//creating leaf nodes to store each transaction
	var leaf1 LeafNode = LeafNode{
		Self_Hash:   calculateHash((transaction1)),
		Transaction: transaction1,
	}

	var leaf2 LeafNode = LeafNode{
		Self_Hash:   calculateHash((transaction2)),
		Transaction: transaction2,
	}

	var leaf3 LeafNode = LeafNode{
		Self_Hash:   calculateHash((transaction3)),
		Transaction: transaction3,
	}

	
	var leaf4 LeafNode = LeafNode{
		Self_Hash:   calculateHash((transaction4)),
		Transaction: transaction4,
	}


	//creating internal nodes to store pairs of transactions
	var left InternalNode = InternalNode{
		Data:        "left",
		Self_Hash:   "",
		Left_child:  &leaf1,
		Right_child: &leaf2,
	}

	var right InternalNode = InternalNode{
		Data:        "right",
		Self_Hash:   "",
		Left_child:  &leaf3,
		Right_child: &leaf4,
	}
	
	//calculating hashes of leaf nodes and updating internal nodes along with the data of the no
	left.Data = left.Left_child.GetHash() + ","+left.Right_child.GetHash()
	right.Data = right.Left_child.GetHash() +","+ right.Right_child.GetHash()

	left.Self_Hash = calculateHash(left.Data)
	right.Self_Hash = calculateHash(right.Data)

	//connecting the two internal nodes to the root node
	//calculating hash of root node and updating the data of the root node
	root.Left_child = &left
	root.Right_child = &right
	root.Data = root.Left_child.GetHash() + ","+ root.Right_child.GetHash()
	root.Self_Hash = calculateHash(root.Data)

	//printing the hashes of all nodes
	fmt.Println("Hash of root node: ",root.Self_Hash)
	fmt.Println("data on root node: ",root.Data)
	fmt.Println("Hash of root's left child node: ",root.Left_child.GetHash())
	fmt.Println("Hash of root's left child  node: ",root.Right_child.GetHash())
	fmt.Println("----------------------------------------")
	fmt.Println("Hash of transactions on root's left child  node: ",root.Left_child.(*InternalNode).Left_child.GetHash())
	fmt.Println("Hash of transactions on root's left child  node: ",root.Left_child.(*InternalNode).Right_child.GetHash())
	fmt.Println("----------------------------------------")
	fmt.Println("Hash of transactions on root's right child  node: ",root.Right_child.(*InternalNode).Left_child.GetHash())
	fmt.Println("Hash of transactions on root's right child  node: ",root.Right_child.(*InternalNode).Right_child.GetHash())


}






