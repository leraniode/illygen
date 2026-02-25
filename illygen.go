// Package illygen is the top-level entry point for the Illygen runtime.
//
// Illygen enables developers to build AI-like intelligence systems that can
// reason, make decisions, and learn — without being full AI models.
// It mimics AI concepts (neural networks, training, knowledge) using
// deterministic, inspectable, resource-light Go machinery.
//
// # Quick Start
//
//	// 1. Define a node
//	type GreeterNode struct{ core.BaseNode }
//	func (n *GreeterNode) Consult(ctx *core.Context) (core.Verdict, error) {
//	    return core.Verdict{Output: "hello"}, nil
//	}
//
//	// 2. Wire a flow
//	flow := core.NewFlow("greeting").
//	    Add(&GreeterNode{core.NewBaseNode("greeter")})
//
//	// 3. Run it
//	rt := runtime.New()
//	rt.Register(flow)
//	out, _ := rt.Run(context.Background(), "greeting", core.NewContext("greeting", "run-1"))
//	fmt.Println(out.LastOutput())
//
// # Packages
//
//   - core:      Node, Flow, Context, Verdict — the building blocks
//   - knowledge: KnowledgeUnit, KnowledgeStore — the intelligence feed
//   - learning:  Trainer (training mode), Explorer (exploring mode)
//   - runtime:   The execution engine — goroutines, concurrency, lifecycle
package illygen
