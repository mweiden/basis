package binary_tree_diameter

import (
	"github.com/mweiden/basis/datastructures"
	"testing"
)

func TestBinaryTreeDiameter(t *testing.T) {
	t.Parallel()
	/*
	       0
	      / \
	     1   2
	    /
	   3

	*/
	tree := datastructures.BinaryTree{}
	zero := datastructures.Node{Id: 0, Val: nil}
	tree.SetNode(zero)
	one, _ := tree.LeftChild(0)
	two, _ := tree.RightChild(0)
	tree.SetNode(one)
	tree.SetNode(two)
	three, _ := tree.LeftChild(one.Id)
	tree.SetNode(three)

	expected := 3
	result := tree.Diameter()
	if expected != result {
		t.Errorf("%v != %v\n", expected, result)
	}

	/*
	 	 0
	        / \
	       1   2
	      / \
	     3	4
	    /	 \
	   7	 10

	*/
	four, _ := tree.RightChild(one.Id)
	tree.SetNode(four)
	seven, _ := tree.LeftChild(three.Id)
	tree.SetNode(seven)
	ten, _ := tree.RightChild(four.Id)
	tree.SetNode(ten)

	expected = 4
	result = tree.Diameter()
	if expected != result {
		t.Errorf("%v != %v\n", expected, result)
	}
}
