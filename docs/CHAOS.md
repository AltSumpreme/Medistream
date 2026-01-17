# Medistream – Failure & Resilience Testing

This document describes the failure scenarios tested (or designed)
to evaluate Medistream’s resilience, degradation behavior, and fault isolation.

The goal is not zero failure, but controlled, predictable failure.


## Failure Philosophy

Medistream is designed to:
- Fail fast at system boundaries
- Degrade gracefully under overload
- Isolate failures across subsystems
- Preserve core functionality under partial outages



## 1. CPU Saturation (Signup Overload)

### Scenario
- High concurrent signup traffic (200+ VUs)
- Each request performs password hashing (bcrypt)

### Expected Behavior
- Throughput plateaus
- Latency increases predictably
- Small percentage of signup failures
- No cascading failures

### Observed Results
- Throughput capped at ~85 workflows/sec
- p95 latency increased to ~3.2s
- ~1.8% signup failures
- Appointment creation remained functional

### Conclusion
CPU-bound workloads degrade gracefully without destabilizing the system.


## 2. Database Contention (Write-Heavy Load)

### Scenario
- High-volume appointment creation
- Concurrent writes to indexed tables

### Expected Behavior
- Increased latency due to lock contention
- No data corruption
- Error rate remains bounded

### Observations
- Latency increased under load
- No transaction failures or deadlocks observed
- System recovered immediately after load subsided


## 3. Redis Queue Backpressure

### Scenario
- High appointment creation rate triggering async tasks
- Redis used for background job queue (Asynq)

### Expected Behavior
- Queue depth increases
- Worker throughput limits processing rate
- API requests remain responsive

### Observations
- API latency unaffected by queue growth
- Failed tasks retried with exponential backoff
- Dead-letter strategy captured exhausted retries

### Conclusion
Async job failures are isolated from request path.



## 4. Invalid Authentication & Malformed Tokens

### Scenario
- Requests with malformed or expired JWTs

### Expected Behavior
- Immediate rejection
- No downstream processing

### Observations
- 401 responses returned early
- No DB or queue interaction



## 5. Rate Limiting Enforcement

### Scenario
- Excessive request rate from a single client

### Expected Behavior
- Requests rejected at middleware
- Core services protected

### Observations
- Rate-limited requests returned 429
- No observable impact on healthy traffic

## 6. Partial Dependency Failure (Planned)

### Scenario
- Redis unavailable
- PostgreSQL degraded or unavailable

### Expected Behavior
- API endpoints fail fast
- Async features disabled
- Clear error propagation

### Status
- Failure scenarios designed but not yet executed



## Key Takeaways

- Medistream favors **bounded failure over availability at all costs**
- Critical paths are isolated from background processing
- Overload scenarios degrade latency before correctness
- Failure modes are observable and predictable



## Future Chaos Experiments

- Inject DB latency and packet loss
- Kill worker processes under load
- Simulate Redis failover
- Test read-only fallback modes

These experiments will further validate production readiness.
