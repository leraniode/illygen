package core

import "time"

// FlowOutput holds the result of a complete flow execution.
// It records every node that was consulted and what each concluded.
type FlowOutput struct {
	FlowID    string
	TraceID   string
	Steps     []Step
	StartedAt time.Time
	EndedAt   time.Time
	Done      bool
}

// Step represents a single node consultation within a flow execution.
type Step struct {
	NodeID  string
	Verdict Verdict
}

func newFlowOutput(flowID, traceID string) *FlowOutput {
	return &FlowOutput{
		FlowID:    flowID,
		TraceID:   traceID,
		StartedAt: time.Now(),
	}
}

func (o *FlowOutput) record(nodeID string, verdict Verdict) {
	o.Steps = append(o.Steps, Step{NodeID: nodeID, Verdict: verdict})
}

func (o *FlowOutput) complete() {
	o.EndedAt = time.Now()
	o.Done = true
}

// LastOutput returns the Output value from the final node consulted.
// Returns nil if no steps were recorded.
func (o *FlowOutput) LastOutput() any {
	if len(o.Steps) == 0 {
		return nil
	}
	return o.Steps[len(o.Steps)-1].Verdict.Output
}

// Duration returns how long the flow execution took.
func (o *FlowOutput) Duration() time.Duration {
	return o.EndedAt.Sub(o.StartedAt)
}
