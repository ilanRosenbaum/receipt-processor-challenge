# Receipt Processor Service

A RESTful web service that processes store receipts and calculates reward points based on specific rules. Built with Go, this service provides an API for receipt processing and points calculation.

## Features

- Receipt validation and processing
- Points calculation based on multiple rules
- In-memory storage with thread-safe operations
- RESTful API with JSON responses
- Test coverage including integration tests

## Prerequisites

Primary requirements:

- Go 1.19 or higher
- Git

Alternative requirement (if not using Go):

- Git
- Docker

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/ilanRosenbaum/receipt-processor-challenge.git
cd receipt-processor
```

2. Choose your preferred method to run the service:

### Using Go (Recommended)

```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

### Alternative: Using Docker

```bash
# Build the Docker image
docker build -t receipt-processor .

# Run the container
docker run -p 8080:8080 receipt-processor
```

The service will start on port 8080 by default.

## API Documentation

The API specification is defined in OpenAPI 3.0 format. See [api.yml](./api.yml) for the complete API documentation, including:

- Available endpoints
- Request/response formats
- Data schemas
- Example payloads

## Points Calculation Rules

Points are awarded based on the following rules:

1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. Points for items with descriptions of length multiple of 3
6. 6 points if the day in the purchase date is odd
7. 10 points if the time of purchase is between 2:00pm and 4:00pm

## Development

### Project Structure

```
.
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── models/          # Data models
│   ├── service/         # Business logic and validation
│   └── store/           # Data storage
├── main.go              # Application entry point
├── api.yml             # API specification
└── README.md
```

### Running Tests

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test ./... -cover
```
