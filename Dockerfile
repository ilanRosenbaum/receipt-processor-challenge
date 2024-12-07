FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o receipt-processor
EXPOSE 8080

# Run the application
CMD ["./receipt-processor"]
