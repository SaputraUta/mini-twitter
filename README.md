# mini-twitter

Mini Twitter backend untuk belajar:
- Microservices architecture (Go)
- Fan-out at write pattern
- Redis timeline cache
- Message queue (NATS)
- Observability (LGTM stack, Day 2)
- Performance/latency analysis (Day 2)

Reference: DDIA Chapter 1 Twitter timeline example + Redis quicklist internals.

## Architecture

```
                ┌──────────┐
                │ Frontend │
                └────┬─────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
   ┌────▼────┐  ┌────▼────┐  ┌───▼─────────┐
   │  posts  │  │timeline │  │fanout-worker│
   │ service │  │ service │  │   (NATS)    │
   └────┬────┘  └────┬────┘  └───┬─────────┘
        │            │           │
        │       ┌────▼────┐      │
        │       │  Redis  │◄─────┘
        │       │(timeline│
        │       │ cache)  │
        │       └─────────┘
        │
   ┌────▼─────┐
   │ Postgres │
   │ (master) │
   └──────────┘
```

## Quick start

```bash
cp .env.example .env
docker compose up -d
```

Verify:
```bash
docker compose ps
docker compose logs -f
```

## Stack

- **Go 1.26+** — services
- **Postgres 16** — master data (users, posts)
- **Redis 7** — timeline cache
- **NATS 2 (JetStream)** — async event bus
- **React** — frontend
- **LGTM** — observability (Day 2)
- **k6** — load testing (Day 2)
