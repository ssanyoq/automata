package main

type Node interface {
	// Returns list of all subsidiary nodes
	GetChildren() []*Node
}

type AST struct {
	Root *Node
}

// Node for representing binary operations such as concatenation
type BinaryOpNode struct {
	Node
	Left      Node
	Right     Node
	Operation Token
}

func (n *BinaryOpNode) GetChildren() []*Node {
	return []*Node{&n.Left, &n.Right}
}

// Node for representing unary operations such as kleene closure
type UnaryOpNode struct {
	Node
	Child     Node
	Operation Token
}

func (n *UnaryOpNode) GetChildren() []*Node {
	return []*Node{&n.Child}
}

// Node for representing capture groups that can be made using brackets
type CaptureGroupNode struct {
	Node
	Number int
	Child  Node
}

func (n *CaptureGroupNode) GetChildren() []*Node {
	return []*Node{&n.Child}
}

// Node for representing r{x,y} -type expressions
type RangeRepeatNode struct {
	Node
	From  *int
	To    *int
	Child Node
}

func (n *RangeRepeatNode) GetChildren() []*Node {
	return []*Node{&n.Child}
}

// Node for representing [a-z] -type expressions
type CharacterRangeNode struct {
	Node
	From *rune
	To   *rune
}

func (n *CharacterRangeNode) GetChildren() []*Node {
	return nil
}

// Node for representing atomic characters
type CharNode struct {
	Node
	Character rune
}

func (n *CharNode) GetChildren() []*Node {
	return nil
}
