// Package graph provides raw graph primitives for flow management.
// This is an internal package — use the core/flow API instead.
package graph

import (
	"fmt"
	"sync"
)

// Connection represents a directed weighted edge between two nodes.
type Connection struct {
	From   string
	To     string
	Weight float64 // 0.0 to 1.0
}

// Graph is a directed weighted graph of node connections.
// It is the underlying structure of a Flow.
type Graph struct {
	mu          sync.RWMutex
	connections map[string][]*Connection // keyed by From node ID
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{
		connections: make(map[string][]*Connection),
	}
}

// Connect adds a directed connection from → to with the given weight.
func (g *Graph) Connect(from, to string, weight float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.connections[from] {
		if c.To == to {
			return fmt.Errorf("illygen/graph: connection %q → %q already exists", from, to)
		}
	}
	g.connections[from] = append(g.connections[from], &Connection{
		From:   from,
		To:     to,
		Weight: weight,
	})
	return nil
}

// Outgoing returns all connections from a given node, sorted by weight descending.
func (g *Graph) Outgoing(from string) []*Connection {
	g.mu.RLock()
	defer g.mu.RUnlock()

	conns := make([]*Connection, len(g.connections[from]))
	copy(conns, g.connections[from])
	sortConnections(conns)
	return conns
}

// AdjustWeight modifies the weight of a specific connection.
// delta is added to the current weight and the result is clamped to [0.0, 1.0].
func (g *Graph) AdjustWeight(from, to string, delta float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, c := range g.connections[from] {
		if c.To == to {
			c.Weight = clamp(c.Weight+delta, 0.0, 1.0)
			return nil
		}
	}
	return fmt.Errorf("illygen/graph: connection %q → %q not found", from, to)
}

// Disconnect removes a connection from → to.
func (g *Graph) Disconnect(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conns := g.connections[from]
	for i, c := range conns {
		if c.To == to {
			g.connections[from] = append(conns[:i], conns[i+1:]...)
			return
		}
	}
}

func sortConnections(conns []*Connection) {
	for i := 1; i < len(conns); i++ {
		for j := i; j > 0 && conns[j].Weight > conns[j-1].Weight; j-- {
			conns[j], conns[j-1] = conns[j-1], conns[j]
		}
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
