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

// ------Dummy Hashing Function------//
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

// CreateMerkleTree creates a balanced Merkle tree.
func CreateMerkleTree(transactions []string) *InternalNode {
	// Ensure the number of transactions is a multiple of 4 by duplicating the last one.
	for len(transactions)%4 != 0 {
		transactions = append(transactions, transactions[len(transactions)-1])
	}

	var leaves []TreeNode

	// Create leaf nodes for each transaction.
	for _, tx := range transactions {
		leaf := &LeafNode{
			Self_Hash:   calculateHash(tx),
			Transaction: tx,
		}
		leaves = append(leaves, leaf)
	}

	// Build the Merkle tree by grouping transactions in pairs.
	for len(leaves) > 1 {
		var level []TreeNode
		for i := 0; i < len(leaves); i += 2 {
			internal := &InternalNode{
				Data:        leaves[i].GetHash() + "," + leaves[i+1].GetHash(),
				Self_Hash:   calculateHash(leaves[i].GetHash() + "," + leaves[i+1].GetHash()),
				Left_child:  leaves[i],
				Right_child: leaves[i+1],
			}
			level = append(level, internal)
		}
		leaves = level
	}

	// The last remaining node is the root of the Merkle tree.
	return leaves[0].(*InternalNode)
}

// ////////------------dummy functionality to create a small merkel tree------------//////////
// func Create_Merkel() *InternalNode {

// 	//Root node of merkel tree
// 	var root InternalNode = InternalNode{
// 		Data:        "root",
// 		Self_Hash:   "",
// 		Left_child:  nil,
// 		Right_child: nil,
// 	}

// 	//Dummy transactions
// 	transaction1 := "Transaction 1"
// 	transaction2 := "Transaction 2"
// 	transaction3 := "Transaction 3"
// 	transaction4 := "Transaction 4"

// 	//creating leaf nodes to store each transaction
// 	var leaf1 LeafNode = LeafNode{
// 		Self_Hash:   calculateHash((transaction1)),
// 		Transaction: transaction1,
// 	}

// 	var leaf2 LeafNode = LeafNode{
// 		Self_Hash:   calculateHash((transaction2)),
// 		Transaction: transaction2,
// 	}

// 	var leaf3 LeafNode = LeafNode{
// 		Self_Hash:   calculateHash((transaction3)),
// 		Transaction: transaction3,
// 	}

// 	var leaf4 LeafNode = LeafNode{
// 		Self_Hash:   calculateHash((transaction4)),
// 		Transaction: transaction4,
// 	}

// 	//creating internal nodes to store pairs of transactions
// 	var left InternalNode = InternalNode{
// 		Data:        "left",
// 		Self_Hash:   "",
// 		Left_child:  &leaf1,
// 		Right_child: &leaf2,
// 	}

// 	var right InternalNode = InternalNode{
// 		Data:        "right",
// 		Self_Hash:   "",
// 		Left_child:  &leaf3,
// 		Right_child: &leaf4,
// 	}

// 	//calculating hashes of leaf nodes and updating internal nodes along with the data of the no
// 	left.Data = left.Left_child.GetHash() + "," + left.Right_child.GetHash()
// 	right.Data = right.Left_child.GetHash() + "," + right.Right_child.GetHash()

// 	left.Self_Hash = calculateHash(left.Data)
// 	right.Self_Hash = calculateHash(right.Data)

// 	//connecting the two internal nodes to the root node
// 	//calculating hash of root node and updating the data of the root node
// 	root.Left_child = &left
// 	root.Right_child = &right
// 	root.Data = root.Left_child.GetHash() + "," + root.Right_child.GetHash()
// 	root.Self_Hash = calculateHash(root.Data)

// 	// //printing the hashes of all nodes
// 	// fmt.Println("Hash of root node: ", root.Self_Hash)
// 	// fmt.Println("data on root node: ", root.Data)
// 	// fmt.Println("Hash of root's left child node: ", root.Left_child.GetHash())
// 	// fmt.Println("Hash of root's left child  node: ", root.Right_child.GetHash())
// 	// fmt.Println("----------------------------------------")
// 	// fmt.Println("Hash of transactions on root's left child  node: ", root.Left_child.(*InternalNode).Left_child.GetHash())
// 	// fmt.Println("Hash of transactions on root's left child  node: ", root.Left_child.(*InternalNode).Right_child.GetHash())
// 	// fmt.Println("----------------------------------------")
// 	// fmt.Println("Hash of transactions on root's right child  node: ", root.Right_child.(*InternalNode).Left_child.GetHash())
// 	// fmt.Println("Hash of transactions on root's right child  node: ", root.Right_child.(*InternalNode).Right_child.GetHash())
// 	return &root
// }

