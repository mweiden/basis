package packet_muxer

import (
	"errors"
	"github.com/mweiden/basis/datastructures"
	"sort"
)

var (
	MXE = errors.New("Muxer empty!")
)

// packet and packet stack (pipe)
type Packet struct {
	id      int
	payload []byte
}

func NewPacketStack(packets ...Packet) datastructures.Stack {
	sort.Slice(
		packets[:],
		func(i, j int) bool {
			return packets[i].id > packets[i].id
		},
	)
	s := datastructures.Stack{}
	for _, p := range packets {
		s.Push(p)
	}
	return s
}

// packet muxer
type PacketMuxer struct {
	pipes   []datastructures.Stack
	minHeap datastructures.Heap
}

type muxerTuple struct {
	pipeInd int
	packet  Packet
}

func NewPacketMuxer(pipes ...datastructures.Stack) PacketMuxer {
	minHeap := datastructures.Heap{
		Compare: func(a interface{}, b interface{}) int {
			return a.(muxerTuple).packet.id - b.(muxerTuple).packet.id
		},
	}
	pm := PacketMuxer{
		pipes:   pipes,
		minHeap: minHeap,
	}
	for i, _ := range pm.pipes {
		pm.pollPipe(i)
	}
	return pm
}

func (pm *PacketMuxer) HasNext() bool {
	return pm.minHeap.Size() > 0
}

func (pm *PacketMuxer) pollPipe(i int) {
	nextFromPipe, pipeErr := pm.pipes[i].Pop()
	if pipeErr != datastructures.EOS {
		pm.minHeap.Insert(muxerTuple{i, nextFromPipe.(Packet)})
	}
}

func (pm *PacketMuxer) Next() (Packet, error) {
	if pm.HasNext() {
		tupleInterface, _ := pm.minHeap.Pop()
		tuple := tupleInterface.(muxerTuple)
		pm.pollPipe(tuple.pipeInd)
		return tuple.packet, nil
	} else {
		return Packet{}, MXE
	}
}
