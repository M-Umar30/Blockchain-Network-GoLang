/*
	20I-0518 - Muhammad Umar
	20I-0450 - Fatima Zubeda
	20L-
	20I-
*/

// Assignment 1: Block Chain Implementation

// TODO: Seperate files for all hash related functions?
// TODO: Calculate hash of block
// TODO: Complete merkel tree implementation
// TODO: Block mining
// TODO: Block validation
// TODO: Chain validation

// Assignment 2: P2P Implementation

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")

	transactions := []string{"Transaction 1", "Transaction 2", "Transaction 3", "Transaction 4"}
	root := CreateMerkleTree(transactions)
	// transactions := []string{"Transaction 1", "Transaction 2", "Transaction 3"}
	// root := CreateMerkleTree(transactions)
	fmt.Println("Initial Merkle Tree:")
	//PrintMerkleTree(root, 0)

	root.DisplayMerkleTree()

	// // Add a new transaction.
	// newTransaction := "New Transaction"
	// transactions = append(transactions, newTransaction)

	// // Update the Merkle tree with the new transaction.
	// root = CreateMerkleTree(transactions)
	// fmt.Printf("\nMerkle Tree after adding a new transaction (%s):\n", newTransaction)
	// root.DisplayMerkleTree()

	// // Add a new transaction
	// newTransaction := "New Transaction"
	// fmt.Printf("\nAdding transaction: %s\n", newTransaction)
	// root.AddTransaction(newTransaction)

	// // Print the tree after adding the transaction
	// fmt.Printf("\nMerkle Tree after adding a new transaction:")
	// root.DisplayMerkleTree()

	// // Verify inclusion of the new transaction
	// fmt.Printf("\nVerifying inclusion of transaction: %s\n", newTransaction)
	// if root.VerifyInclusion(newTransaction) {
	// 	fmt.Println("Verification: Transaction is included.")
	// } else {
	// 	fmt.Println("Verification: Transaction is NOT included.")
	// }

	// // Get the root hash
	// rootHash := root.GetRootHash()
	// fmt.Printf("\nRoot Hash after adding a new transaction: %s\n", rootHash)

	// // Add a new transaction that is not in the Merkle tree
	// nonIncludedTransaction := "Non Included Transaction"
	// fmt.Printf("\nVerifying inclusion of transaction: %s\n", nonIncludedTransaction)
	// if root.VerifyInclusion(nonIncludedTransaction) {
	// 	fmt.Println("Verification: Transaction is included.")
	// } else {
	// 	fmt.Println("Verification: Transaction is NOT included.")
	// }

	// Update a transaction
	oldTransaction := "Transaction 1"
	newTransaction2 := "Updated Transaction 1"
	fmt.Printf("\nUpdating transaction: %s to %s\n", oldTransaction, newTransaction2)
	root.UpdateTransaction(oldTransaction, newTransaction2)

	// Print the tree after updating the transaction
	fmt.Printf("\nUpdated Merkle Tree:")
	root.DisplayMerkleTree()

	// Verify inclusion of the updated transaction
	fmt.Printf("\nVerifying inclusion of transaction: %s\n", newTransaction2)
	if root.VerifyInclusion(newTransaction2) {
		fmt.Println("Verification: Transaction is included.")
	} else {
		fmt.Println("Verification: Transaction is NOT included.")
	}
	// // Test removing a transaction
	// fmt.Println("\nRemoving transaction: Transaction 2")
	// // Specify the transaction ID to be removed
	// transactionToRemove := "Transaction 2"

	// // Remove the specified transaction
	// root, removed := root.RemoveTransaction(transactionToRemove)
	// fmt.Println("\nUpdated Merkle Tree:")
	// root.DisplayMerkleTree()
	// print("\n", removed)

	// // Test verifying inclusion after removal
	// fmt.Println("\nVerifying inclusion of removed transaction: Transaction 2")
	// if root.VerifyInclusion("Transaction 2") {
	// 	fmt.Println("Verification: Transaction is included.")
	// } else {
	// 	fmt.Println("Verification: Transaction is NOT included.")
	// }

}
