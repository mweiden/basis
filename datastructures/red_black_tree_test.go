package datastructures

import (
	"reflect"
	"testing"
)

func TestRedBlackTreeRotation(t *testing.T) {
	t.Parallel()

	rotationTestNodes := func(forLeft bool) (*RedBlackTree, *RBNode, *RBNode, *RBNode, *RBNode, *RBNode, map[*RBNode]string) {
		tree := NewRedBlackTree()
		n := testNode()
		l := testNode()
		r := testNode()
		ll := testNode()
		rr := testNode()
		tree.root = n
		n.left = l
		l.parent = n
		n.right = r
		r.parent = n
		if forLeft {
			r.left = ll
			ll.parent = r
			r.right = rr
			rr.parent = r
		} else {
			l.left = ll
			ll.parent = l
			l.right = rr
			rr.parent = l
		}
		m := map[*RBNode]string{
			n:            "n",
			l:            "l",
			r:            "r",
			ll:           "ll",
			rr:           "rr",
			leafSentinel: "leafSentinel",
			nil:          "nil",
		}
		return tree, n, l, r, ll, rr, m
	}

	// rotate left
	tree, n, l, r, ll, rr, m := rotationTestNodes(true)
	tree.rotateLeft(n)

	npAssert(t, "rotateLeft, l.left", l.left, leafSentinel, m)
	npAssert(t, "rotateLeft, l.right", l.right, leafSentinel, m)
	npAssert(t, "rotateLeft, l.parent", l.parent, n, m)
	npAssert(t, "rotateLeft, ll.left", ll.left, leafSentinel, m)
	npAssert(t, "rotateLeft, ll.right", ll.right, leafSentinel, m)
	npAssert(t, "rotateLeft, ll.parent", ll.parent, n, m)
	npAssert(t, "rotateLeft, n.left", n.left, l, m)
	npAssert(t, "rotateLeft, n.right", n.right, ll, m)
	npAssert(t, "rotateLeft, n.parent", n.parent, r, m)
	npAssert(t, "rotateLeft, r.left", r.left, n, m)
	npAssert(t, "rotateLeft, r.right", r.right, rr, m)
	npAssert(t, "rotateLeft, r.parent", r.parent, nil, m)
	npAssert(t, "rotateLeft, rr.left", rr.left, leafSentinel, m)
	npAssert(t, "rotateLeft, rr.right", rr.right, leafSentinel, m)
	npAssert(t, "rotateLeft, rr.parent", rr.parent, r, m)

	// rotate right
	tree, n, l, r, ll, rr, m = rotationTestNodes(false)
	tree.rotateRight(n)

	npAssert(t, "rotateRight, l.left", l.left, ll, m)
	npAssert(t, "rotateRight, l.right", l.right, n, m)
	npAssert(t, "rotateRight, l.parent", l.parent, nil, m)
	npAssert(t, "rotateRight, ll.left", ll.left, leafSentinel, m)
	npAssert(t, "rotateRight, ll.right", ll.right, leafSentinel, m)
	npAssert(t, "rotateRight, ll.parent", ll.parent, l, m)
	npAssert(t, "rotateRight, n.left", n.left, rr, m)
	npAssert(t, "rotateRight, n.right", n.right, r, m)
	npAssert(t, "rotateRight, n.parent", n.parent, l, m)
	npAssert(t, "rotateRight, r.left", r.left, leafSentinel, m)
	npAssert(t, "rotateRight, r.right", r.right, leafSentinel, m)
	npAssert(t, "rotateRight, r.parent", r.parent, n, m)
	npAssert(t, "rotateRight, rr.left", rr.left, leafSentinel, m)
	npAssert(t, "rotateRight, rr.right", rr.right, leafSentinel, m)
	npAssert(t, "rotateRight, rr.parent", rr.parent, n, m)

}

