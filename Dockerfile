FROM golang:1.24-alpine@sha256:12c199a889439928e36df7b4c5031c18bfdad0d33cdeae5dd35b2de369b5fbf5 AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app


# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

RUN go install github.com/AikidoSec/firewall-go/cmd/zen-go@v0.1.0

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -toolexec "zen-go toolexec" -o main .

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