// // AddTransaction adds a new transaction to the Merkle tree.
// func (root *InternalNode) AddTransaction(transaction string) {
// 	// Create a new leaf node for the transaction
// 	newLeaf := LeafNode{
// 		Self_Hash:   calculateHash(transaction),
// 		Transaction: transaction,
// 	}

// 	// Create a new internal node to store the newLeaf and the right child of the root
// 	newInternal := InternalNode{
// 		Data:        "newInternal",
// 		Self_Hash:   "",
// 		Left_child:  &newLeaf,
// 		Right_child: root.Right_child,
// 	}

// 	// Update data and hash of the newInternal node
// 	newInternal.Data = newInternal.Left_child.GetHash() + "," + newInternal.Right_child.GetHash()
// 	newInternal.Self_Hash = calculateHash(newInternal.Data)

// 	// Update data and hash of the root node
// 	root.Left_child = &newInternal
// 	root.Data = root.Left_child.GetHash() + "," + root.Right_child.GetHash()
// 	root.Self_Hash = calculateHash(root.Data)

// }

// // AddTransaction adds a new transaction to the Merkle tree.
// func (root *InternalNode) AddTransaction(transaction string) {
// 	// Create a new leaf node for the transaction
// 	newLeaf := LeafNode{
// 		Self_Hash:   calculateHash(transaction),
// 		Transaction: transaction,
// 	}

// 	// Check if the tree has an odd number of leaf nodes
// 	if root.Right_child == nil {
// 		// Duplicate the last leaf node (New Transaction)
// 		duplicateLeaf := LeafNode{
// 			Self_Hash:   calculateHash(transaction),
// 			Transaction: transaction,
// 		}
// 		root.Right_child = &duplicateLeaf
// 	} else {
// 		// Create a new internal node for the newLeaf and the right child of the root
// 		newInternal := InternalNode{
// 			Data:        "newInternal",
// 			Self_Hash:   "",
// 			Left_child:  &newLeaf,
// 			Right_child: root.Right_child,
// 		}

// 		// Update data and hash of the newInternal node
// 		newInternal.Data = newInternal.Left_child.GetHash() + "," + newInternal.Right_child.GetHash()
// 		newInternal.Self_Hash = calculateHash(newInternal.Data)

// 		// Replace the right child of the root with the newInternal node
// 		root.Right_child = &newInternal
// 	}

// 	// Create a new internal node for the combined hash of the left and right children of the root
// 	combinedInternal := InternalNode{
// 		Data:        "combinedInternal",
// 		Self_Hash:   "",
// 		Left_child:  root.Left_child,
// 		Right_child: root.Right_child,
// 	}

// 	// Update data and hash of the combinedInternal node
// 	combinedInternal.Data = combinedInternal.Left_child.GetHash() + "," + combinedInternal.Right_child.GetHash()
// 	combinedInternal.Self_Hash = calculateHash(combinedInternal.Data)

// 	// Replace the left child of the root with the combinedInternal node
// 	root.Left_child = &combinedInternal

// 	// Update data and hash of the root node
// 	root.Data = root.Left_child.GetHash() + "," + root.Right_child.GetHash()
// 	root.Self_Hash = calculateHash(root.Data)
// }

// VerifyInclusion verifies if a transaction is included in the Merkle tree.
func (root *InternalNode) VerifyInclusion(transaction string) bool {
	// Create a leaf node for the given transaction
	transactionLeaf := LeafNode{
		Self_Hash:   calculateHash(transaction),
		Transaction: transaction,
	}
	print("hash: ", transactionLeaf.Self_Hash)
	print("transaction: \n", transactionLeaf.Transaction)

	return verifyInclusionRecursive(root, &transactionLeaf)
}

