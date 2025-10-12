#!/usr/bin/env bash
set -euo pipefail

# Get directory of this script (works on macOS bash 3.2+)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]:-$0}" )" && pwd )"
PARENT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"

# Check if .env exists
if [ ! -f "$SCRIPT_DIR/.env" ]; then
    echo "ERROR: .env file not found!"
    echo "Please copy .env.example to .env and fill in your credentials:"
    echo "  cd \"$SCRIPT_DIR\""
    echo "  cp .env.example .env"
    echo "  # Edit .env with your AWS credentials and test zone"
    exit 1
fi

# Source .env to validate it
set -a
. "$SCRIPT_DIR/.env"
set +a

# Validate required variables
if [ -z "${AWS_ACCESS_KEY_ID:-}" ] || [ -z "${AWS_SECRET_ACCESS_KEY:-}" ] || [ -z "${TEST_ZONE:-}" ]; then
    echo "ERROR: Required environment variables not set in .env!"
    echo "Please ensure AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and TEST_ZONE are set."
    exit 1
fi

echo "=== Building Docker image ==="
echo "Building from: $PARENT_DIR"
docker build -t caddy-route53-test -f "$SCRIPT_DIR/Dockerfile" "$PARENT_DIR"

echo ""
echo "=== Running Docker container ==="
echo "Test zone: $TEST_ZONE"
echo "Test domains: caddydns-wildtest.$TEST_ZONE, *.caddydns-wildtest.$TEST_ZONE"
echo "Debug logs: $SCRIPT_DIR/debug/"
echo ""

# Ensure debug directory exists and clean old logs
mkdir -p "$SCRIPT_DIR/debug"
rm -f "$SCRIPT_DIR/debug"/*.log

# Run docker with environment variables and debug mount
# Use -i only (not -t) to avoid TTY requirement when run non-interactively
# Mount debug directory as read-write for logs
docker run --rm -i \
    --name caddy-dns-route53-test \
    -v "$SCRIPT_DIR/debug:/debug:rw" \
    -e AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY \
    -e TEST_ZONE \
    caddy-route53-test
