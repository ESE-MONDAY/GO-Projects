Fair call. Let's strip away the corporate buzzwords and generic documentation fluff. Here is a clean, punchy, developer-focused README that reads like it was written by an engineer, not an instruction manual.

---

# Movies CRUD with Live Prometheus & Grafana

A lightning-fast, zero-database REST API built to showcase a production-grade monitoring pipeline. The API stores data in a native Go slice and uses custom middleware to track metrics. Everything runs inside a single Docker Compose network to bypass messy local firewall issues.

## The Architecture

The entire app stack is sandboxed inside Docker to ensure zero configuration friction:

* **Go App (Port 8080):** Uses Gorilla Mux for routing, `log/slog` for structured JSON logs, and standard middleware to time incoming HTTP requests.
* **Prometheus (Port 9090):** Scrapes the Go app's `/metrics` endpoint every 2 seconds.
* **Grafana (Port 3000):** Visualizes the Prometheus time-series data using dynamic line charts.

---

## File Tree

```text
.
├── docker-compose.yml     # Orchestrates the Go app, Prom, and Grafana
├── Dockerfile             # Multi-stage build for a tiny Go container image
├── go.mod                 # Module dependencies
├── main.go                # API logic & monitoring middleware
└── prometheus.yml         # Tells Prometheus where to scrape data

```

---

## Quick Start

### 1. Boot the Stack

Spin up all three containers in the background. Docker will automatically build the Go binary:

```bash
docker compose up --build -d

```

Verify they are all running: `docker ps`

### 2. Smash it with Traffic

To populate your charts with actual latency distribution data, run a stress test using `hey`:

```bash
hey -n 30000 -c 100 http://localhost:8080/movies

```

---

## API Spec

| Method | Endpoint | Description | Payload Example |
| --- | --- | --- | --- |
| **GET** | `/movies` | Get all movies | *None* |
| **GET** | `/movies/{id}` | Get single movie | *None* |
| **POST** | `/movies` | Add a movie | `{"title":"Interstellar","director":"Nolan","year":2014}` |
| **PUT** | `/movies/{id}` | Update a movie | `{"title":"The Matrix 2","director":"Wachowski","year":2003}` |
| **DELETE** | `/movies/{id}` | Delete a movie | *None* |
| **GET** | `/metrics` | Raw Prometheus data | *None* |

### Quick Test Commands

```bash
# Get all movies
curl http://localhost:8080/movies

# Create a movie
curl -X POST -H "Content-Type: application/json" -d '{"title":"Inception","director":"Nolan","year":2010}' http://localhost:8080/movies

```

---

## Logging Format (`log/slog`)

Every incoming HTTP request dumps a structured JSON line to stdout. It looks like this:

```json
{"time":"2026-06-14T23:06:12Z","level":"INFO","msg":"http_request","method":"GET","path":"/movies","status":200,"latency":852300}

```

---

## Grafana Dashboard setup

1. Go to `http://localhost:3000` (User: `admin` / Pass: `admin`).
2. Navigate to **Connections** -> **Data Sources** -> **Add data source** and select **Prometheus**.
3. Set the connection URL to: `http://prometheus:9090` and hit **Save & Test**.
4. Create a new panel, select **Time series** visualization, set the time window to **Last 5 minutes**, and paste any of these production queries:

### Useful PromQL Queries

#### Real-time p95 Latency (Slowest 5% of users)

```text
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[1m])) by (le))

```

#### Real-time p99 Latency (Worst-case performance anomalies)

```text
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[1m])) by (le))

```

#### Throughput Volume by HTTP Status Code (Spotting 404s/500s)

```text
sum(rate(http_request_duration_seconds_count[1m])) by (status)

```

#### Go Goroutine Tracker (Monitoring memory/thread leaks)

```text
go_goroutines

```