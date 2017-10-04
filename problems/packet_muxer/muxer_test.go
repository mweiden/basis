package packet_muxer

import (
	"reflect"
	"testing"
)

func TestNewPacketMuxer(t *testing.T) {
	t.Parallel()
	p1 := NewPacketStack(Packet{id: 8}, Packet{id: 6}, Packet{id: 3}, Packet{id: 1}, Packet{id: 1})
	p2 := NewPacketStack(Packet{id: 7}, Packet{id: 6}, Packet{id: 2}, Packet{id: 2}, Packet{id: 1})
	p3 := NewPacketStack(Packet{id: 9}, Packet{id: 5}, Packet{id: 4}, Packet{id: 3}, Packet{id: 3})
	p4 := NewPacketStack(Packet{id: 7}, Packet{id: 4}, Packet{id: 3}, Packet{id: 2}, Packet{id: 2})

	type Result struct {
		packet Packet
		err    error
	}

	var results []Result
	expected := []Result{
		Result{Packet{id: 1}, nil},
		Result{Packet{id: 1}, nil},
		Result{Packet{id: 1}, nil},
		Result{Packet{id: 2}, nil},
		Result{Packet{id: 2}, nil},
		Result{Packet{id: 2}, nil},
		Result{Packet{id: 2}, nil},
		Result{Packet{id: 3}, nil},
		Result{Packet{id: 3}, nil},
		Result{Packet{id: 3}, nil},
		Result{Packet{id: 3}, nil},
		Result{Packet{id: 4}, nil},
		Result{Packet{id: 4}, nil},
		Result{Packet{id: 5}, nil},
		Result{Packet{id: 6}, nil},
		Result{Packet{id: 6}, nil},
		Result{Packet{id: 7}, nil},
		Result{Packet{id: 7}, nil},
		Result{Packet{id: 8}, nil},
		Result{Packet{id: 9}, nil},
		Result{Packet{}, MXE},
	}
	pm := NewPacketMuxer(p1, p2, p3, p4)
	for pm.HasNext() {
		next, err := pm.Next()
		results = append(results, Result{next, err})
	}

	next, err := pm.Next()
	results = append(results, Result{next, err})

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("%v != %v", results, expected)
	}
}
