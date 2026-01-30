# Go-Stream: Real-Time BTC/USDT Trading Stack

A high-performance, containerized streaming architecture designed to capture, store and visualize real-time BTC/USDT market data. This project demonstrates a complete data engineering pipeline using Go, Kafka and QuestDB, featuring a professional TradingView dashboard.

## System Architecture

The project follows a microservices-based "Event-Driven" architecture:

* Producer (Go): Connects to the OKX WebSocket API, captures real-time BTC/USDT "tick" data and publishes it to a Kafka topic.
* Consumer (Go): Subscribes to the Kafka topic and performs high-speed batch writes into QuestDB.
* API Backend (Go): Aggregates raw tick data into 1-minute (M1) candlesticks using QuestDB's SAMPLE BY SQL syntax and serves it via http://localhost:3000/candles
* Frontend (Go): A dedicated static file server providing the Goosey LLC - QuestDB View professional terminal on http://localhost:8080

## Tech Stack

- Language: Go (Golang 1.24+)
- Message Broker: Apache Kafka
- Database: QuestDB (Time-series optimized)
- Frontend: TradingView Lightweight Charts v4
- Containerization: Docker
- API Framework: Fiber (Go)

## Installation & Deployment

### Prerequisites
- Dockerfile supported

### Fast Track
Clone the repository and spin up the entire stack with a Docker command.

The system will automatically initialize Kafka topics, QuestDB tables and start the data flow from OKX.

## Service Endpoints

| Service | URL | Description |
| :--- | :--- | :--- |
| Terminal UI | http://localhost:8080 | Pro Trading Dashboard (Goosey LLC) |
| Backend API | http://localhost:3000/candles | JSON Candlestick Data Endpoint |
| QuestDB UI | http://localhost:9000 | Database Console & SQL Editor |

## Time-Series Aggregation

The system utilizes QuestDB's specialized SQL syntax to generate candlesticks on-the-fly from raw tick data:

SELECT
    (extract(epoch from timestamp))::long as time,
    first(price) as open,
    max(price) as high,
    min(price) as low,
    last(price) as close
FROM btc_ticks
SAMPLE BY 1m;

## Key Features
- Real-time Flow: Sub-millisecond data processing from OKX to QuestDB
- Modern UI: Professional dark-mode terminal with "Goosey LLC" branding and live pulse indicators
- Scalable Design: Kafka-driven architecture allows for adding multiple pairs (ETH, SOL, etc.) easily
- Dockerized: Fully isolated environment, no local database or Kafka installation required

---
Created as an example project for high-frequency data streaming in Go and built with focus on network performance and time-series efficiency.
