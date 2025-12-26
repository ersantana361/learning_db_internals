# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go workspace files
COPY go.work go.work.sum* ./

# Copy all modules
COPY apps/api ./apps/api
COPY packages ./packages
COPY projects ./projects

# Build the API
WORKDIR /app/apps/api
RUN go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ENV PORT=8080

CMD ["./server"]
