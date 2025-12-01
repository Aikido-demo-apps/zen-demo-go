FROM golang:1.24-alpine@sha256:0466223b8544fb7d4ff04748acc4d75a608234bf4e79563bff208d2060c0dd79 AS builder

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
FROM alpine:3.22.2@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy binary and static files
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 3000

RUN apk add --no-cache shadow && \
    useradd -U -u 1000 appuser && \
    chown -R 1000:1000 /app
USER 1000

CMD ["./main"]
