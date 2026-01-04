
---

# 🚀 THE PROJECT (Core)

## **“Microscope” — A Self-Observing, Self-Breaking Go System**

> A Go-based micro-platform that **spawns, throttles, breaks, heals, and lies about itself** —
> while Prometheus watches everything in real time.

This is **not** a monitoring project.
This is a **system designed to be monitored**.

---

## 🔥 High-Level Idea

You build a **Go service orchestrator** that:

• Spawns **hundreds of tiny Go workers** (containers or goroutines)
• Intentionally introduces **chaos** (latency, memory leaks, CPU spikes)
• Uses **Prometheus metrics to detect its own failure patterns**
• Automatically **adapts behavior** based on metrics
• Runs in **absurdly constrained resources**

Think:

> “What if a system could watch Prometheus… and react?”

---

## 🧠 What Makes This “Cool” (By Your Definition)

| Your benchmark          | This project equivalent                         |
| ----------------------- | ----------------------------------------------- |
| Minecraft in 300KB RAM  | Go workers capped at **5–10MB RSS**             |
| 1k Docker in 1GB        | **1000 goroutines pretending to be services**   |
| Impossible resource use | **Self-throttling via metrics feedback loop**   |
| Meme but real           | **Prometheus drives decisions, not dashboards** |

---

## 🏗️ Core Architecture

```
microscope/
├── cmd/
│   ├── controller/      # Master brain
│   ├── worker/          # Tiny services
│
├── internal/
│   ├── metrics/         # Prometheus instrumentation
│   ├── chaos/           # Fault injectors
│   ├── scheduler/       # Adaptive logic
│   ├── limiter/         # CPU / mem throttling
│
├── deploy/
│   ├── docker/
│   ├── k8s/             # Optional but insane
│
└── docs/
```

---

## 🔬 What You Build (Step-by-Step)

### **Phase 1 – Ultra-Minimal Go Services**

Each worker:

* Exposes `/metrics`
* Handles fake “requests”
* Uses **custom counters, histograms, summaries**

You’ll learn:

* `prometheus/client_golang`
* Labels (correctly, not cardinality hell)
* Process metrics vs custom metrics

---

### **Phase 2 – Metrics That Actually Matter**

Workers emit:

* Request latency
* Memory growth rate
* Goroutine count
* Error bursts
* Synthetic SLOs

Controller scrapes Prometheus and learns:

* “This worker is lying”
* “This worker is slow but stable”
* “This one is about to OOM”

You learn:

* Histogram buckets
* Quantiles vs percentiles
* RED / USE methods
* Recording rules

---

### **Phase 3 – Chaos Engineering (But Tiny)**

Workers can:

* Leak memory slowly
* Sleep randomly
* Spin CPU
* Panic occasionally

Controller:

* Detects patterns via PromQL
* Kills / restarts / rate-limits workers

You learn:

* PromQL deeply (rate, increase, irate, predict_linear)
* Alert rules
* Why alerts are hard

---

### **Phase 4 – Feedback Loop (🔥 Production Brain)**

This is the **rare part**.

Controller:

* Queries Prometheus HTTP API
* Uses metrics to make **runtime decisions**
* Adjusts:

  * Worker count
  * Request rate
  * Resource usage
  * Chaos level

This teaches:

* Prometheus as **control plane**
* Go HTTP clients + JSON decoding
* Why most systems don’t do this (but should)

---

### **Phase 5 – Ridiculous Constraints Mode**

Run the entire system with:

* 256MB RAM
* CPU quota
* Network delay

Goals:

* Keep 99% latency under X
* Survive chaos
* Never OOM

You learn:

* Memory profiling
* Go GC tuning
* `GOMEMLIMIT`
* Why “efficient Go” actually matters

---

## 🧪 Metrics You’ll Implement (Real Ones)

Not toy metrics like `requests_total`.

You’ll implement:

* **Error budget burn rate**
* **Latency SLO violation rate**
* **Worker health score (derived metric)**
* **Instability index** (variance over time)
* **Chaos tolerance score**

These are **interview / GSoC / LFX gold**.

---

## 🧨 Optional Insane Add-Ons (Pick One)

### 🔥 Add-On 1: **“Metrics Liar Detector”**

Workers randomly falsify metrics.

Controller:

* Detects inconsistencies
* Flags suspicious workers

You learn:

* Trust boundaries
* Cross-metric validation
* Observability limits

---

### 🔥 Add-On 2: **“Prometheus Under Attack”**

Simulate:

* Metric explosions
* Label cardinality attacks
* Scrape overload

You learn:

* Why Prometheus dies in production
* How to defend it
* How orgs screw this up

---

## 🧠 What You’ll Learn (Mapped Explicitly)

### Go:

* Concurrency (goroutines, channels)
* HTTP servers & clients
* Memory profiling
* Graceful shutdown
* Context propagation
* Production logging

### Prometheus:

* Client instrumentation
* PromQL (deep)
* Alerting philosophy
* Recording rules
* Cardinality management
* Performance limits
* Prometheus HTTP API

### Production Thinking:

* SLOs vs SLIs
* Feedback systems
* Failure patterns
* Why dashboards lie
* Why alerts wake people up at 3AM

---

