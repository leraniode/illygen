// Intent Detection — Illygen v0.1.1 Official Example
//
// Demonstrates the complete Illygen v0.1.1 public API:
//
//	input → action
//
// The flow classifies intent from user text, queries the KnowledgeStore
// for a matching response, and returns it with a confidence score.
//
// Concepts shown:
//   - NewNode / NewFlow / NewEngine
//   - Context.String(), Context.Set()
//   - KnowledgeStore + domain-scoped queries
//   - illygen.Knowledge(ctx) inside a node
//   - Result.Confidence for full transparency
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	illygen "github.com/leraniode/illygen"
)

func main() {
	// ── Knowledge ──────────────────────────────────────────────────
	// Feed the store with facts across named domains.
	// Nodes query by domain during execution — they never touch the whole store.
	store := illygen.NewKnowledgeStore()

	must(store.Add("greet-1", "greetings", map[string]any{
		"response": "Hi! I'm Illygen — a lightweight intelligence engine. How can I help?",
	}))
	must(store.Add("greet-2", "greetings", map[string]any{
		"response": "Hello! What can I do for you today?",
	}))
	must(store.Add("bye-1", "farewells", map[string]any{
		"response": "Goodbye! Come back anytime — I'll keep learning.",
	}))
	must(store.Add("fact-illygen", "facts", map[string]any{
		"topic":    "illygen",
		"response": "Illygen is a lightweight intelligence engine built in Go. It uses flows, nodes, and knowledge to reason — no AI models needed.",
	}))
	must(store.Add("fact-node", "facts", map[string]any{
		"topic":    "node",
		"response": "A node is the atomic unit of reasoning in Illygen. Define its logic as a Go function, connect it in a flow, and the engine does the rest.",
	}))
	must(store.Add("fact-flow", "facts", map[string]any{
		"topic":    "flow",
		"response": "A flow is a connected net of nodes — Illygen's reasoning pipeline. Like a neural network, it reshapes itself as it learns.",
	}))
	must(store.Add("fact-knowledge", "facts", map[string]any{
		"topic":    "knowledge",
		"response": "Knowledge is the feed of intelligence in Illygen. Stored as KnowledgeUnits with structured facts and a trust weight.",
	}))
	must(store.Add("fact-signal", "facts", map[string]any{
		"topic":    "signal",
		"response": "A signal (coming in v0.2) carries sparse features and a full execution trace through the flow — richer than a plain Result.",
	}))

	// ── Nodes ──────────────────────────────────────────────────────

	// inputNode: reads text, classifies intent, routes to actionNode.
	inputNode := illygen.NewNode("input", func(ctx illygen.Context) illygen.Result {
		text := strings.ToLower(strings.TrimSpace(ctx.String("text")))

		switch {
		case isGreeting(text):
			ctx.Set("intent", "greeting")
			return illygen.Result{Next: "action", Confidence: 0.95}
		case isFarewell(text):
			ctx.Set("intent", "farewell")
			return illygen.Result{Next: "action", Confidence: 0.95}
		case isQuestion(text):
			ctx.Set("intent", "question")
			ctx.Set("query", text)
			return illygen.Result{Next: "action", Confidence: 0.80}
		default:
			ctx.Set("intent", "unknown")
			return illygen.Result{Next: "action", Confidence: 0.30}
		}
	})

	// actionNode: reads intent from context, queries knowledge, returns response.
	// Uses illygen.Knowledge(ctx) to access the store injected by the engine.
	actionNode := illygen.NewNode("action", func(ctx illygen.Context) illygen.Result {
		intent := ctx.String("intent")
		ks := illygen.Knowledge(ctx)

		switch intent {
		case "greeting":
			if units := ks.Domain("greetings"); len(units) > 0 {
				return illygen.Result{Value: units[0].Fact("response"), Confidence: 0.95}
			}

		case "farewell":
			if units := ks.Domain("farewells"); len(units) > 0 {
				return illygen.Result{Value: units[0].Fact("response"), Confidence: 0.95}
			}

		case "question":
			query := ctx.String("query")
			for _, unit := range ks.Domain("facts") {
				topic, _ := unit.Fact("topic").(string)
				if topic != "" && strings.Contains(query, topic) {
					return illygen.Result{Value: unit.Fact("response"), Confidence: 0.85}
				}
			}
			return illygen.Result{
				Value:      "I don't have knowledge about that yet — but I'm built to learn.",
				Confidence: 0.30,
			}
		}

		return illygen.Result{
			Value:      "I'm not sure how to respond. Try: hi · what is illygen · bye.",
			Confidence: 0.20,
		}
	})

	// ── Flow ───────────────────────────────────────────────────────
	flow := illygen.NewFlow().
		Add(inputNode).
		Add(actionNode).
		Link("input", "action", 1.0)

	// ── Engine ─────────────────────────────────────────────────────
	engine := illygen.NewEngine(store)

	// ── REPL ───────────────────────────────────────────────────────
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   Illygen v0.1.1 — Intent Detection Example      ║")
	fmt.Println("║                                                  ║")
	fmt.Println("║   Try: hi · what is illygen · what is a node     ║")
	fmt.Println("║        what is a flow · what is a signal · bye   ║")
	fmt.Println("║   Type 'exit' to quit                            ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Println("Illygen: Goodbye!")
			break
		}

		result, err := engine.Run(flow, illygen.Context{"text": text})
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Illygen [%.0f%%]: %v\n\n", result.Confidence*100, result.Value)
	}
}

// ── Helpers ────────────────────────────────────────────────────────

func isGreeting(s string) bool {
	for _, g := range []string{"hi", "hello", "hey", "yo", "howdy", "greetings"} {
		if s == g || strings.HasPrefix(s, g+" ") {
			return true
		}
	}
	return false
}

func isFarewell(s string) bool {
	for _, f := range []string{"bye", "goodbye", "see you", "later", "farewell", "ciao"} {
		if strings.Contains(s, f) {
			return true
		}
	}
	return false
}

func isQuestion(s string) bool {
	for _, q := range []string{"what", "who", "how", "why", "when", "where", "tell me", "explain"} {
		if strings.HasPrefix(s, q) {
			return true
		}
	}
	return false
}

// must panics if err is non-nil. Used at startup to catch bad knowledge setup immediately.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
