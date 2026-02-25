// Package runtime provides the execution engine for Illygen flows.
// It orchestrates goroutines, manages flow lifecycles, and handles memory.
package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/leraniode/illygen/core"
)

// Runtime is the heart of Illygen.
// It registers flows, runs them concurrently, and manages their lifecycle.
type Runtime struct {
	mu    sync.RWMutex
	flows map[string]*core.Flow
}

// New creates a new Illygen runtime.
func New() *Runtime {
	return &Runtime{
		flows: make(map[string]*core.Flow),
	}
}

// Register adds a flow to the runtime.
func (r *Runtime) Register(flow *core.Flow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flows[flow.ID()]; exists {
		return fmt.Errorf("illygen/runtime: flow %q already registered", flow.ID())
	}
	r.flows[flow.ID()] = flow
	return nil
}

// Run executes a flow by ID in a goroutine and blocks until it completes.
// Use RunAsync for non-blocking execution.
func (r *Runtime) Run(ctx context.Context, flowID string, flowCtx *core.Context) (*core.FlowOutput, error) {
	flow, err := r.get(flowID)
	if err != nil {
		return nil, err
	}

	type result struct {
		output *core.FlowOutput
		err    error
	}

	ch := make(chan result, 1)

	go func() {
		out, err := flow.Execute(flowCtx)
		ch <- result{out, err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("illygen/runtime: flow %q cancelled: %w", flowID, ctx.Err())
	case res := <-ch:
		return res.output, res.err
	}
}

// RunAsync executes a flow in a goroutine and returns immediately.
// The caller receives the output and any error through the returned channels.
func (r *Runtime) RunAsync(flowID string, flowCtx *core.Context) (<-chan *core.FlowOutput, <-chan error) {
	outCh := make(chan *core.FlowOutput, 1)
	errCh := make(chan error, 1)

	go func() {
		flow, err := r.get(flowID)
		if err != nil {
			errCh <- err
			return
		}

		output, err := flow.Execute(flowCtx)
		if err != nil {
			errCh <- err
			return
		}
		outCh <- output
	}()

	return outCh, errCh
}

// RunMany executes multiple flows concurrently and waits for all to finish.
// Returns a map of flowID → FlowOutput and a combined slice of any errors.
func (r *Runtime) RunMany(ctx context.Context, runs []FlowRun) (map[string]*core.FlowOutput, []error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make(map[string]*core.FlowOutput)
		errors  []error
	)

	for _, run := range runs {
		wg.Add(1)
		go func(fr FlowRun) {
			defer wg.Done()
			out, err := r.Run(ctx, fr.FlowID, fr.Ctx)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = append(errors, err)
			} else {
				results[fr.FlowID] = out
			}
		}(run)
	}

	wg.Wait()
	return results, errors
}

// FlowRun bundles a flow ID with its execution context for RunMany.
type FlowRun struct {
	FlowID string
	Ctx    *core.Context
}

func (r *Runtime) get(flowID string) (*core.Flow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flow, ok := r.flows[flowID]
	if !ok {
		return nil, fmt.Errorf("illygen/runtime: flow %q not found", flowID)
	}
	return flow, nil
}
