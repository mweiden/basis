package datastructures

import (
	"errors"
)

var (
	EOH = errors.New("End of Heap!")
)

type Heap struct {
	tree    BinaryTree
	Compare func(a interface{}, b interface{}) int
}

func (h *Heap) Insert(val interface{}) {
	// append to end and then bubble up
	h.tree.ary = append(h.tree.ary, val)
	i := h.tree.Size() - 1
	for i > 0 {
		parent, _ := h.tree.Parent(i)
		child, _ := h.tree.GetNode(i)
		if h.Compare(parent.Val, child.Val) <= 0 {
			break
		} else {
			h.tree.Swap(parent.Id, i)
		}
		i = parent.Id
	}
}

func (h *Heap) Pop() (interface{}, error) {
	if h.tree.Size() == 0 {
		return nil, EOH
	}
	// pull from front, replace with last, then bubble down
	root, _ := h.tree.GetNode(0)
	leaf, _ := h.tree.GetNode(h.tree.Size() - 1)
	leaf.Id = 0
	h.tree.SetNode(leaf)
	h.tree.ary = h.tree.ary[:h.tree.Size()-1]

	testParentId := leaf.Id
	for testParentId < h.tree.Size()-1 {
		parentId := testParentId
		parent, _ := h.tree.GetNode(parentId)
		left, leftErr := h.tree.LeftChild(parentId)
		right, rightErr := h.tree.RightChild(parentId)
		if leftErr == nil && h.Compare(left.Val, parent.Val) < 0 {
			parentId = left.Id
		}
		if rightErr == nil && h.Compare(right.Val, parent.Val) < 0 {
			parentId = right.Id
		}
		if parentId != testParentId {
			h.tree.Swap(parentId, testParentId)
			testParentId = parentId
		} else {
			break
		}
	}
	return root.Val, nil
}

func (h *Heap) Peek() (interface{}, error) {
	if h.tree.Size() == 0 {
		return nil, EOH
	} else {
		node, _ := h.tree.GetNode(0)
		return node.Val, nil
	}
}

func (h *Heap) Size() int {
	return h.tree.Size()
}
