package main

type Color int 

const (
	Red Color = 1 
	Black Color = 2
)

type Node struct {
	parent *Node
	left *Node
	right *Node
	value int
	color Color
}

type RBTree struct {
	root *Node
}

func (t *RBTree) Insert(value int) {
	newNode := &Node{
		value: value,
	}
	if t.root == nil {
		newNode.color = Red
		t.root = newNode
	} else {
		node := t.root
		loop := true 
		for loop {
			if node.value > value {
				
			}
		}

	}


}



