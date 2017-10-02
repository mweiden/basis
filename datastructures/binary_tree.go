package datastructures

import (
	"errors"
	"math"
)


var (
	DNE = errors.New("Node does not exist!")
	EMPTY = tombstone{}
)

type BinaryTree struct {
	ary []interface{}
}

type tombstone struct {}

type Node struct {
	Id  int
	Val interface{}
}

func (b *BinaryTree) Diameter() int {
	root, err := b.GetNode(0)
	if err == DNE {
		return 0
	}
	toSearch := Stack{}
	toSearch.Push(root.Id)
	var diameters []int

	for !toSearch.Empty() {
		popped, _ := toSearch.Pop()
		parentId := popped.(int)
		left, leftErr := b.LeftChild(parentId)
		right, rightErr := b.RightChild(parentId)
		currentDiameter := 0
		if leftErr == nil {
			currentDiameter += 1 + b.maxDepth(left.Id)
			toSearch.Push(left.Id)
		}
		if rightErr == nil {
			currentDiameter += 1 + b.maxDepth(right.Id)
			toSearch.Push(right.Id)
		}
		diameters = append(diameters, currentDiameter)
	}
	return max(diameters)
}

func max(s []int) int {
	mx := math.MinInt32
	for _, val := range s {
		if val > mx {
			mx = val
		}
	}
	return mx
}

// TODO: memoize this
func (b *BinaryTree) maxDepth(id int) int {
	type Tuple struct {
		id    int
		depth int
	}
	startingNode, err := b.GetNode(id)
	if err == DNE {
		return 0
	}
	toExpand := Stack{}
	toExpand.Push(Tuple{startingNode.Id, 0})
	maxDepth := 0

	check := func(id int, currentDepth int, err error) {
		if err == nil {
			newDepth := currentDepth + 1
			if newDepth > maxDepth {
				maxDepth = newDepth
			}
			toExpand.Push(Tuple{id, newDepth})
		}
	}

	for !toExpand.Empty() {
		popped, _ := toExpand.Pop()
		parent := popped.(Tuple)
		left, leftErr := b.LeftChild(parent.id)
		right, rightErr := b.RightChild(parent.id)
		check(right.Id, parent.depth, rightErr)
		check(left.Id, parent.depth, leftErr)
	}

	return maxDepth
}

func (b *BinaryTree) Size() int {
	return len(b.ary)
}

func (b *BinaryTree) GetNode(ind int) (Node, error) {
	if ind >= b.Size() || b.ary[ind] == EMPTY {
		return Node{ind, nil}, DNE
	} else {
		return Node{ind, b.ary[ind]}, nil
	}
}

func (b *BinaryTree) SetNode(n Node) {
	diff := b.Size() - n.Id
	if diff > 0 {
		b.ary[n.Id] = n.Val
	} else {
		if diff == 0 {
			b.ary = append(b.ary, n.Val)
		} else {
			tmp := make([]interface{}, n.Id+1)
			for i, _ := range tmp {
				tmp[i] = EMPTY
			}
			copy(tmp, b.ary)
			tmp[n.Id] = n.Val
			b.ary = tmp
		}
	}
}

func (b *BinaryTree) LeftChild(i int) (Node, error) {
	ind := i*2 + 1
	return b.GetNode(ind)
}

func (b *BinaryTree) RightChild(i int) (Node, error) {
	ind := i*2 + 2
	return b.GetNode(ind)
}

func (b *BinaryTree) Parent(i int) (Node, error) {
	ind := (i - 1) / 2
	return b.GetNode(ind)
}

func (b *BinaryTree) Swap(i int, j int) {
	ival := b.ary[i]
	b.ary[i] = b.ary[j]
	b.ary[j] = ival
}
