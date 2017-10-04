package datastructures

import (
	"fmt"
)

type ListNode struct {
	next *ListNode
	val  int
}

type LinkedList struct {
	start *ListNode
	end   *ListNode
}

func NewLinkedList() LinkedList {
	ll := LinkedList{}
	ll.start = nil
	ll.end = nil
	return ll
}

func (ll *LinkedList) Prepend(val int) {
	newStart := ListNode{nil, val}
	if ll.start != nil {
		newStart.next = ll.start
	}
	ll.start = &newStart
}

func (ll *LinkedList) Append(val int) {
	newEnd := ListNode{nil, val}
	if ll.start == nil {
		ll.start = &newEnd
	}
	if ll.end != nil {
		ll.end.next = &newEnd
	}
	ll.end = &newEnd
}

func (ll *LinkedList) Reverse() {
	// reverse the oder of the list
	var left *ListNode = nil
	right := ll.start

	for right != nil {
		tmp := right.next
		right.next = left
		left = right
		right = tmp
	}
	// reverse the end pointers
	tmp := ll.start
	ll.start = ll.end
	ll.end = tmp
}

func (ll *LinkedList) ToString() string {
	listString := ""
	printOp := func(currentNode *ListNode) {
		var pos string
		if ll.start == currentNode {
			pos = "s"
		} else if ll.end == currentNode {
			pos = "e"
		}
		listString += fmt.Sprintf("[%d]%s -> ", currentNode.val, pos)
	}
	ll.Walk(printOp)
	listString += "nil"
	return listString
}

func (ll *LinkedList) Walk(op func(*ListNode)) {
	var currentNode *ListNode = ll.start
	for (ll.start == nil && ll.end == nil) || currentNode != nil {
		op(currentNode)
		currentNode = currentNode.next
	}
}
