FROM golang:1.26-alpine@sha256:ce864e7223ac17b1775e6fd0b4c0db580c2eb50e7953a427916379e4b92a1628 AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app


# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go tool zen-go go build -o main .

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
