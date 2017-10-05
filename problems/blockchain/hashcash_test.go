package blockchain

import (
	"testing"
)

func TestHashCash_ProofOfWork(t *testing.T) {
	t.Parallel()
	difficulty := 2
	hc := HashCash{difficulty}

	for _, lastProof := range []uint{0, 1, 22, 333} {
		proof := hc.ProofOfWork(lastProof)
		solution := hash([]byte(solutionKey(lastProof, proof)))
		if solution[:difficulty] != "00" {
			t.Errorf("solution %s does not have %d leading zeros!", solution, difficulty)
		}
	}
}
