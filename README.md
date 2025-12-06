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

```Go
shardIndex = hash(username)[0] % numShards
```

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

```sql
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
```

### Description

- **Client**: Sends HTTP requests to create or fetch users.  
- **API Handler**: Translates HTTP/JSON into domain calls.  
- **Service**: Applies business rules and coordinates actions.  
- **Router**: Selects the appropriate shard using a deterministic hash.  
- **Shard**: Executes SQL queries against its assigned PostgreSQL instance.  
- **Database**: Stores user records for the shard.

---

## C4 Model — Level 1: System Context

### Purpose
The system provides user creation and retrieval through a sharded PostgreSQL backend.

### Primary Actor
- **Client (User or External Application)**  
  Sends HTTP requests to create or retrieve a user.

### System Under Design (SUD)
- **Sharded User Service**  
  A Go-based backend responsible for distributing user data across multiple database shards.

### External Systems
- **PostgreSQL Shards**  
  Independent database instances acting as storage nodes.  
  The service routes users deterministically to one of these shards.

### High-Level Context Diagram (text form)

```sql
             +------------------------+
             |        Client          |
             |  (Browser/Postman/cURL)|
             +-----------+------------+
                         |
            HTTP Requests (Create/Get User)
                         |
             +-----------v------------+
             |  Sharded User Service  |
             |  (System Under Design) |
             +-----------+------------+
                         |
    +-------------+----------------------------+
    |       |            |              |      |
    |       |            |              |      |
    +-------v----+ +-----v------+ +-----v------+
    | PostgreSQL | | PostgreSQL | | PostgreSQL |
    |   Shard 0  | |   Shard 1  | |   Shard 2  |
    +------------+ +------------+ +------------+
       (Data is deterministically distributed)
```

```mermaid
graph TD
    A[Client (Browser/Postman/cURL)] --> B[API Layer]
    B --> C[Shared User Service]
    C --> D[PostgresSQL shard 0..2]
```

### Summary
- The client interacts only with the Sharded User Service.
- The service abstracts all storage logic and shard selection.
- The underlying PostgreSQL instances hold the actual data.

---

## C4 Model — Level 2: Container Diagram

This level details the internal structure of the Sharded User Service.

### Containers within the System

- **API (HTTP Server / Main Application)**
  - Exposes HTTP endpoints.
  - Initializes all dependencies (handlers, services, router, shards).

- **Handler Layer**
  - Parses HTTP requests (JSON, query params).
  - Validates basic inputs.
  - Delegates execution to the Service layer.

- **Service Layer (Business Logic)**
  - Contains domain logic (e.g., create user, get user).
  - Validates domain rules.
  - Selects the appropriate shard via the Router.
  - An orchestrator between Handler and Storage layers.

- **Router (Shard Selector)**
  - Applies hashing logic (`SHA-256 % number_of_shards`).
  - Routes each operation to its correct shard.
  - Encapsulates the sharding strategy.

- **Shard (Storage Implementation)**
  - Executes SQL statements.
  - Connected to one physical PostgreSQL instance.
  - Represents one node in the sharded architecture.

- **PostgreSQL Instances**
  - Independent databases.
  - Each shard stores only the portion of data assigned to it.

---

### Container-level Diagram (text form)

```sql
+-------------------------------------------------------------+
|                 Sharded User Service (Go)                   |
|                                                             |
|     +------------------+       +---------------+            |
|     |     API Layer    | ----> | Handler Layer |            |
|     +------------------+       +---------------+            |
|                                        |                    |
|                                        v                    |
|                                +---------------+            |
|                                | Service Layer |            |
|                                +---------------+            |
|                                        |                    |
|                                        v                    |
|                               +-----------------+           |
|                               |      Router     |           |
|                               |  (Shard Picker) |           |
|                               +-----------------+           |
|                                      / | \                  |
|                                     /  |  \                 |
|                                    v   v   v                |
|                       +---------+ +---------+ +---------+   |
|                       | Shard 0 | | Shard 1 | | Shard 2 |   |
|                       +---------+ +---------+ +---------+   |
|                            |           |           |        |
|                            v           v           v        |
|                        PostgreSQL  PostgreSQL  PostgreSQL   |
|                           DB0         DB1         DB2       |
+-------------------------------------------------------------+
```

### Summary

- The API exposes endpoints but contains no business logic.  
- The Handler is only responsible for HTTP concerns.  
- The Service contains all domain rules and orchestrates operations.  
- The Router encapsulates the sharding strategy.  
- Each Shard is a thin wrapper around SQL operations.  
- PostgreSQL instances are isolated, independent storage nodes.

---

## ADR #1 — Using Hash Mod N for Sharding

### Status
Accepted — educational purpose.

---

### Context
The project requires distributing user data across multiple PostgreSQL instances.  
A deterministic strategy is needed so that the same username always maps to the same shard.

The goal is **simplicity and clarity**, not production-grade scalability.

---

### Decision
Use SHA-256 hashing on the username and take the modulo of the number of shards:

```Go
shardIndex = hash(username)[0] % numShards
```

The `Router` component encapsulates this logic and exposes a simple interface:

- `SaveUser(ctx, username, pwd)`
- `GetUser(ctx, username)`

---

### Rationale
This choice provides:

- Extremely simple implementation  
- Deterministic mapping of users to shards  
- Easy to understand and teach  
- No external dependencies  
- Good for prototyping distributed storage concepts  

---

### Consequences

#### Positive
- Fast and deterministic routing  
- Minimal code complexity  
- Easy to test and reason about  
- Good starting point for learning sharding concepts  

#### Negative
- **Not scalable** for production systems  
- Adding/removing shards requires redistributing all data  
- No automatic load balancing  
- No shard hot-spot protection  
- No replica support  

---

### Alternatives Considered
1. **Consistent Hashing**
   - Much better scalability
   - Avoids full data redistribution
   - Used by large-scale systems (Kafka, DynamoDB, Cassandra, Redis Cluster)
   - Rejected because it adds complexity beyond the scope of this project.

2. **Lookup Tables (User → Shard)**
   - Flexible and dynamic
   - Requires a central metadata service
   - Harder to maintain in small demo projects

---

### Decision Summary
Hash Mod N is chosen **deliberately** as a simple, clear, educational approach to demonstrate:

- routing  
- multi-shard design  
- system layering  
- database distribution  

