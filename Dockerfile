# Build stage
FROM golang:1.24.3 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o expense-tracker

# Run stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/expense-tracker .
COPY --from=builder /app/templates ./templates
COPY .env ./
EXPOSE 8000
CMD ["./expense-tracker"]
