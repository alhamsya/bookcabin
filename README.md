# <img align="right" src="https://avatars.githubusercontent.com/u/56905970?s=60&v=4" alt="bookcabin" title="bookcabin" /> bookcabin

## 👀 Architecture
- architecture repository: Hexagonal Architecture
- web framework: fiber (https://gofiber.io/)
- retry mechanism: library (https://github.com/avast/retry-go)
- mock APIs: mockoon (https://mockoon.com)
- the service uses two-level caching (L1 + L2):
  - Redis (L2 Cache)
    - used as a distributed cache
    - ensures consistency across multiple pods
  - In-Memory Cache (L1 Cache)
    - Used for ultra-fast access
    - Acts as a fallback when Redis experiences like network disruption or temporary unavailability

## 🎯 Structure
- `cmd`: directory for main entry points or commands of the application
- `internal`: directory for containing application code that should not exposed to external packages
- `core`: directory that contains the central business logic of the application
  - `domain`: directory that contains domain models/entities representing core business concepts
  - `port`: directory that contains defined interfaces or contracts that adapters must follow
  - `service`: directory that contains the business logic or services of the application
- `pkg`: shared managing to support service and utilities
```
cmd
└── rest
internal
├── adapter
│   ├── airline
│   │   ├── airasia
│   │   ├── batik
│   │   ├── garuda
│   │   └── lion
│   ├── handler
│   │   └── rest
│   └── redis
└── core
    ├── domain
    │   ├── airline
    │   ├── config
    │   ├── constant
    │   ├── flight
    │   ├── request
    │   └── response
    ├── port
    │   ├── repository
    │   └── service
    └── service
        └── flight
pkg
├── manager
│   ├── config
│   ├── graceful
│   ├── protocol
│   ├── response
│   └── xhttp
└── util
```

## ⚡️ Prerequisites
1. **Go**: version 1.25.5 or higher is required
2. **Make**: running Makefile commands
3. **Docker**: running redis and mockoon service

## ⚙️ Installation and Setup
1. clone this repository:
    ```bash
    git clone https://github.com/alhamsya/bookcabin.git
    cd bookcabin
    ```
2. download dependencies
    ```bash
    go mod download
    ```
3. first setup in local environment
    ```bash
    make start-local
    ```

## ⚙️ Running Tests
1. start service in local environment
    ```bash
    go run ./cmd/main.go run rest
    ```