func TestRedBlackTreeInsertion(t *testing.T) {
	t.Parallel()

	resetTestNodes := func() (*RedBlackTree, *RBNode, *RBNode, *RBNode, *RBNode, map[*RBNode]string) {
		n := testNode()
		p := testNode()
		g := testNode()
		u := testNode()
		tree := NewRedBlackTree()
		tree.insertCaseChain = false
		tree.root = g
		g.parent = nil
		g.left = p
		p.parent = g
		g.right = u
		u.parent = g
		n.parent = p
		m := map[*RBNode]string{
			n:            "n",
			p:            "p",
			g:            "g",
			u:            "u",
			leafSentinel: "leafSentinel",
			nil:          "nil",
		}
		return tree, n, p, g, u, m
	}

	// insert case 1 should color a node black if it has no parent
	tree, n, p, g, u, m := resetTestNodes()
	n.color = RED
	n.parent = n

	tree.insertCase1(n)
	expectedResult := RED
	result := n.color

	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}

	n.parent = nil
	tree.insertCase1(n)
	expectedResult = BLACK
	result = n.color

	if expectedResult != result {
		t.Errorf("%v != %v", expectedResult, result)
	}

	// insert case 2 does nothing

	// insert case 3
	tree, n, p, g, u, m = resetTestNodes()
	tree.insertCaseChain = false
	p.left = n
	p.color = RED
	n.color = RED
	u.color = RED

	tree.insertCase3(n)

	colorAssert(t, "insertCase3, g", g, RED)
	colorAssert(t, "insertCase3, p", p, BLACK)
	colorAssert(t, "insertCase3, u", u, BLACK)
	colorAssert(t, "insertCase3, n", n, RED)
	npAssert(t, "insertCase3, g.left", g.left, p, m)
	npAssert(t, "insertCase3, g.right", g.right, u, m)
	npAssert(t, "insertCase3, n.parent", n.parent, p, m)
	npAssert(t, "insertCase3, u.parent", u.parent, g, m)
	npAssert(t, "insertCase3, n.left", n.left, leafSentinel, m)
	npAssert(t, "insertCase3, p.parent", p.parent, g, m)
	npAssert(t, "insertCase3, p.left", p.left, n, m)
	npAssert(t, "insertCase3, p.right", p.right, leafSentinel, m)
	npAssert(t, "insertCase3, n.right", n.right, leafSentinel, m)
	npAssert(t, "insertCase3, u.left", u.left, leafSentinel, m)
	npAssert(t, "insertCase3, u.right", u.right, leafSentinel, m)

	// insert case 4
	tree, n, p, g, u, m = resetTestNodes()
	tree.insertCaseChain = false
	p.right = n
	p.color = RED
	n.color = RED

	tree.insertCase4(n)

	colorAssert(t, "insertCase4, g", g, BLACK)
	colorAssert(t, "insertCase4, p", p, RED)
	colorAssert(t, "insertCase4, u", u, BLACK)
	colorAssert(t, "insertCase4, n", n, RED)
	npAssert(t, "insertCase4, g.left", g.left, n, m)
	npAssert(t, "insertCase4, g.right", g.right, u, m)
	npAssert(t, "insertCase4, n.parent", n.parent, g, m)
	npAssert(t, "insertCase4, u.parent", u.parent, g, m)
	npAssert(t, "insertCase4, p.parent", p.parent, n, m)
	npAssert(t, "insertCase4, p.left", p.left, leafSentinel, m)
	npAssert(t, "insertCase4, p.right", p.right, leafSentinel, m)
	npAssert(t, "insertCase4, n.right", n.right, leafSentinel, m)
	npAssert(t, "insertCase4, n.left", n.left, p, m)
	npAssert(t, "insertCase4, u.left", u.left, leafSentinel, m)
	npAssert(t, "insertCase4, u.right", u.right, leafSentinel, m)

	// case 4 step 2
	tree, n, p, g, u, m = resetTestNodes()
	tree.insertCaseChain = false
	p.left = n
	p.color = RED
	n.color = RED

	tree.insertCase4Step2(n)

	colorAssert(t, "insertCase4Step2, g", g, RED)
	colorAssert(t, "insertCase4Step2, p", p, BLACK)
	colorAssert(t, "insertCase4Step2, u", u, BLACK)
	colorAssert(t, "insertCase4Step2, n", n, RED)
	npAssert(t, "insertCase4Step2, g.left", g.left, leafSentinel, m)
	npAssert(t, "insertCase4Step2, g.right", g.right, u, m)
	npAssert(t, "insertCase4Step2, n.parent", n.parent, p, m)
	npAssert(t, "insertCase4Step2, u.parent", u.parent, g, m)
	npAssert(t, "insertCase4Step2, p.parent", p.parent, nil, m)
	npAssert(t, "insertCase4Step2, p.left", p.left, n, m)
	npAssert(t, "insertCase4Step2, p.right", p.right, g, m)
	npAssert(t, "insertCase4Step2, n.right", n.right, leafSentinel, m)
	npAssert(t, "insertCase4Step2, n.left", n.left, leafSentinel, m)
	npAssert(t, "insertCase4Step2, u.left", u.left, leafSentinel, m)
	npAssert(t, "insertCase4Step2, u.right", u.right, leafSentinel, m)

	// test full insert
	tree = NewRedBlackTree()

	tree.Insert(NumNode{10}, 10)
	tree.Insert(NumNode{5}, 5)
	tree.Insert(NumNode{25}, 25)
	tree.Insert(NumNode{20}, 20)
	tree.Insert(NumNode{15}, 15)
	tree.Insert(NumNode{13}, 13)
	tree.Insert(NumNode{12}, 12)
	tree.Insert(NumNode{11}, 11)

	expectedInsertTraversal := []interface{}{5, 10, 11, 12, 13, 15, 20, 25}
	traversal := tree.Traverse()
	if !reflect.DeepEqual(expectedInsertTraversal, traversal) {
		t.Errorf("%v != %v", traversal, expectedInsertTraversal)
	}
}

