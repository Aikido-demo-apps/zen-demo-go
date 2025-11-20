FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

RUN go install github.com/DataDog/orchestrion

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -toolexec "orchestrion toolexec" -o main .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy binary and static files
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 3000

CMD ["./main"]
