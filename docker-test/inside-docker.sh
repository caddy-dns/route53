#!/usr/bin/env bash
set -euo pipefail

echo "=== Inside Docker Container ==="
echo "Working directory: $(pwd)"
echo ""

# Verify Caddy binary is available
echo "=== Verifying Caddy Installation ==="
if ! command -v caddy >/dev/null 2>&1; then
    echo "ERROR: Caddy binary not found in PATH!"
    exit 1
fi

caddy version
echo "✓ Caddy binary ready"
echo ""

echo "=== Preparing Caddyfile ==="
# Substitute environment variables in Caddyfile.template
envsubst < Caddyfile.template > Caddyfile

echo "Generated Caddyfile:"
echo "---"
cat Caddyfile
echo "---"
echo ""

echo "=== Starting Caddy in background ==="
# Verify debug directory is mounted and writable
if [ ! -d "/debug" ]; then
    echo "ERROR: /debug directory not mounted!"
    exit 1
fi

# Test write access
if ! touch /debug/.test 2>/dev/null; then
    echo "ERROR: /debug directory not writable!"
    ls -la /debug
    exit 1
fi
rm -f /debug/.test
echo "Logs will be written to: /debug/"

# Create log file first to avoid race condition
touch /debug/caddy-stdout.log
caddy run --config Caddyfile --adapter caddyfile > /debug/caddy-stdout.log 2>&1 &
CADDY_PID=$!

# Give Caddy a moment to start
sleep 3

# Check if Caddy process is still running
if ! kill -0 $CADDY_PID 2>/dev/null; then
    echo "ERROR: Caddy process died immediately after start!"
    echo "Check logs:"
    cat /debug/caddy-stdout.log 2>/dev/null || echo "No stdout log"
    exit 1
fi

echo "✓ Caddy process running (PID: $CADDY_PID)"

echo ""
echo "=== Waiting for certificate issuance ==="
echo "This may take up to 3 minutes..."

# Wait for certificate to be issued
MAX_WAIT=180
WAITED=0
CERT_ISSUED=false

# Check multiple possible certificate paths (staging vs production)
CERT_PATHS=(
    "/root/.local/share/caddy/certificates"
)

while [ $WAITED -lt $MAX_WAIT ]; do
    sleep 5
    WAITED=$((WAITED + 5))

    # Check all possible certificate paths
    for cert_path in "${CERT_PATHS[@]}"; do
        if [ -d "$cert_path" ]; then
            CERT_COUNT=$(find "$cert_path" -name "*.crt" 2>/dev/null | wc -l)
            # We expect 2 certificates: one for caddydns-wildtest.* and one for *.caddydns-wildtest.*
            if [ "$CERT_COUNT" -ge 2 ]; then
                echo "✓ Both certificates found in $cert_path! ($CERT_COUNT .crt files)"
                CERT_ISSUED=true
                break 2
            elif [ "$CERT_COUNT" -gt 0 ]; then
                echo "  Found $CERT_COUNT certificate(s), waiting for 2..."
            fi
        fi
    done

    if [ "$CERT_ISSUED" = false ]; then
        echo "  Waiting... (${WAITED}s/${MAX_WAIT}s)"
    fi
done

if [ "$CERT_ISSUED" = false ]; then
    echo ""
    echo "WARNING: Certificates not found after ${MAX_WAIT}s"
    echo "Checked paths:"
    for cert_path in "${CERT_PATHS[@]}"; do
        echo "  - $cert_path"
    done
    echo ""
    echo "Last 50 lines of Caddy log:"
    echo "---"
    tail -50 /debug/caddy.log 2>/dev/null || echo "Log file not found"
    echo "---"
fi

echo ""
echo "=== Testing HTTPS endpoints ==="

# Test domains
TEST_DOMAIN="caddydns-wildtest.${TEST_ZONE}"
TEST_WILDCARD="www.caddydns-wildtest.${TEST_ZONE}"

echo ""
echo "Testing: $TEST_DOMAIN"
echo "(Validating staging certificate)"
RESPONSE=$(curl -sv --connect-timeout 5 --max-time 10 --resolve "${TEST_DOMAIN}:443:127.0.0.1" "https://${TEST_DOMAIN}" 2>&1 || true)
if echo "$RESPONSE" | grep -q "Hello from"; then
    echo "✓ SUCCESS: $TEST_DOMAIN responded correctly"
    # Check if certificate was validated
    if echo "$RESPONSE" | grep -q "SSL certificate verify ok"; then
        echo "  ✓ Certificate validated successfully"
    fi
    echo "  Response: $(echo "$RESPONSE" | grep "Hello from")"
else
    echo "✗ FAILED: $TEST_DOMAIN did not respond correctly"
    echo "  Response: $RESPONSE"
fi

echo ""
echo "Testing: $TEST_WILDCARD (wildcard)"
echo "(Validating staging certificate)"
RESPONSE=$(curl -sv --connect-timeout 5 --max-time 10 --resolve "${TEST_WILDCARD}:443:127.0.0.1" "https://${TEST_WILDCARD}" 2>&1 || true)
if echo "$RESPONSE" | grep -q "Hello from"; then
    echo "✓ SUCCESS: $TEST_WILDCARD responded correctly"
    # Check if certificate was validated
    if echo "$RESPONSE" | grep -q "SSL certificate verify ok"; then
        echo "  ✓ Certificate validated successfully"
    fi
    echo "  Response: $(echo "$RESPONSE" | grep "Hello from")"
else
    echo "✗ FAILED: $TEST_WILDCARD did not respond correctly"
    echo "  Response: $RESPONSE"
fi

echo ""
echo "=== Checking certificate details ==="
echo "Checking all possible certificate storage locations:"
for cert_path in "${CERT_PATHS[@]}"; do
    if [ -d "$cert_path" ]; then
        echo "Found certificates in: $cert_path"
        find "$cert_path" -name "*.crt" -o -name "*.key" | head -20
    fi
done

echo ""
echo "=== Test complete! ==="
echo "Caddy is still running. Press Ctrl+C to stop."
echo ""
echo "Full logs available on host at: docker-test/debug/"
echo "  - caddy.log (Caddy structured logs)"
echo "  - caddy-stdout.log (stdout/stderr)"
echo ""
echo "To test manually, connect to this container:"
echo "  docker exec -it caddy-dns-route53-test sh"
echo ""
echo "Then inside the container:"
echo "  curl --resolve caddydns-wildtest.${TEST_ZONE}:443:127.0.0.1 https://caddydns-wildtest.${TEST_ZONE}"
echo "  curl --resolve www.caddydns-wildtest.${TEST_ZONE}:443:127.0.0.1 https://www.caddydns-wildtest.${TEST_ZONE}"
echo ""

# Function to handle shutdown signals
cleanup() {
    echo ""
    echo "=== Shutting down ==="
    if [ -n "$CADDY_PID" ]; then
        echo "Stopping Caddy (PID: $CADDY_PID)..."
        kill -TERM "$CADDY_PID" 2>/dev/null || true
        wait "$CADDY_PID" 2>/dev/null || true
    fi
    exit 0
}

# Trap signals to ensure clean shutdown
trap cleanup SIGTERM SIGINT

# Keep container running by waiting for Caddy process
wait $CADDY_PID
