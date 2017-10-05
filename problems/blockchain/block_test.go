package blockchain

import (
	"testing"
)

func TestBlock_ToJson(t *testing.T) {
	t.Parallel()
	block := Block{
		Index:        0,
		Timestamp:    10,
		Transactions: []Transaction{{Sender: "1", Recipient: "2", Amount: 3}},
		Proof:        100,
		PreviousHash: "ABC",
	}
	result := string(block.ToJson())
	expected := `{"index":0,"timestamp":10,"transactions":[{"sender":"1","recipient":"2","amount":3}],"proof":100,"previous_hash":"ABC"}`
	if result != expected {
		t.Errorf("%v != %v", result, expected)
	}
}

func TestBlock_Hash(t *testing.T) {
	t.Parallel()
	expected := "b5d4045c3f466fa91fe2cc6abe79232a1a57cdf104f7a26e716e0a1e2789df78"
	result := hash([]byte("ABC"))
	if result != expected {
		t.Errorf("%v != %v", result, expected)
	}
}
