// Package runtime contains the internal execution engine for Illygen.
// Users interact with Engine in the illygen package — not this directly.
package runtime

import "fmt"

// maxVisits is the maximum number of times a single node can be visited
// in one flow execution before the engine declares a cycle and returns an error.
// This guards against infinite loops caused by circular Next routing.
const maxVisits = 50

// Step records what happened at a single node during execution.
type Step struct {
	NodeID     string
	Value      any
	Confidence float64
	Next       string
}

// ExecutionTrace is the complete record of a flow execution.
// Steps holds every node visited in order.
// Final holds the last step — its Value and Confidence are returned to the caller.
// In future versions, the learning logic will read the trace to adjust weights.
type ExecutionTrace struct {
	Steps []Step
	Final Step
	Done  bool
}

// NodeExecutor is the function the engine passes to Execute.
// It decouples the executor from illygen package types, keeping internal/ clean.
type NodeExecutor func(nodeID string) (value any, confidence float64, next string, err error)

// Execute runs the flow starting from entry, calling executor for each node,
// and following the returned next node until execution ends.
//
// Algorithm:
//
//	current = entry
//	loop:
//	  (value, confidence, next) = executor(current)
//	  record step
//	  current = next
//	until current == ""
func Execute(entry string, executor NodeExecutor) (*ExecutionTrace, error) {
	trace := &ExecutionTrace{}
	current := entry
	visited := make(map[string]int)

	for current != "" {
		visited[current]++
		if visited[current] > maxVisits {
			return nil, fmt.Errorf(
				"illygen/runtime: node %q visited %d times — possible cycle detected",
				current, visited[current],
			)
		}

		value, confidence, next, err := executor(current)
		if err != nil {
			return nil, fmt.Errorf("illygen/runtime: node %q failed: %w", current, err)
		}

		step := Step{
			NodeID:     current,
			Value:      value,
			Confidence: confidence,
			Next:       next,
		}
		trace.Steps = append(trace.Steps, step)
		trace.Final = step

		current = next
	}

	trace.Done = true
	return trace, nil
}
