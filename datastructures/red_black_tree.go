package datastructures

import (
	"math/rand"
)

const (
	RED = iota
	BLACK
)

type RBNode struct {
	Key    Comparable
	Value  interface{}
	color  int
	left   *RBNode
	right  *RBNode
	parent *RBNode
}

type SentinelKey struct{}

func (k SentinelKey) Compare(interface{}) int {
	return -1
}

var sentinelKey SentinelKey = SentinelKey{}
var leafSentinel *RBNode = &RBNode{Key: sentinelKey, Value: nil, color: BLACK}

type RedBlackTree struct {
	root            *RBNode
	insertCaseChain bool
	deleteCaseChain bool
}

func NewRedBlackTree() *RedBlackTree {
	return &RedBlackTree{root: nil, insertCaseChain: true, deleteCaseChain: true}
}

func (t *RedBlackTree) Delete(key Comparable) {
	n := t.find(key)
	if n == nil {
		return
	}
	s := rand.Float32()
	t.deleteNode(n, s)
}

func (t *RedBlackTree) deleteNode(n *RBNode, s float32) {
	if hasTwoLeafChildren(n) {
		if n.parent.right == n {
			n.parent.right = leafSentinel
		} else {
			n.parent.left = leafSentinel
		}
	} else if isInternal(n) {
		var m *RBNode
		if s < 0.5 {
			// select from left subtree
			m = maxOfSubtree(n.left)
		} else {
			// select from right subtree
			m = minOfSubtree(n.right)
		}
		n.Value = m.Value
		t.deleteNode(m, s)
	} else {
		// has at least one child
		if t.deleteCaseChain {
			t.deleteOneChild(n)
		}
	}
}

func (t *RedBlackTree) deleteOneChild(n *RBNode) {
	child := n.right
	if isLeaf(n.right) {
		child = n.left
	}
	t.replaceNode(n, child)
	if n.color == BLACK {
		if child.color == RED {
			child.color = BLACK
		} else if t.deleteCaseChain {
			t.deleteCase1(child)
		}
	}
}

func (t *RedBlackTree) replaceNode(dest *RBNode, src *RBNode) {
	dest.parent = src.parent
	dest.left = src.left
	dest.right = src.right
	dest.Value = src.Value
	dest.Key = src.Key
	dest.color = src.color
}

func (t *RedBlackTree) deleteCase1(n *RBNode) {
	if n.parent == nil && t.deleteCaseChain {
		t.deleteCase2(n)
	}
}

func (t *RedBlackTree) deleteCase2(n *RBNode) {
	s := sibling(n)
	if s.color == RED {
		n.parent.color = RED
		s.color = BLACK
		if n == n.parent.left {
			t.rotateLeft(n.parent)
		} else {
			t.rotateRight(n.parent)
		}
	}
	if t.deleteCaseChain {
		t.deleteCase3(n)
	}
}

func (t *RedBlackTree) deleteCase3(n *RBNode) {
	s := sibling(n)
	if n.parent.color == BLACK &&
		s.color == BLACK &&
		s.left.color == BLACK &&
		s.right.color == BLACK {
		s.color = RED
		if t.deleteCaseChain {
			t.deleteCase1(n.parent)
		}
	} else if t.deleteCaseChain {
		t.deleteCase4(n)
	}
}

func (t *RedBlackTree) deleteCase4(n *RBNode) {
	s := sibling(n)
	if n.parent.color == RED &&
		s.color == BLACK &&
		s.left.color == BLACK &&
		s.right.color == BLACK {
		s.color = RED
		n.parent.color = BLACK
	} else if t.deleteCaseChain {
		t.deleteCase5(n)
	}
}

func (t *RedBlackTree) deleteCase5(n *RBNode) {
	s := sibling(n)
	if s.color == BLACK {
		if n == n.parent.left &&
			s.right.color == BLACK &&
			s.left.color == RED {
			s.color = RED
			s.left.color = BLACK
			t.rotateRight(s)
		} else if n.parent.right == n &&
			s.left.color == BLACK &&
			s.right.color == RED {
			s.color = RED
			s.right.color = BLACK
			t.rotateLeft(s)
		}
	}
	if t.deleteCaseChain {
		t.deleteCase6(n)
	}
}

func (t *RedBlackTree) deleteCase6(n *RBNode) {
	s := sibling(n)
	s.color = n.parent.color
	n.parent.color = BLACK

	if n == n.parent.left {
		s.right.color = BLACK
		t.rotateLeft(n.parent)
	} else {
		s.left.color = BLACK
		t.rotateRight(n.parent)
	}
}

func isLeaf(n *RBNode) bool {
	return leafSentinel == n
}