func verifyInclusionRecursive(node TreeNode, transactionLeaf *LeafNode) bool {
	switch n := node.(type) {
	case *InternalNode:
		leftHash := calculateHash(n.Left_child.GetHash() + "," + n.Right_child.GetHash())
		rightHash := calculateHash(n.Right_child.GetHash() + "," + n.Left_child.GetHash())

		return verifyInclusionRecursive(n.Left_child, transactionLeaf) ||
			verifyInclusionRecursive(n.Right_child, transactionLeaf) ||
			(leftHash == n.Self_Hash && verifyInclusionRecursive(n.Right_child, transactionLeaf)) ||
			(rightHash == n.Self_Hash && verifyInclusionRecursive(n.Left_child, transactionLeaf))
	case *LeafNode:
		return n.GetHash() == transactionLeaf.GetHash()
	default:
		return false
	}
}

// PrintTree prints the hashes of all nodes in the Merkle tree.
func (root *InternalNode) PrintTree() {
	fmt.Println("Hash of root node: ", root.Self_Hash)
	printNode(root.Left_child, "Root's Left Child")
	printNode(root.Right_child, "Root's Right Child")
}

func printNode(node TreeNode, nodeName string) {
	if internalNode, ok := node.(*InternalNode); ok {
		fmt.Println("Hash of", nodeName, "node: ", internalNode.Self_Hash)
		printNode(internalNode.Left_child, nodeName+"'s Left Child")
		printNode(internalNode.Right_child, nodeName+"'s Right Child")
	} else if leafNode, ok := node.(*LeafNode); ok {
		fmt.Println("Hash of", nodeName, "node (Transaction): ", leafNode.Self_Hash, "(", leafNode.Transaction, ")")
	}
}

// GetRootHash returns the hash of the root node.
func (root *InternalNode) GetRootHash() string {
	return root.Self_Hash
}

// GetMerklePath returns the Merkle path for a given transaction in the Merkle tree.
func (root *InternalNode) GetMerklePath(transaction string) []string {
	var path []string
	getMerklePathRecursive(root, transaction, &path)
	return path
}

func getMerklePathRecursive(node TreeNode, targetTransaction string, path *[]string) bool {
	switch n := node.(type) {
	case *InternalNode:
		leftHash := calculateHash(n.Left_child.GetHash() + "," + n.Right_child.GetHash())
		rightHash := calculateHash(n.Right_child.GetHash() + "," + n.Left_child.GetHash())

		if getMerklePathRecursive(n.Left_child, targetTransaction, path) {
			*path = append(*path, rightHash)
			return true
		} else if getMerklePathRecursive(n.Right_child, targetTransaction, path) {
			*path = append(*path, leftHash)
			return true
		} else if leftHash == n.Self_Hash {
			*path = append(*path, rightHash)
			return false
		} else if rightHash == n.Self_Hash {
			*path = append(*path, leftHash)
			return false
		}
	case *LeafNode:
		if n.Transaction == targetTransaction {
			return true
		}
	}

	return false
}

// VerifyMerklePath verifies the Merkle path for a given transaction.
func VerifyMerklePath(rootHash string, transaction string, path []string) bool {
	currentHash := calculateHash(transaction)

	for _, siblingHash := range path {
		currentHash = calculateHash(currentHash + "," + siblingHash)
	}

	return currentHash == rootHash
}

// Helper function to update the hashes of nodes up the tree
func updateNodeHashes(node *InternalNode, leafNode *LeafNode) {
	if node != nil && leafNode != nil {
		updateNodeHash(node)
	}
}

// Helper function to get the parent of a node in the Merkle tree
func getParent(root *InternalNode, node TreeNode) TreeNode {
	if root != nil && node != nil {
		if root.Left_child == node || root.Right_child == node {
			return root
		}

		if leftParent := getParent(root.Left_child.(*InternalNode), node); leftParent != nil {
			return leftParent
		}

		return getParent(root.Right_child.(*InternalNode), node)
	}

	return nil
}

// findLeafNode finds a leaf node with the given transaction in the Merkle tree.
func findLeafNode(root TreeNode, targetTransaction string) (*InternalNode, *LeafNode, bool) {
	if root == nil {
		return nil, nil, false
	}

	switch root := root.(type) {
	case *InternalNode:
		leftParent, leftLeaf, leftFound := findLeafNode(root.Left_child, targetTransaction)
		if leftFound {
			return leftParent, leftLeaf, true
		}

		rightParent, rightLeaf, rightFound := findLeafNode(root.Right_child, targetTransaction)
		if rightFound {
			return rightParent, rightLeaf, true
		}

	case *LeafNode:
		if root.Transaction == targetTransaction {
			return nil, root, true
		}
	}

	return nil, nil, false
}

