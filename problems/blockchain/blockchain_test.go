package blockchain

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBlockchain_NewBlock(t *testing.T) {
	t.Parallel()
	var timestamp int64 = 1000
	transactions := []Transaction{
		{
			Sender:    "sender",
			Recipient: "recipient",
			Amount:    1,
		},
	}
	getTimestamp := func() int64 { return timestamp }
	bc := NewBlockchain(getTimestamp)

	bc.NewTransaction("sender", "recipient", 1)
	bc.NewBlock(200, nil)

	expectedChain := []Block{
		{
			Index:        0,
			Timestamp:    timestamp,
			Transactions: []Transaction{},
			Proof:        0,
			PreviousHash: "GENESIS",
		},
		{
			Index:        1,
			Timestamp:    timestamp,
			Transactions: transactions,
			Proof:        200,
			PreviousHash: "f75f875fef6b33c089d2c35b93287c2617db148703e2b04eb994e5b5140182cd",
		},
	}
	expectedTransactions := []Transaction{}

	if fmt.Sprintf("%v", bc.chain) != fmt.Sprintf("%v", expectedChain) {
		t.Errorf("%v != %v", bc.chain, expectedChain)
	}
	if !reflect.DeepEqual(bc.currentTransactions, expectedTransactions) {
		t.Errorf("%v != %v", bc.currentTransactions, expectedTransactions)
	}
}

func TestBlockchain_NewTransaction(t *testing.T) {
	t.Parallel()
	var timestamp int64 = 1000
	transactions := []Transaction{
		{
			Sender:    "sender",
			Recipient: "recipient",
			Amount:    1,
		},
	}
	getTimestamp := func() int64 { return timestamp }
	bc := NewBlockchain(getTimestamp)

	bc.NewTransaction("sender", "recipient", 1)

	if !reflect.DeepEqual(bc.currentTransactions, transactions) {
		t.Errorf("%v != %v", bc.currentTransactions, transactions)
	}
}
