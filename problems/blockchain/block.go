package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

type Block struct {
	Index        uint          `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        uint          `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

func (b *Block) ToJson() []byte {
	bytes, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return bytes
}

func (b *Block) Hash() string {
	return hash(b.ToJson())
}

func hash(b []byte) string {
	var bytes []byte
	for _, b := range sha256.Sum256(b) {
		bytes = append(bytes, b)
	}
	return hex.EncodeToString(bytes)
}
