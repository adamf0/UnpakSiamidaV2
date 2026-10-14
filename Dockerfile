# ==========================================================
# STAGE 1: BUILD STAGE
# Multi-stage build for Go application
# ==========================================================
FROM golang:bookworm AS builder

WORKDIR /app

# Enable CGO_ENABLED=0 for a fully static, standalone binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Compile binary with optimizations:
# -ldflags="-w -s" removes debug information & symbols, significantly reducing binary size
RUN go build -ldflags="-w -s" -o /app/app .

# ==========================================================
# STAGE 2: RUNTIME STAGE
# Production lightweight runtime environment
# ==========================================================
FROM debian:bookworm-slim

WORKDIR /app

# Install essential runtime dependencies:
# - wkhtmltopdf & fonts: for generating PDF reports (Berita Acara & KTS)
# - ca-certificates: for secure outbound HTTPS requests
# - tzdata: for correct timezone handling
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    wkhtmltopdf \
    fonts-dejavu-core \
    fontconfig \
    && rm -rf /var/lib/apt/lists/*

# Set default timezone to Western Indonesia Time (WIB)
ENV TZ=Asia/Jakarta

# Create required directories for uploads and runtime assets
RUN mkdir -p /app/upload/tmp \
             /app/upload/prokerv2 \
             /app/upload/proker_v3 \
             /app/assets \
             /app/modules

# Copy compiled binary from builder
COPY --from=builder /app/app /app/app

# Copy runtime assets (logos) and module templates (HTML templates for PDF generation)
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/modules /app/modules

# Ensure binary is executable
RUN chmod +x /app/app

# Application listens on port 3000
EXPOSE 3000

# Start server
CMD ["/app/app"]
