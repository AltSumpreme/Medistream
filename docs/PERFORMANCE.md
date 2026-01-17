# Medistream – Performance & Scalability Analysis

This document summarizes the performance characteristics, capacity limits,
and overload behavior of the Medistream backend.

---

## Workload Model

I benchmarked realistic user workflows using **k6** as well as **wrk**, focusing on:

1. **Signup + Authenticated Appointment Creation**
   - Password hashing (bcrypt)
   - User creation (PostgreSQL)
   - JWT issuance
   - Appointment creation (PostgreSQL)
   - Async enqueue (Redis + Asynq)

2. **Appointments-Only (Steady State)**
   - Pre-authenticated requests
   - Isolates business logic from auth overhead

---

## Test Environment

- Backend: Go (Gin)
- Database: PostgreSQL
- Cache / Queue: Redis + Asynq
- Metrics: Prometheus + Grafana
- Execution: Local Docker Compose
- Load Tool: k6

---

## Results Summary

### Signup → Appointment Workflow

| Concurrent Users | Throughput (workflows/sec) | p95 Latency | Error Rate |
|------------------|---------------------------|-------------|------------|
| 50 VUs           | ~84                       | ~0.8s       | 0%         |
| 100 VUs          | ~84                       | ~1.6s       | 0%         |
| 200 VUs          | ~85                       | ~3.2s       | ~1.8%      |

**Observations**
- Throughput saturates at ~85 workflows/sec
- Increasing concurrency increases latency, not throughput
- Errors appear only after saturation
- System degrades gracefully under overload

**Root Cause**
- CPU-bound password hashing during signup

---

### Appointments-Only (Steady State)

- Significantly higher throughput than signup workflows
- Lower p95 latency
- Bottleneck shifts from CPU-bound auth to DB writes and indexes

---

## Key Takeaways

- Signup is intentionally expensive and rate-limited
- Core business workflows scale independently of auth
- The system exhibits predictable saturation and graceful degradation
- No cascading failures or instability under load


