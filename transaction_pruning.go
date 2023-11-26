package main

// TODO: Check block verification function
// TODO: Check functions in that block
// TODO: Remove transactions from list of verified transactions

//Checking for transactions in the block and removing them from the list of node transactions
func prunning(node_transactions []string, recv_block Block) [] string {
	//fmt.Println("Prunning")
	
	//iterating over the block trans
	for i := 0; i < len(recv_block.Transactions); i++ {
		//iterating over the current node transactions
		for j := 0; j < len(node_transactions); j++ {

			//if the transaction in block is found in the node transactions
			if recv_block.Transactions[i] == node_transactions[j] {

				//then remove it by appending the list of node transactions without the block transaction
				node_transactions = append(node_transactions[:j], node_transactions[j+1:]...)
				
				break
			}
		}
	}
	return node_transactions

}