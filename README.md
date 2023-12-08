# DA Checker

## Introduction

DA Checker is an open source tool designed for monitoring and validating data associated with validators da nodes' metrics. 
Uses Go backend, MySQL database, and a serverless deployment model.

## Features

- **Metrics Scraping and Processing**: Regularly scrapes metrics from a Prometheus endpoint, processes them, and updates the database with the latest status of validators and their peer IDs.
- **RESTful API**: Offers a set of RESTful endpoints to interact with validator data using the Gin web framework.
- **Database Integration**: Utilizes a MySQL database for persistent storage of validator data and historical peer ID statuses.
- **Serverless Deployment**: Facilitates easy and scalable deployment using serverless architecture, enhancing maintainability.

## Package Structure

- `pkg/database`: Contains database connection setup and initialization logic. Implements a singleton pattern.
- `pkg/models`: Defines the data models used throughout the application, including Validator, PeerID, and PeerIDHistory.
- `pkg/handlers`: Implements the RESTful API logic, handling HTTP requests and responses.
- `pkg/metrics`: Houses the core business logic for scraping Prometheus metrics, processing them, and updating the database.
- `pkg/repository`: Implements the repository pattern for abstracting database access.
- `pkg/vpi`: Contains the logic for Validator-PeerID matching. The end result is scraped by the agent, so it can be used in combination with otlp metrics sent by validators. Useful for variable based dashboards in Grafana.

## Database

It stores validators' data, peer IDs, and historical statuses, facilitating efficient data retrieval and manipulation in later dashboards.

## Getting Started

### Prerequisites

- Go 1.21.1
- MySQL
- Serverless Framework

### Installation

Clone the repository & install:
   ```bash
   git clone https://github.com/celestiaorg/da-checker.git
   go mod tidy
   go run -v cmd/main.go
   ```
Run locally:

1. Define your env vars. The whole file is defined in .env file 
2. If you want to load .env automatically, you can use [godotenv](github.com/joho/godotenv)
