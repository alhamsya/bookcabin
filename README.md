# <img align="right" src="https://avatars.githubusercontent.com/u/56905970?s=60&v=4" alt="bookcabin" title="bookcabin" /> bookcabin

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alhamsya/bookcabin)
[![Sourcegraph](https://sourcegraph.com/github.com/alhamsya/bookcabin/-/badge.svg)](https://sourcegraph.com/github.com/alhamsya/bookcabin?badge)
[![Documentation](https://godoc.org/github.com/alhamsya/bookcabin?status.svg)](https://godoc.org/github.com/alhamsya/bookcabin)
[![codecov](https://codecov.io/gh/alhamsya/bookcabin/graph/badge.svg?token=PIN65DKRGQ)](https://codecov.io/gh/alhamsya/bookcabin)
[![Go Report Card](https://goreportcard.com/badge/github.com/alhamsya/bookcabin)](https://goreportcard.com/report/github.com/alhamsya/bookcabin)
[![License](https://img.shields.io/github/license/alhamsya/bookcabin?color=blue)](https://raw.githubusercontent.com/alhamsya/bookcabin/master/LICENSE)

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
cmd/
└── rest
internal/
├── adapter/
│   ├── airline/
│   │   ├── airasia
│   │   ├── batik
│   │   ├── garuda
│   │   └── lion
│   ├── handler/
│   │   └── rest
│   └── redis
├── core/
│   ├── domain/
│   │   ├── airline
│   │   ├── config
│   │   ├── constant
│   │   ├── flight
│   │   ├── request
│   │   └── response
│   ├── port/
│   │   ├── repository
│   │   └── service
│   └── service/
│       └── flight
└── mock/
    ├── repository
    └── service
pkg/
├── manager/
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

## ⚡️ Mock Documentation
by default host mock: `localhost:3000`
1. **Garuda Indonesia**: `GET {{HOST_MOCK}}/airline/gruda`
2. **Lion Air**: `GET {{HOST_MOCK}}/airline/lion`
3. **Batik Air**: `GET {{HOST_MOCK}}/airline/batik`
4. **AirAsia**: `GET {{HOST_MOCK}}/airline/airasia`

## ⚡️ API Documentation
### Search Flights
**Endpoint**: `POST /v1/flights/search`

**Request Body**:
```json
{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "arrivalDate": "",
    "passengers": 1,
    "cabinClass": "economy",
    "sort": {
        "key": "price",
        "order": "desc"
    },
    "filters": {
        "minPrice": 2000,
        "maxPrice": 2000000,
        "stops": [
            1
        ],
        "airlines": [],
        "maxDurationMinutes": 1000
    }
}
```

**Response** `OK (200)`:
```json
{
    "data": {
        "data": [
            {
                "ID": "ID7042_BATIK",
                "Provider": "BATIK",
                "Airline": {
                    "Name": "Batik Air",
                    "Code": "ID"
                },
                "Route": {
                    "Origin": "CGK",
                    "Destination": "DPS"
                },
                "Schedule": {
                    "DepartureTime": "2025-12-15T18:45:00+07:00",
                    "ArrivalTime": "2025-12-15T23:50:00+08:00",
                    "DepartureTs": 1765799100,
                    "ArrivalTs": 1765813800
                },
                "Duration": {
                    "TotalMinutes": 245,
                    "Formatted": "4h 5m"
                },
                "Stops": 1,
                "Price": {
                    "Amount": 950000,
                    "Currency": "IDR"
                },
                "SeatsAvailable": 41,
                "BestValueScore": 3877.5510204081634
            }
            // ... other
        ],
        "metadata": {
            "total_result": 3,
            "sort": {
                "key": "price",
                "order": "desc"
            }
        }
    },
    "message": "success searching flight successfully"
}
```

**Response** `Bad request (400)`:
```json
{
    "data": null,
    "message": "please check your request",
    "error": {
        "errors": [
            {
                "field": "origin",
                "tag": "min",
                "value": "3"
            }
        ]
    }
}
```

**Response** `Internal Server Error (500)`:
```json
{
    "data": null,
    "message": "please try again"
}
```