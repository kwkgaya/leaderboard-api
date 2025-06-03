
# leaderboard-api

A leaderboard/competition web API in Go.

---

## Introduction

The problem statement is described in the `code-assessment-1.pdf` file.

---

## Prerequisites

- Visual Studio Code with the Go extension
- Go modules (install dependencies using `setup.sh`)
- Docker

---

## How to Run

**Option 1: Using Docker Compose**

1. Run `docker compose up` from a terminal (Linux, macOS, or Windows with Docker Desktop).
2. Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view the API documentation.

**Option 2: From within VS Code**

You can also start and debug the application directly from Visual Studio Code.

---

## Prometheus metrics

- Prometheus metrics of the service are accessible from [http://localhost:8080/metrics](http://localhost:8080/metrics)
- The metrics follow the guidelines from [this guide](https://prometheus.io/docs/practices/naming/)

### Available metrics

- `leaderboard_competitions_created_total` - Total number of competitions created
- `leaderboard_competitions_started_total` - Total number of competitions started
- TODO: Add more metrics

## Design Decisions and Trade-offs

- Invalid or empty arguments return HTTP status `400 Bad Request`, even if not specified in the API documentation.
- In-memory state is used to hold players and competitions. Adding players is not a thread-safe operation, but this is not an issue because players are always loaded at system startup. Access to the competitions map is synchronized using a mutex.
- Mutexes are used to synchronize critical paths. For higher performance, a message-processing model using goroutines and channels could be implemented.
- The minimum number of participants to start a competition is assumed to be 2.
- If a match is not found for a player within 30 seconds, a ticker fires every second to attempt matching and start the competition. This ticker currently keeps firing until a match is found. In the future, the ticker should stop after a configurable timeout.
- Constants are configured in the `constants.go` file in the `leaderboard/internal/config` package. Some constants are variables to allow changes during testing. In the future, all constants should be read from configuration (environment variables, command line, or config file).
- Some packages do not define interfaces (to save development time). These packages should be refactored to use interfaces.

---

## TODO

- Implement country code-based matching
- Refactor to use interfaces and a mock framework instead of mock function assignments
- Move more domain logic from the leaderboard and matchmaking packages to the model
- Improve performance by using less eager locking mechanisms or channel-based message passing
- Implement health endpoint