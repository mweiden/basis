package blockchain

import (
	"fmt"
	"strings"
)

type HashCash struct {
	Difficulty int
}

func (hc *HashCash) ProofOfWork(lastProof uint) uint {
	var proof uint = 0
	for !hc.ValidProof(lastProof, proof) {
		proof += 1
	}
	return proof
}

func (hc *HashCash) ValidProof(lastProof uint, proof uint) bool {
	guess := hash([]byte(solutionKey(lastProof, proof)))
	return guess[:2] == strings.Repeat("0", hc.Difficulty)
}

func solutionKey(lastProof uint, proof uint) string {
	return fmt.Sprintf(`namespace:%s:%s`, lastProof, proof)
}
