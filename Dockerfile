# Stage 1: Build environment
FROM golang:1.24.1-bookworm AS builder

WORKDIR /app

# Create non-root user and group in build stage
RUN addgroup --gid 222 --system app \
    && adduser --uid 222 --system --group app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Runtime environment
FROM scratch

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app /app/app

# Copy user and group definitions from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set ownership (optional, as scratch has no chown command)
# Since we're using a non-root user, ensure the binary is accessible
USER 222:222

EXPOSE 8000
CMD ["/app/app"]
