# Illygen

> A Go-based library and runtime for building intelligence systems.

Illygen enables developers to build AI-like systems that can **reason, make decisions, and learn** — without being full AI models. It mimics the concepts used in AI using deterministic, inspectable, resource-light Go machinery.

**It is not a replacement for AI.** It is a lightweight alternative for domains where AI is overkill — embedded systems, edge computing, domain-specific reasoning engines, and smart automation.

---

## Concepts

| Concept | What it is |
|---|---|
| **Node** | A single unit of reasoning. You consult it, it returns a verdict. Like a neuron. |
| **Flow** | A net of connected nodes — the reasoning pipeline. Like a neural network. |
| **Knowledge** | The feed of intelligence. The more, the smarter the system. |
| **Learning** | Training (high impact, offline) and Exploring (low impact, online). |

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/leraniode/illygen/core"
    "github.com/leraniode/illygen/runtime"
)

// Define a node
type UserProfilerNode struct{ core.BaseNode }

func (n *UserProfilerNode) Consult(ctx *core.Context) (core.Verdict, error) {
    if ctx.Get("is_programmer") == true {
        return core.Verdict{Route: "go_intro", Output: "programmer detected", Weight: 0.9}, nil
    }
    return core.Verdict{Route: "prog_intro", Output: "new to programming", Weight: 0.7}, nil
}

type GoIntroNode struct{ core.BaseNode }

func (n *GoIntroNode) Consult(ctx *core.Context) (core.Verdict, error) {
    return core.Verdict{Output: "Welcome, here's Go: https://go.dev"}, nil
}

type ProgrammingIntroNode struct{ core.BaseNode }

func (n *ProgrammingIntroNode) Consult(ctx *core.Context) (core.Verdict, error) {
    return core.Verdict{Output: "Welcome! Programming is the art of telling computers what to do."}, nil
}

func main() {
    // Wire the flow
    flow := core.NewFlow("onboarding").
        Add(&UserProfilerNode{core.NewBaseNode("profiler")}).
        Add(&GoIntroNode{core.NewBaseNode("go_intro")}).
        Add(&ProgrammingIntroNode{core.NewBaseNode("prog_intro")}).
        Connect("profiler", "go_intro").
        Connect("profiler", "prog_intro")

    // Boot the runtime
    rt := runtime.New()
    rt.Register(flow)

    // Run it
    ctx := core.NewContext("onboarding", "run-001")
    ctx.Set("is_programmer", true)

    out, err := rt.Run(context.Background(), "onboarding", ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println(out.LastOutput())
    // Output: Welcome, here's Go: https://go.dev
}
```

---

## Packages

```
illygen/
├── core/        Node, Flow, Context, Verdict — the building blocks
├── knowledge/   KnowledgeUnit, KnowledgeStore — the intelligence feed
├── learning/    Trainer + Explorer — the two learning modes
├── runtime/     The execution engine
└── internal/    Graph primitives (not for direct use)
```

---

## Roadmap

| Version | Scope |
|---|---|
| **v0.1** | Node, Flow, Context, Runtime, KnowledgeStore ← *you are here* |
| v0.2 | Training logic, weighted connection reinforcement |
| v0.3 | Exploring logic, incremental runtime adaptation |
| v0.4+ | Community |

---

## Organisation

Illygen is developed under [leraniode](https://github.com/leraniode).

Contributions, ideas, and feedback are welcome. Read [DESIGN.md](./DESIGN.md) to understand the vision before contributing.

---

## License

MIT
