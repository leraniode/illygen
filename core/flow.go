package core

import (
	"fmt"

	"github.com/leraniode/illygen/internal/graph"
)

// Flow is a net of connected nodes — the reasoning pipeline.
// Like a neural network, it is not static: it reshapes itself as learning happens.
type Flow struct {
	id    string
	nodes map[string]Node
	graph *graph.Graph
	entry string // the first node to consult
}

// NewFlow creates a new, empty flow with the given ID.
func NewFlow(id string) *Flow {
	return &Flow{
		id:    id,
		nodes: make(map[string]Node),
		graph: graph.New(),
	}
}

// ID returns the flow's identifier.
func (f *Flow) ID() string {
	return f.id
}

// Add registers a node into the flow.
// The first node added becomes the entry point.
func (f *Flow) Add(node Node) *Flow {
	f.nodes[node.ID()] = node
	if f.entry == "" {
		f.entry = node.ID()
	}
	return f
}

// Entry sets the entry node explicitly.
// Useful when you want a specific node to be consulted first.
func (f *Flow) Entry(nodeID string) *Flow {
	f.entry = nodeID
	return f
}

// Connect wires two nodes together with a default weight of 1.0.
func (f *Flow) Connect(from, to string) *Flow {
	return f.ConnectWeighted(from, to, 1.0)
}

// ConnectWeighted wires two nodes with an explicit weight.
func (f *Flow) ConnectWeighted(from, to string, weight float64) *Flow {
	if err := f.graph.Connect(from, to, weight); err != nil {
		// connection already exists — silently skip
		_ = err
	}
	return f
}

// Execute runs the flow starting from the entry node.
// It walks the graph, consulting each node in sequence,
// following the routes returned in each Verdict.
//
// Execute is synchronous. The runtime wraps this in goroutines.
func (f *Flow) Execute(ctx *Context) (*FlowOutput, error) {
	if f.entry == "" {
		return nil, fmt.Errorf("illygen/flow %q: no entry node defined", f.id)
	}

	output := newFlowOutput(f.id, ctx.TraceID)
	current := f.entry

	for current != "" {
		node, ok := f.nodes[current]
		if !ok {
			return nil, fmt.Errorf("illygen/flow %q: node %q not found", f.id, current)
		}

		verdict, err := node.Consult(ctx)
		if err != nil {
			return nil, fmt.Errorf("illygen/flow %q: node %q: %w", f.id, current, err)
		}

		output.record(current, verdict)

		// Follow the route specified in the verdict.
		// If the route is empty, the flow ends here.
		if verdict.Route == "" {
			break
		}

		// Validate that the route exists as a connection in the graph.
		found := false
		for _, conn := range f.graph.Outgoing(current) {
			if conn.To == verdict.Route {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("illygen/flow %q: node %q returned unknown route %q",
				f.id, current, verdict.Route)
		}

		current = verdict.Route
	}

	output.complete()
	return output, nil
}

// Graph exposes the underlying graph for the learning logic to adjust.
// This is intentionally low-level — use the learning package instead.
func (f *Flow) Graph() *graph.Graph {
	return f.graph
}
