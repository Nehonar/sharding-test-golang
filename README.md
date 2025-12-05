# go-sharding-basic

Example system in Go that implements basic sharding using `hash(username) % N` to distribute users across multiple PostgreSQL shards.

This project showcases:

- Clear separation of layers (API, Handler, Service, Router, Storage, Shards)
- Clean and maintainable design
- Deterministic routing for reads and writes
- Real sharding across multiple PostgreSQL databases
- Architecture designed to scale horizontally
- Proper usage of Go modules, interfaces, and dependency boundaries

---

## Architecture Overview

Main flow (CreateUser / GetUser):

Client → Handler → Service → Router → Shard → Database

Each layer has a single responsibility:

- **Handler**: HTTP layer (JSON ↔ domain)
- **Service**: business logic
- **Router**: selects the correct shard
- **Shard**: executes SQL
- **Database**: persists user data
- **Models**: pure data structures

---

## Architectural Decision: Hash Mod N

Sharding is implemented by applying SHA-256 to the username and taking the modulo of the number of shards:

shardIndex = hash(username)[0] % numShards

This guarantees that the same username always maps to the same shard.

---

## Limitations

This sharding method is **NOT scalable for production environments**, because:

- Changing the number of shards causes all users to map to different shards.
- No controlled redistribution of data.
- No replicas or fault tolerance.
- No dynamic load balancing.

This project is an educational example, **not** a production-ready implementation.

---

## Running the System

1. Start the shards:

```bash
docker compose up -d
```

2. Run the server:

```bash
go run cmd/server/main.go
```

3. Create a user:

```bash
curl -X POST http://localhost:8080/create-user

-H "Content-Type: application/json"
-d '{"user":"user-test","password":"1234"}'
```

4. Fetch a user:

```bash
curl "http://localhost:8080/get-user?user=user-test"
```

---

## Educational Purpose

This project is part of a portfolio designed to demostrate:

- Software architecture skills

- Understanding of distributed systems

- Knowledge of real-world architectural patterns

- Ability to document and communicate system design effectively


---


## 6-Box Architecture Diagram

       Client
         ↓
     API Handler
         ↓
 Service (Business Logic)
         ↓
  Router (Shard Selector)
         ↓
Shard (PostgreSQL Instance)
         ↓
      Database

### Description

- **Client**: Sends HTTP requests to create or fetch users.  
- **API Handler**: Translates HTTP/JSON into domain calls.  
- **Service**: Applies business rules and coordinates actions.  
- **Router**: Selects the appropriate shard using a deterministic hash.  
- **Shard**: Executes SQL queries against its assigned PostgreSQL instance.  
- **Database**: Stores user records for the shard.


---