func TestRedBlackTreeDeletion(t *testing.T) {
	t.Parallel()

	resetTestNodes := func() (*RedBlackTree, *RBNode, *RBNode, *RBNode, *RBNode, *RBNode, map[*RBNode]string) {
		n := testNode()
		p := testNode()
		s := testNode()
		sl := testNode()
		sr := testNode()
		tree := NewRedBlackTree()
		tree.deleteCaseChain = false
		tree.root = p
		p.parent = nil
		p.left = n
		n.parent = p
		p.right = s
		s.parent = p
		s.left = sl
		sl.parent = s
		s.right = sr
		sr.parent = s
		m := map[*RBNode]string{
			n:            "n",
			p:            "p",
			s:            "s",
			sl:           "sl",
			sr:           "sr",
			leafSentinel: "leafSentinel",
			nil:          "nil",
		}
		return tree, n, p, s, sl, sr, m
	}

	// delete case 1 is trivial

	// delete case 2
	tree, n, p, s, sl, sr, m := resetTestNodes()
	s.color = RED

	tree.deleteCase2(n)

	colorAssert(t, "deleteCase2, s", s, BLACK)
	colorAssert(t, "deleteCase2, n", n, BLACK)
	colorAssert(t, "deleteCase2, p", p, RED)
	colorAssert(t, "deleteCase2, sl", sl, BLACK)
	colorAssert(t, "deleteCase2, sr", sr, BLACK)
	npAssert(t, "deleteCase2, n.parent", n.parent, p, m)
	npAssert(t, "deleteCase2, n.right", n.right, leafSentinel, m)
	npAssert(t, "deleteCase2, n.left", n.left, leafSentinel, m)
	npAssert(t, "deleteCase2, p.parent", p.parent, s, m)
	npAssert(t, "deleteCase2, p.left", p.left, n, m)
	npAssert(t, "deleteCase2, p.right", p.right, sl, m)
	npAssert(t, "deleteCase2, s.parent", s.parent, nil, m)
	npAssert(t, "deleteCase2, s.left", s.left, p, m)
	npAssert(t, "deleteCase2, s.right", s.right, sr, m)

	// delete case 3
	tree, n, p, s, sl, sr, m = resetTestNodes()

	tree.deleteCase3(n)

	colorAssert(t, "deleteCase3, s", s, RED)
	colorAssert(t, "deleteCase3, n", n, BLACK)
	colorAssert(t, "deleteCase3, p", p, BLACK)
	colorAssert(t, "deleteCase3, sl", sl, BLACK)
	colorAssert(t, "deleteCase3, sr", sr, BLACK)
	npAssert(t, "deleteCase3, n.parent", n.parent, p, m)
	npAssert(t, "deleteCase3, n.right", n.right, leafSentinel, m)
	npAssert(t, "deleteCase3, n.left", n.left, leafSentinel, m)
	npAssert(t, "deleteCase3, p.parent", p.parent, nil, m)
	npAssert(t, "deleteCase3, p.left", p.left, n, m)
	npAssert(t, "deleteCase3, p.right", p.right, s, m)
	npAssert(t, "deleteCase3, s.parent", s.parent, p, m)
	npAssert(t, "deleteCase3, s.left", s.left, sl, m)
	npAssert(t, "deleteCase3, s.right", s.right, sr, m)

	// delete case 4
	tree, n, p, s, sl, sr, m = resetTestNodes()
	p.color = RED

	tree.deleteCase4(n)

	colorAssert(t, "deleteCase4, s", s, RED)
	colorAssert(t, "deleteCase4, n", n, BLACK)
	colorAssert(t, "deleteCase4, p", p, BLACK)
	colorAssert(t, "deleteCase4, sl", sl, BLACK)
	colorAssert(t, "deleteCase4, sr", sr, BLACK)
	npAssert(t, "deleteCase4, n.parent", n.parent, p, m)
	npAssert(t, "deleteCase4, n.right", n.right, leafSentinel, m)
	npAssert(t, "deleteCase4, n.left", n.left, leafSentinel, m)
	npAssert(t, "deleteCase4, p.parent", p.parent, nil, m)
	npAssert(t, "deleteCase4, p.left", p.left, n, m)
	npAssert(t, "deleteCase4, p.right", p.right, s, m)
	npAssert(t, "deleteCase4, s.parent", s.parent, p, m)
	npAssert(t, "deleteCase4, s.left", s.left, sl, m)
	npAssert(t, "deleteCase4, s.right", s.right, sr, m)

	// delete case 5
	tree, n, p, s, sl, sr, m = resetTestNodes()
	sl.color = RED

	tree.deleteCase5(n)

	colorAssert(t, "deleteCase5, s", s, RED)
	colorAssert(t, "deleteCase5, sl", sl, BLACK)
	colorAssert(t, "deleteCase5, sr", sr, BLACK)
	npAssert(t, "deleteCase5, s.parent", s.parent, sl, m)
	npAssert(t, "deleteCase5, s.left", s.left, leafSentinel, m)
	npAssert(t, "deleteCase5, s.right", s.right, sr, m)
	npAssert(t, "deleteCase5, sl.parent", sl.parent, p, m)
	npAssert(t, "deleteCase5, sl.left", sl.left, leafSentinel, m)
	npAssert(t, "deleteCase5, sl.right", sl.right, s, m)
	npAssert(t, "deleteCase5, sr.parent", sr.parent, s, m)
	npAssert(t, "deleteCase5, sr.left", sr.left, leafSentinel, m)
	npAssert(t, "deleteCase5, sr.right", sr.right, leafSentinel, m)

	// delete case 6
	tree, n, p, s, sl, sr, m = resetTestNodes()
	sr.color = RED
	s.left = leafSentinel

	tree.deleteCase6(n)

	colorAssert(t, "deleteCase6, p", p, BLACK)
	colorAssert(t, "deleteCase6, sr", sr, BLACK)
	colorAssert(t, "deleteCase6, n", n, BLACK)
	npAssert(t, "deleteCase6, s.parent", s.parent, nil, m)
	npAssert(t, "deleteCase6, s.left", s.left, p, m)
	npAssert(t, "deleteCase6, s.right", s.right, sr, m)
	npAssert(t, "deleteCase6, sr.parent", sr.parent, s, m)
	npAssert(t, "deleteCase6, sr.left", sr.left, leafSentinel, m)
	npAssert(t, "deleteCase6, sr.right", sr.right, leafSentinel, m)
	npAssert(t, "deleteCase6, p.parent", p.parent, s, m)
	npAssert(t, "deleteCase6, p.left", p.left, n, m)
	npAssert(t, "deleteCase6, p.right", p.right, leafSentinel, m)

	// test deletion
	tree = NewRedBlackTree()

	tree.Insert(NumNode{10}, 10)
	tree.Insert(NumNode{5}, 5)
	tree.Insert(NumNode{25}, 25)
	tree.Insert(NumNode{20}, 20)
	tree.Insert(NumNode{15}, 15)
	tree.Insert(NumNode{13}, 13)
	tree.Insert(NumNode{12}, 12)
	tree.Insert(NumNode{11}, 11)
	tree.Delete(NumNode{20})

	expectedInsertTraversal := []interface{}{5, 10, 11, 12, 13, 15, 25}
	traversal := tree.Traverse()
	if !reflect.DeepEqual(expectedInsertTraversal, traversal) {
		t.Errorf("%v != %v", traversal, expectedInsertTraversal)
	}

	tree.Delete(NumNode{12})
	expectedInsertTraversal = []interface{}{5, 10, 11, 13, 15, 25}
	traversal = tree.Traverse()
	if !reflect.DeepEqual(expectedInsertTraversal, traversal) {
		t.Errorf("%v != %v", traversal, expectedInsertTraversal)
	}

}

type NumNode struct {
	i int
}

func (n NumNode) Compare(other interface{}) int {
	return n.i - other.(NumNode).i
}

func testNode() *RBNode {
	return &RBNode{Key: sentinelKey, Value: nil, color: BLACK, left: leafSentinel, right: leafSentinel}
}

func npAssert(t *testing.T, name string, p1 *RBNode, p2 *RBNode, m map[*RBNode]string) {
	if p1 != p2 {
		t.Errorf("Node pointer assertion failure (%s): %s != %s", name, m[p1], m[p2])
	}
}

func colorAssert(t *testing.T, name string, n *RBNode, c2 int) {
	if n.color != c2 {
		t.Errorf("Color assertion failure (%s): %v != %v", name, n.color, c2)
	}
}
