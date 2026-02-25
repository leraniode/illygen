// Package core defines the fundamental building blocks of Illygen.
package core

// Node is the atomic unit of reasoning in Illygen.
// A node is consulted during flow execution — it receives context,
// evaluates against its knowledge, and returns a Verdict.
//
// Think of it as a neuron: it fires, produces a signal, and passes it forward.
type Node interface {
	// ID returns the unique identifier of this node within a flow.
	ID() string

	// Consult is called by the flow engine during execution.
	// The node inspects the context, applies its knowledge, and returns a Verdict.
	Consult(ctx *Context) (Verdict, error)
}

// Verdict is the output of a node's consultation.
// It tells the flow engine where to go next and how confident the node is.
type Verdict struct {
	// Route is the ID of the next node to consult.
	// An empty Route signals the end of this reasoning path.
	Route string

	// Output is whatever this node concluded — can be any value.
	Output any

	// Weight represents the node's confidence in this verdict (0.0 to 1.0).
	// The learning logic uses this to adjust connection strengths over time.
	Weight float64
}

// BaseNode provides a default implementation of the Node interface.
// Embed this in your own nodes to avoid boilerplate.
type BaseNode struct {
	id string
}

// NewBaseNode creates a new BaseNode with the given ID.
func NewBaseNode(id string) BaseNode {
	return BaseNode{id: id}
}

// ID returns the node's identifier.
func (b BaseNode) ID() string {
	return b.id
}
