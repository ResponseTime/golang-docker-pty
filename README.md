# 🐚 Remote Shell Platform (Go + Docker)

A multi-tenant remote shell platform built with Go that dynamically provisions isolated Docker containers per session and provides full terminal access over TCP connections.

## ✨ Features

- 🔒 **Multi-Tenant Architecture**  
  Handles multiple client sessions simultaneously, each in an isolated Docker container.

- ⚙️ **Dynamic Docker Provisioning**  
  Automatically creates and starts a fresh Ubuntu-based Docker container for each new incoming connection.

- 💻 **PTY Support**  
  Uses [creack/pty](https://github.com/creack/pty) to provide fully interactive Bash shells with support for:
  - Shell prompts  
  - Line editing  
  - Signal handling (Ctrl+C, etc.)

- 🔁 **Bidirectional I/O Streaming**  
  Bridges the TCP connection with the container’s PTY master using `io.Copy` for low-latency, real-time interaction.

- ⏱ **Automatic Session Expiry**  
  Uses `context.WithTimeout` to enforce time limits on sessions and automatically clean up containers after timeout.

## 🚀 How It Works

1. Listens on a TCP port for incoming client connections.
2. Upon connection, spawns a new Docker container running an interactive Bash shell.
3. Sets up a PTY inside the container and bridges it with the client over TCP.
4. Automatically tears down the session and container after the allocated time.

## 🛠 Tech Stack

- **Language:** Go
- **Containerization:** Docker
- **Terminal Emulation:** [creack/pty](https://github.com/creack/pty)
- **Concurrency:** Goroutines and context for timeouts

## 🧪 Usage

### Prerequisites

- Docker installed and running
- Go installed (`go 1.20+` recommended)

### Run

```bash
go run main.go
