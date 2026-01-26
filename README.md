# 🎟️ TicketEngine: High-Frequency Distributed Reservation System

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Redis](https://img.shields.io/badge/Redis-Lua_Scripting-DC382D?style=flat&logo=redis)](https://redis.io/)
[![Postgres](https://img.shields.io/badge/Postgres-ACID_Compliant-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Orchestration-326CE5?style=flat&logo=kubernetes)](https://kubernetes.io/)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=flat&logo=docker)](https://www.docker.com/)

**TicketEngine** is a high-performance, concurrent event ticketing API designed to handle massive traffic spikes typical of "flash sales" (e.g., Taylor Swift concert bookings).

Unlike standard booking systems that fail under load, TicketEngine utilizes a **multi-layer locking architecture** combining **Redis Lua Scripting** (for optimistic memory-speed locking) and **PostgreSQL Row-Level Locking** (for pessimistic data integrity).

---

## ⚡ Key Highlights & Performance

* **10,000+ Requests Per Second (RPS):** Benchmarked using k6 on standard hardware.
* **Sub-15ms P95 Latency:** Achieved by offloading 99% of rejection traffic to the Redis Gatekeeper.
* **Zero Double-Bookings:** Proven concurrency safety via Atomic Distributed Locking.
* **Self-Healing:** 60-second TTL on locks prevents "zombie seats" during server crashes.

---

## 🏗️ System Architecture

The system follows a **Horizontal Slice Architecture** optimized for write-heavy workloads.

### The "Gatekeeper" Pattern
To protect the database from connection storms, traffic flows through a strict funnel:

1.  **Level 1: The Bouncer (Redis + Lua)**
    * Incoming requests trigger an atomic Lua script in Redis.
    * **Logic:** "If Key `seat:50` exists, return 0. Else, set Key and return 1."
    * **Latency:** ~1ms.
    * **Outcome:** 99.9% of concurrent requests are rejected here without ever touching the primary database.

2.  **Level 2: The Vault (PostgreSQL)**
    * Only the single "winner" from Level 1 proceeds to the database.
    * A `FOR UPDATE` row lock ensures serialized access for final persistence.
    * A `UNIQUE` constraint on the `bookings` table serves as the ultimate fail-safe.

3.  **Level 3: The Broadcaster (WebSockets)**
    * Upon successful booking, a Go channel pushes the update to the `Hub`.
    * The `Hub` broadcasts the "Sold Out" status to all connected frontend clients in real-time.

---

## 🛠️ Tech Stack

### Core Backend
* **Language:** Golang (Chi Router)
* **Database:** PostgreSQL (pgx driver with connection pooling)
* **Caching/Locking:** Redis (Go-Redis v9)
* **Real-Time:** Native WebSockets (Gorilla)

### Infrastructure & DevOps
* **Containerization:** Docker (Multi-stage builds)
* **Orchestration:** Kubernetes (StatefulSets for DB/Redis, Deployments for Go App)
* **Load Balancing:** Nginx (Round-robin distribution)
* **Observability:** Prometheus (Metrics scraping) & Grafana (Dashboards)
* **Load Testing:** k6 (scripted heavy-load simulation)

---

## 🚀 Engineering Deep Dive: The Locking Mechanism

The core challenge of this project was handling **Race Conditions**.
*Scenario: 500 users try to book Seat #A1 at the exact same millisecond.*

### ❌ The Naive Approach (Why it fails)
```go
// BAD CODE: A race condition nightmare
if db.GetSeatStatus(id) == "Available" {
    // 500 users enter this block simultaneously
    db.BookSeat(id) 
}
```
# 🎟️ TicketEngine — High-Concurrency Booking System

## ✅ The TicketEngine Approach (Redis Atomic Lock)

We utilize a custom **Lua script** to enforce atomicity. Redis guarantees that Lua scripts are executed **sequentially**.

~~~lua
-- The "Check-And-Set" Operation
if redis.call("EXISTS", KEYS[1]) == 1 then
    return 0 -- Failed: Seat is already locked
end
redis.call("SET", KEYS[1], ARGV[1], "EX", 60) -- Success: Lock acquired for 60s
return 1
~~~

This acts as a **distributed mutex**.  
Even if **5+ Go replicas** are running behind a load balancer, they all respect a **single source of truth** in Redis.

---

## 📊 Observability & Metrics

The system exposes a `/metrics` endpoint scraped by **Prometheus**.

### Key Metrics
- `http_request_duration_seconds` (P95, P99)
- `active_connections` (WebSocket count)
- `booking_success_rate` vs `booking_conflict_rate`


---

## 🏃 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+

---

### 🚀 Quick Start (Docker Compose)

~~~bash
git clone https://github.com/Abhinnavverma/High-Frequency-Distributed-Ticketing-Engine.git
docker-compose up --build -d
~~~

API available at: `http://localhost:8080`

---

### 🔥 Run Load Tests

~~~bash
k6 run scripts/load_test.js
~~~

---

## 🛣️ Future Roadmap

- [ ] Implement Redlock Algorithm
- [ ] Add gRPC endpoints
- [ ] Integrate Kafka

---
