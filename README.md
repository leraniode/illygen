<p align="center">
  <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/images/illygenbrandingimage.svg" width="1024" height="300" alt="Illygen"/>
</p>

# Illygen

<p align="center">
   <a href="https://github.com/leraniode">
      <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/badges/partofleraniode.svg" alt="Part of Leraniode" width="190" />
   </a>
   <a href="https://github.com/leraniode/illygen">
      <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/badges/illygenintelligenceleraniode.svg" alt="Illygen Intelligence" width="230" />
    </a>
</p>
<p align="center">
    <a href="https://github.com/leraniode/illygen/issues">
      <img src="https://img.shields.io/github/issues/leraniode/illygen" alt="GitHub Issues" />
    </a>
    <a href="https://github.com/leraniode/illygen/stargazers">
      <img src="https://img.shields.io/github/stars/leraniode/illygen" alt="GitHub Stars" />
    </a>
    <a href="https://github.com/leraniode/illygen/network/members">
      <img src="https://img.shields.io/github/forks/leraniode/illygen" alt="GitHub Forks" />
    </a>
    <a href="https://github.com/leraniode/illygen/blob/main/LICENSE">
       <img src="https://img.shields.io/github/license/leraniode/illygen" alt="GitHub License" />
    </a>
    <a href="https://github.com/leraniode/illygen/actions/workflows/ci.yml">
       <img src="https://img.shields.io/github/actions/workflow/status/leraniode/illygen/ci.yml" alt="GitHub Actions CI" />
    </a>
    <a href="https://pkg.go.dev/github.com/leraniode/illygen">
       <img src="https://img.shields.io/badge/pkg.go.dev-illygen-blue" alt="pkg.go.dev" />
    </a>
    <a href="https://goreportcard.com/report/github.com/leraniode/illygen">
       <img src="https://goreportcard.com/badge/github.com/leraniode/illygen" alt="Go Report Card" />
    </a>
    <a href="https://github.com/leraniode/illygen/commits/main">
       <img src="https://img.shields.io/github/last-commit/leraniode/illygen" alt="Last Commit" />
    </a>
</p>

> A Go-based library and runtime for building intelligence systems.

Illygen enables developers to build AI-like systems that can **reason, make decisions, and learn** — without being full AI models. It mimics the concepts used in AI using deterministic, inspectable, resource-light Go machinery.

**It is not a replacement for AI.** It is a lightweight alternative for domains where AI is overkill — embedded systems, edge computing, domain-specific reasoning engines, and smart automation.

---

## Concepts

| Concept       | What it is                                                                       |
| ------------- | -------------------------------------------------------------------------------- |
| **Node**      | A single unit of reasoning. You consult it, it returns a verdict. Like a neuron. |
| **Flow**      | A net of connected nodes — the reasoning pipeline. Like a neural network.        |
| **Knowledge** | The feed of intelligence. The more, the smarter the system.                      |

---

## Installation

To add it to your project, simply use:

```sh
go get github.com/leraniode/illygen@latest
```

Or add it to your `go.mod` by importing `github.com/leraniode/illygen` in your code and running `go mod tidy`.


---

## Quick Start

```go
package main

import (
    "fmt"

    illygen "github.com/leraniode/illygen"
)

func main() {
    // Define nodes as plain Go functions
    profiler := illygen.NewNode("profiler", func(ctx illygen.Context) illygen.Result {
        if ctx.Get("is_programmer") == true {
            return illygen.Result{Next: "go_intro", Confidence: 0.9}
        }
        return illygen.Result{Next: "prog_intro", Confidence: 0.7}
    })

    goIntro := illygen.NewNode("go_intro", func(ctx illygen.Context) illygen.Result {
        return illygen.Result{
            Value:      "Welcome! Here's Go: https://go.dev",
            Confidence: 1.0,
        }
    })

    progIntro := illygen.NewNode("prog_intro", func(ctx illygen.Context) illygen.Result {
        return illygen.Result{
            Value:      "Welcome! Programming is the art of telling computers what to do.",
            Confidence: 1.0,
        }
    })

    // Wire the flow
    flow := illygen.NewFlow().
        Add(profiler).
        Add(goIntro).
        Add(progIntro).
        Link("profiler", "go_intro", 1.0).
        Link("profiler", "prog_intro", 1.0)

    // Run it
    engine := illygen.NewEngine()

    result, err := engine.Run(flow, illygen.Context{
        "is_programmer": true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(result.Value)
    // Output: Welcome! Here's Go: https://go.dev
}
```

---

## Contribution

> [!NOTE]
> Contributions are welcome! Please open an issue or submit a pull request.

If you have ideas, suggestions, or want to contribute code, please feel free to:

- Join the discussion in the [Leraniode Discussions](https://github.com/leraniode/illygen/discussions)
- Open an issue for bugs or feature requests
- Submit a pull request with your changes

---

## License

MIT

---

A [leraniode](https://github.com/leraniode) project.

<p align="left">
   <a href="https://github.com/leraniode">
       <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/footer/leraniodeproductbrandimage.png" width="400" alt="Leraniode"/>
   </a>
</p>