// Helper function to update the hash of nodes up the tree
func updateNodeHash(node TreeNode) string {
	if internalNode, ok := node.(*InternalNode); ok {
		leftHash := updateNodeHash(internalNode.Left_child)
		rightHash := updateNodeHash(internalNode.Right_child)
		internalNode.Self_Hash = calculateHash(leftHash + "," + rightHash)
		return internalNode.Self_Hash
	} else if leafNode, ok := node.(*LeafNode); ok {
		return leafNode.GetHash()
	}
	return ""
}

func findLeafNodeRecursive(node TreeNode, targetTransaction string) (*LeafNode, bool) {
	switch n := node.(type) {
	case *InternalNode:
		leftLeaf, leftFound := findLeafNodeRecursive(n.Left_child, targetTransaction)
		rightLeaf, rightFound := findLeafNodeRecursive(n.Right_child, targetTransaction)

		if leftFound {
			return leftLeaf, true
		} else if rightFound {
			return rightLeaf, true
		}
	case *LeafNode:
		if n.Transaction == targetTransaction {
			return n, true
		}
	}

	return nil, false
}

// UpdateTransaction updates a transaction in the Merkle tree.
func (root *InternalNode) UpdateTransaction(oldTransaction, newTransaction string) {
	// Find the leaf node corresponding to the old transaction
	leafNode, found := findLeafNodeRecursive(root, oldTransaction)
	if !found {
		fmt.Printf("Transaction not found: %s\n", oldTransaction)
		return
	}

	// Update the transaction in the leaf node
	leafNode.Transaction = newTransaction
	leafNode.Self_Hash = calculateHash(newTransaction)

	// Update hashes of all nodes along the path to the root
	updateNodeHashes(root, leafNode)

	// Update the root hash
	root.Data = root.Left_child.GetHash() + "," + root.Right_child.GetHash()
	root.Self_Hash = calculateHash(root.Data)
}

// DisplayMerkleTree prints a visual representation of the Merkle tree.
func (root *InternalNode) DisplayMerkleTree() {
	displayTree(root, "")
}
func displayTree(node TreeNode, prefix string) {
	if node == nil {
		fmt.Printf("%s├── Nil Node\n", prefix)
		return
	}

	fmt.Printf("%s├── %s\n", prefix, getNodeString(node))

	switch n := node.(type) {
	case *InternalNode:
		if n.Left_child != nil {
			displayTree(n.Left_child, prefix+"│   ")
		}
		if n.Right_child != nil {
			displayTree(n.Right_child, prefix+"│   ")
		}
	}
}

func getNodeString(node TreeNode) string {
	if node == nil {
		return "Nil Node"
	}

	switch n := node.(type) {
	case *InternalNode:
		return fmt.Sprintf("Internal Node (Hash: %s)", n.Self_Hash)
	case *LeafNode:
		return fmt.Sprintf("Leaf Node (Hash: %s, Transaction: %s)", n.Self_Hash, n.Transaction)
	default:
		return "Unknown Node Type"
	}
}

// // RemoveTransaction removes a transaction from the Merkle tree.
// func (root *InternalNode) RemoveTransaction(transactionID string) (*InternalNode, bool) {
// 	// Use the helper function to find the leaf node to be removed
// 	parent, leafToRemove, found := findLeafNode(root, transactionID)

// 	if !found {
// 		return root, false // Transaction not found
// 	}

// 	// Use the helper function to remove the leaf node
// 	root = removeLeafNode(root, parent, leafToRemove)

// 	// Recalculate hashes up the tree
// 	updateNodeHashes(root, leafToRemove)

// 	return root, true // Transaction successfully removed
// }

// // Helper function to remove a leaf node from the Merkle tree.
// func removeLeafNode(root *InternalNode, parent *InternalNode, leafNode *LeafNode) *InternalNode {
// 	if root == nil || leafNode == nil {
// 		return root
// 	}

// 	if parent == nil {
// 		// The leafNode is the root
// 		return nil
// 	}

// 	// Identify if the leafNode is the left or right child of the parent
// 	isLeftChild := parent.Left_child == leafNode

// 	// Replace the leafNode with nil in the parent
// 	if isLeftChild {
// 		parent.Left_child = nil
// 	} else {
// 		parent.Right_child = nil
// 	}

// 	return root
// }
