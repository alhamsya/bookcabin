# <img align="right" src="https://avatars.githubusercontent.com/u/56905970?s=60&v=4" alt="bookcabin" title="bookcabin" /> bookcabin

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