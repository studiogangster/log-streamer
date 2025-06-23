# 🚀 log-streamer

A blazing-fast REST microservice that mimics Unix's `head -f` (head with follow) for large log files.  
Easily stream, preview, and monitor logs in real-time—right from your browser or any HTTP client!

---

![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue?logo=go)
![Dockerized](https://img.shields.io/badge/Docker-ready-blue?logo=docker)
![MIT License](https://img.shields.io/badge/license-MIT-green)

---

## ✨ Features

- **Real-time log streaming**: Follow log files as they grow, just like `tail -f` or `head -f`.
- **REST API**: Simple HTTP endpoints for easy integration.
- **Handles large files**: Efficiently reads and streams even massive logs.
- **Docker support**: Run anywhere with a single command.
- **Easy to use**: Minimal setup, instant results.

---

## 🚀 Quick Start

### 1. Clone & Build

```bash
git clone https://github.com/yourusername/log-streamer.git
cd log-streamer
go build -o log_reader main.go
```

### 2. Run with Docker

```bash
docker build -t log-streamer .
docker run -p 8080:8080 -v /path/to/logs:/logs log-streamer
```

### 3. Run Locally

```bash
go run main.go
```

---

## 📡 API Usage

### Stream the first N lines and follow a log file

```bash
curl "http://localhost:8080/read?file=/logs/yourfile.log&lines=20&follow=true"
```

- `file`: Path to the log file (inside the container or local path)
- `lines`: Number of lines to read from the start (default: 10)
- `follow`: Set to `true` to keep streaming new lines

### Example Response

```json
{
  "lines": [
    "2024-06-23 12:00:00 INFO Starting service...",
    "2024-06-23 12:00:01 INFO Listening on port 8080",
    "...",
    "2024-06-23 12:00:10 INFO Service ready."
  ]
}
```

---

## 🛠️ Project Structure

```
.
├── main.go              # Entry point
├── readlogs/            # Log reading logic
├── sanitize/            # Input sanitization
├── Dockerfile           # Docker support
├── docker-compose.yaml  # Docker Compose config
└── test_log_files/      # Sample log files
```

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!  
Feel free to open an issue or submit a pull request.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

> Made with ❤️ in Go
