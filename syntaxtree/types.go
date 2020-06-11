package syntaxtree

// Node is an interface that must be implemented by every node in the tree
type Node interface {
	GetTokenLiteral() string
}

// Statement is a type of node which provides statement functionality.
// Statements do not produce values. E.g: bind a value to a name: let x = 6;
type Stmt interface {
	Node
	// stmtNode() ensures that only statement nodes can be
	// assigned to a Stmt.
	stmtNode()
}

// Expression is a type of node which provides expression functionality
// Expressions produce values. E.g: 6; sum(6,6)
type Expr interface {
	Node
	// expressionNode() ensures that only expression/type nodes can be
	// assigned to an Expr.
	exprNode()
}

