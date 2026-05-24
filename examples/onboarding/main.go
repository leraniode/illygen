// Onboarding Flow — Illygen v0.1.1 Example
//
// Demonstrates conditional routing based on context values:
//
//	profiler → go_intro
//	         ↘ prog_intro
//
// The profiler node reads a bool from context (ctx.Bool) to decide
// which intro node to route to. This shows how flows branch
// based on what's known about the user.
//
// Concepts shown:
//   - ctx.Bool() typed accessor
//   - Result.Next for explicit routing
//   - Multi-branch flows
//   - Reusing an engine across multiple runs
package main

import (
	"fmt"

	illygen "github.com/leraniode/illygen"
)

func main() {
	// ── Nodes ──────────────────────────────────────────────────────

	// profilerNode: reads user context and routes to the right intro.
	profilerNode := illygen.NewNode("profiler", func(ctx illygen.Context) illygen.Result {
		if ctx.Bool("is_programmer") {
			return illygen.Result{Next: "go_intro", Confidence: 0.95}
		}
		return illygen.Result{Next: "prog_intro", Confidence: 0.80}
	})

	// goIntroNode: for developers — point them to Go resources.
	goIntroNode := illygen.NewNode("go_intro", func(ctx illygen.Context) illygen.Result {
		name := ctx.String("name")
		greeting := "Welcome"
		if name != "" {
			greeting = "Welcome, " + name
		}
		return illygen.Result{
			Value: fmt.Sprintf(
				"%s! Since you're already a developer, here's Go: https://go.dev\n"+
					"Illygen itself is built in Go — you'll feel right at home.",
				greeting,
			),
			Confidence: 1.0,
		}
	})

	// progIntroNode: for newcomers — give them the big picture first.
	progIntroNode := illygen.NewNode("prog_intro", func(ctx illygen.Context) illygen.Result {
		name := ctx.String("name")
		greeting := "Welcome"
		if name != "" {
			greeting = "Welcome, " + name
		}
		return illygen.Result{
			Value: fmt.Sprintf(
				"%s! Programming is the art of telling computers what to do.\n"+
					"It's one of the most powerful skills you can learn — and you're starting now.",
				greeting,
			),
			Confidence: 1.0,
		}
	})

	// ── Flow ───────────────────────────────────────────────────────
	flow := illygen.NewFlow().
		Add(profilerNode).
		Add(goIntroNode).
		Add(progIntroNode).
		Link("profiler", "go_intro", 1.0).
		Link("profiler", "prog_intro", 1.0)

	engine := illygen.NewEngine()

	// ── Run for a programmer ───────────────────────────────────────
	fmt.Println("── Run 1: programmer ──────────────────────────────")
	r1, err := engine.Run(flow, illygen.Context{
		"is_programmer": true,
		"name":          "Ada",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%.0f%%] %v\n\n", r1.Confidence*100, r1.Value)

	// ── Run for a newcomer ─────────────────────────────────────────
	fmt.Println("── Run 2: newcomer ────────────────────────────────")
	r2, err := engine.Run(flow, illygen.Context{
		"is_programmer": false,
		"name":          "Jordan",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%.0f%%] %v\n\n", r2.Confidence*100, r2.Value)

	// ── Run with no context — safe nil handling ─────────────────────
	fmt.Println("── Run 3: no context ──────────────────────────────")
	r3, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%.0f%%] %v\n\n", r3.Confidence*100, r3.Value)
}

// must panics if err is non-nil. Used at startup to catch bad knowledge setup immediately.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
