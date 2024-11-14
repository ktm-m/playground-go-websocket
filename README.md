# Playground Go WebSocket

## Description

This is a playground project to implement Gorilla WebSocket with Echo framework and Gin framework.

## Technologies Stack

- Go
- Echo framework
- Gin framework
- Gorilla WebSocket
- Docker

## Prerequisites

- Go
- Docker
- Makefile

## Environment Variables

Create a `./config/config.yaml` file in the root directory and add the following variables:

```yaml
app:
  name: "playground-go-websocket"
  echoPort: "8081"
  ginPort: "8082"
  version: "1.0.0"
  trustProxies:
    - "127.0.0.1"
  upgrader:
    readBufferSize: 1024
    writeBufferSize: 1024
    checkOrigin: true

log:
  level: "debug"
```

## How to run

Run the following commands:

```bash
make run
```

## How to run with Docker Compose

Run the following commands:

```bash
make run-build
```

## How to stop Docker Compose

Run the following commands:

```bash
make stop-build
```
## How to use this project

Open your browser and go to the following URLs:

- [http://localhost:8081/echo/html](http://localhost:8081/echo/html)
