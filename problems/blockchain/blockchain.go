package blockchain

type Blockchain struct {
	hashCash       HashCash
	chain               []Block
	currentTransactions []Transaction
	getTimestamp        func() int64
}

func NewBlockchain(getTimestamp func() int64) *Blockchain {
	bc := Blockchain{
		hashCash: HashCash{Difficulty: 2},
		getTimestamp: getTimestamp,
	}
	bc.NewBlock(0, "GENESIS")
	return &bc
}

func (b *Blockchain) NewTransaction(sender string, recipient string, amount int) *Transaction {
	transation := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	b.currentTransactions = append(b.currentTransactions, transation)
	return &transation
}

func (b *Blockchain) NewBlock(proof uint, maybePreviousHash interface{}) Block {
	var prev string
	if maybePreviousHash == nil {
		prev = b.chain[len(b.chain)-1].Hash()
	} else {
		prev = maybePreviousHash.(string)
	}

	block := Block{
		Index:        uint(len(b.chain)),
		Timestamp:    b.getTimestamp(),
		Transactions: b.currentTransactions,
		Proof:        proof,
		PreviousHash: prev,
	}
	b.currentTransactions = b.currentTransactions[:0]
	b.chain = append(b.chain, block)
	return block
}

func (b *Blockchain) Chain() []Block {
	return b.chain
}

func (b *Blockchain) LastBlock() Block {
	return b.chain[len(b.chain)-1]
}

func (b *Blockchain) ValidChain() bool {
	lastBlock := b.chain[0]
	currentIndex := 1
	for currentIndex < len(b.chain) {
		block := b.chain[currentIndex]
		if block.PreviousHash != lastBlock.Hash() {
			return false
		}
		if !b.hashCash.ValidProof(lastBlock.Proof, block.Proof) {
			return false
		}
		lastBlock = block
		currentIndex += 1
	}
	return true
}

