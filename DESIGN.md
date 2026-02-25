# Illygen — Design Document

> A Go-based library and runtime for building intelligence systems.

## What Illygen Is

Illygen enables developers to build AI-like systems that can reason, make decisions, and learn —
without being full AI models. It mimics the *concepts* used in AI (neural networks, training,
knowledge) using deterministic, inspectable, resource-light machinery.

It is not a replacement for AI. It is a lightweight alternative for domains where AI is overkill.

---

## Core Concepts

### Node
A node is the atomic unit of reasoning in Illygen — like a neuron in a neural network.
You consult a node with context, and it returns a verdict. Each node holds knowledge relevant
to its domain and uses that knowledge to make decisions.

### Flow
A flow is a net of connected nodes — the reasoning pipeline. Like a neural network, it is not
static. It reshapes itself as learning happens. Connections have weights that shift over time,
routes strengthen or weaken, and the graph evolves.

### Knowledge
Knowledge is the feed of intelligence. The more knowledge, the more intelligent the system.
Knowledge is stored as `KnowledgeUnit` — lightweight structures similar to tensors, holding
structured facts, a domain, and a trust weight.

### Logic
Illygen has four built-in logics that form the engine:

| Logic | Responsibility |
|---|---|
| **Flow Logic** | How nodes connect, how signals pass, branching and merging |
| **Context Logic** | State scoping and lifetime across a flow execution |
| **Learning Logic** | Weight adjustment and graph mutation based on outcomes |
| **Runtime Logic** | Goroutine orchestration, memory allocation, lifecycle |

---

## Learning

Learning comes in two forms:

### Training (Low-level, High Impact)
- Happens offline — before or during testing
- Developer feeds structured data directly into the KnowledgeStore
- Can significantly reshape the flow graph: reweight connections, restructure paths
- Think: *building the brain before deployment*

### Exploring (High-level, Low Impact)
- Happens online — during real system usage
- System observes its own flow outcomes and incrementally adjusts
- Small weight nudges, soft route preferences, pattern accumulation
- Bounded: Exploring can never override what Training established
- Think: *the brain refining itself through experience*

Training writes the skeleton. Exploring adds muscle memory.

---

## Memory Model

Three explicit memory scopes exist in the runtime:

| Scope | Lifetime | Purpose |
|---|---|---|
| `NodeMemory` | Single node execution | Goroutine-local state |
| `FlowMemory` | Single flow run | Shared context across nodes |
| `KnowledgeStore` | Persistent | Feeds all nodes, namespaced by domain |

---

## Tradeoffs

Illygen wins where real AI loses:

| | Real AI | Illygen |
|---|---|---|
| Resources | Massive | Lightweight |
| Inspectability | Low | High |
| Flexibility | Unlimited | Domain-scoped |
| Predictability | Low | High |
| Training data | Massive datasets | Curated knowledge feeds |