func hasTwoLeafChildren(n *RBNode) bool {
	return n.left == leafSentinel && n.right == leafSentinel
}

func isInternal(n *RBNode) bool {
	return n.left != leafSentinel && n.right != leafSentinel
}

func (t *RedBlackTree) find(key Comparable) *RBNode {
	n := t.root
	for n != leafSentinel && n.Key != key {
		if key.Compare(n.Key) < 0 {
			n = n.left
		} else {
			n = n.right
		}
	}
	if n == leafSentinel {
		return nil
	}
	return n
}

func minOfSubtree(n *RBNode) *RBNode {
	m := n
	for m.left != leafSentinel {
		m = m.left
	}
	return m
}

func maxOfSubtree(n *RBNode) *RBNode {
	m := n
	for m.right != leafSentinel {
		m = m.right
	}
	return m
}

func (t *RedBlackTree) Insert(key Comparable, value interface{}) {
	n := &RBNode{Key: key, Value: value} // TODO: color
	t.insertRecurse(t.root, n)
	t.insertRepairTree(n)
	t.root = n
	for t.root.parent != nil {
		t.root = t.root.parent
	}
}

func (t *RedBlackTree) insertRecurse(root *RBNode, n *RBNode) {
	if root != nil && n.Key.Compare(root.Key) < 0 {
		if root.left != leafSentinel {
			t.insertRecurse(root.left, n)
			return
		}
		root.left = n
	} else if root != nil {
		if root.right != leafSentinel {
			t.insertRecurse(root.right, n)
			return
		}
		root.right = n
	}
	n.parent = root
	n.left = leafSentinel
	n.right = leafSentinel
	n.color = RED
}

func (t *RedBlackTree) insertRepairTree(n *RBNode) {
	if n.parent == nil {
		t.insertCase1(n)
	} else if n.parent.color == BLACK {
		// do nothing in case 2
	} else if uncle(n).color == RED {
		t.insertCase3(n)
	} else {
		t.insertCase4(n)
	}
}

func (t *RedBlackTree) insertCase1(n *RBNode) {
	if n.parent == nil {
		n.color = BLACK
	}
}

func (t *RedBlackTree) insertCase3(n *RBNode) {
	n.parent.color = BLACK
	uncle(n).color = BLACK
	grandparent(n).color = RED
	if t.insertCaseChain {
		t.insertRepairTree(grandparent(n))
	}
}

func (t *RedBlackTree) insertCase4(n *RBNode) {
	p := n.parent
	g := grandparent(n)

	if n == g.left.right {
		t.rotateLeft(p)
		n = n.left
	} else if n == g.right.left {
		t.rotateRight(p)
		n = n.right
	}
	if t.insertCaseChain {
		t.insertCase4Step2(n)
	}
}

func (t *RedBlackTree) insertCase4Step2(n *RBNode) {
	p := n.parent
	g := grandparent(n)

	if n == p.left {
		t.rotateRight(g)
	} else {
		t.rotateLeft(g)
	}
	p.color = BLACK
	g.color = RED
}

func (t *RedBlackTree) rotateRight(n *RBNode) {
	newN := n.left
	if newN == leafSentinel {
		panic(1)
	}
	n.left = newN.right
	n.left.parent = n
	newN.right = n
	newN.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = newN
		} else {
			n.parent.right = newN
		}
	}
	n.parent = newN
}

func (t *RedBlackTree) rotateLeft(n *RBNode) {
	newN := n.right
	if newN == leafSentinel {
		panic(1) // since the leaves of a red-black tree are empty, they cannot become internal nodes
	}
	n.right = newN.left
	n.right.parent = n
	newN.left = n
	newN.parent = n.parent
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = newN
		} else {
			n.parent.right = newN
		}
	}
	n.parent = newN
}

func grandparent(n *RBNode) *RBNode {
	p := n.parent
	if p == nil {
		return nil
	}
	return p.parent
}

func sibling(n *RBNode) *RBNode {
	p := n.parent
	if n == p.left {
		return p.right
	}
	return p.left
}

func uncle(n *RBNode) *RBNode {
	p := n.parent
	g := grandparent(n)
	if g == nil {
		return nil
	}
	return sibling(p)
}

func (t *RedBlackTree) Traverse() []interface{} {
	return inOrderTraverse(t.root)
}

func inOrderTraverse(n *RBNode) []interface{} {
	if n == leafSentinel {
		return nil
	}
	var values []interface{}
	values = append(values, inOrderTraverse(n.left)...)
	values = append(values, n.Value)
	values = append(values, inOrderTraverse(n.right)...)
	return values
}
