package illygen

import (
	"fmt"

	"github.com/leraniode/illygen/internal/graph"
)

// Flow is a net of connected nodes — the reasoning pipeline.
// Internally it is a directed weighted graph: nodes are vertices,
// Links are edges. The engine walks this graph during execution.
//
// Build a flow with the fluent API:
//
//	flow := illygen.NewFlow().
//	    Add(inputNode).
//	    Add(outputNode).
//	    Link("input", "output", 1.0)
//
// The first node passed to Add becomes the entry point unless
// overridden by Entry.
type Flow struct {
	nodes map[string]*Node
	graph *graph.Graph
	entry string
}

// NewFlow creates a new empty Flow.
func NewFlow() *Flow {
	return &Flow{
		nodes: make(map[string]*Node),
		graph: graph.New(),
	}
}

// Add registers a node into the flow.
// The first node added becomes the entry point automatically.
// Passing a nil node panics immediately with a clear message.
// Returns the Flow for chaining.
func (f *Flow) Add(node *Node) *Flow {
	if node == nil {
		panic("illygen: Flow.Add called with nil node")
	}
	f.nodes[node.ID()] = node
	if f.entry == "" {
		f.entry = node.ID()
	}
	return f
}

// Link connects two nodes with a weight (0.0 to 1.0).
// When a node's Result does not specify a Next, the engine follows
// the highest-weight link automatically.
// Duplicate links are silently ignored — the first call wins.
// Returns the Flow for chaining.
func (f *Flow) Link(from, to string, weight float64) *Flow {
	if err := f.graph.Add(from, to, weight); err != nil {
		// duplicate edge — silently skip for fluent API usability
		_ = err
	}
	return f
}

// Entry explicitly sets which node the flow starts from.
// By default the first node passed to Add is the entry.
// Use Entry when you add nodes out of order or want to change
// the starting point after building the flow.
func (f *Flow) Entry(nodeID string) *Flow {
	f.entry = nodeID
	return f
}

// node retrieves a node by ID. Returns an error if not found.
func (f *Flow) node(id string) (*Node, error) {
	n, ok := f.nodes[id]
	if !ok {
		return nil, fmt.Errorf("illygen: node %q not found in flow", id)
	}
	return n, nil
}

// entryNode returns the entry node. Returns an error if no entry is defined.
func (f *Flow) entryNode() (*Node, error) {
	if f.entry == "" {
		return nil, fmt.Errorf("illygen: flow has no entry node — call Add() first")
	}
	return f.node(f.entry)
}
