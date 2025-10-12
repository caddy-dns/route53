# Docker Test for caddy-dns-route53

This directory contains an automated test environment for validating the caddy-dns-route53 module with real ACME DNS challenges.

## What It Tests

The test suite:
1. Builds Caddy with xcaddy using the local module source
2. Configures Caddy to use Route53 DNS challenge for TLS certificates
3. Requests certificates from Let's Encrypt **staging** environment for test domains
4. Validates that certificates are issued successfully
5. Tests HTTPS endpoints with curl (validates staging certificates)

Test domains used:
- `caddydns-wildtest.<your-zone>` (regular domain)
- `*.caddydns-wildtest.<your-zone>` (wildcard domain)

## Prerequisites

- Docker installed and running
- AWS Route53 hosted zone with proper permissions
- AWS credentials with Route53 access

## Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your AWS credentials and test zone:
   ```bash
   AWS_ACCESS_KEY_ID=your-access-key
   AWS_SECRET_ACCESS_KEY=your-secret-key
   TEST_ZONE=example.com
   ```

   **Important**: Use a real hosted zone you control. The test will create and delete DNS records.

## Running the Test

Execute the test script:
```bash
./run.sh
```

The script will:
1. Validate your `.env` configuration
2. Build the Docker image (uses `caddy:2.10.2-builder-alpine` with xcaddy pre-installed)
3. Run the container with your AWS credentials (no port mapping needed for DNS challenge)
4. Build Caddy with the local caddy-dns-route53 module
5. Start Caddy in background and request certificates via DNS challenge
6. Wait up to 3 minutes for certificate issuance
7. Test HTTPS endpoints with curl
8. Keep container running (Ctrl+C to stop)

## What Happens Inside the Container

1. **Caddy Build**: Uses xcaddy to build Caddy with your local module:
   ```
   xcaddy build --with github.com/caddy-dns/route53=/caddy-dns-src
   ```

2. **Configuration**: Generates Caddyfile from template with your environment variables

3. **Certificate Issuance**: Caddy requests TLS certificates using Route53 DNS challenge

4. **Validation**: Tests endpoints with curl using `--resolve` to override DNS

## Expected Output

Successful test output looks like:
```
=== Building Caddy with xcaddy ===
Caddy built successfully!

=== Starting Caddy ===
=== Waiting for certificate issuance ===
✓ Both certificates found! (2 .crt files)

=== Testing HTTPS endpoints ===
Testing: caddydns-wildtest.example.com
✓ SUCCESS: caddydns-wildtest.example.com responded correctly
  ✓ Certificate validated successfully

Testing: www.caddydns-wildtest.example.com (wildcard)
✓ SUCCESS: www.caddydns-wildtest.example.com responded correctly
  ✓ Certificate validated successfully
```

## Troubleshooting

### Container fails with "Parent directory not mounted"
- The run.sh script should automatically mount the parent directory
- Verify Docker has permissions to access your project directory

### Certificate issuance times out
- Check AWS credentials are valid and have Route53 permissions
- Verify TEST_ZONE exists in your AWS account
- Check Route53 rate limits (Let's Encrypt has rate limits too)
- Review Caddy logs shown after timeout

### HTTPS tests fail but certificates issued
- This may indicate a Caddy routing issue
- Check the generated Caddyfile in container output
- Verify domain names match your TEST_ZONE

### AWS permission errors
Ensure your AWS credentials have these Route53 permissions:
- `route53:ListHostedZones`
- `route53:GetChange`
- `route53:ChangeResourceRecordSets`

## Testing from Inside the Container

While the container is running, you can connect to it and test manually:

```bash
# In another terminal, connect to the running container
docker exec -it caddy-dns-route53-test sh

# Inside the container, test the endpoints
curl --resolve caddydns-wildtest.example.com:443:127.0.0.1 https://caddydns-wildtest.example.com
curl --resolve www.caddydns-wildtest.example.com:443:127.0.0.1 https://www.caddydns-wildtest.example.com
```

**Note**: No `-k` flag needed inside the container since staging root certificates are installed. Ports are not mapped to the host, so testing must be done from inside the container.

## Configuration Options

The Caddyfile template uses:
- **Let's Encrypt Staging**: Uses `acme-staging-v02.api.letsencrypt.org` to avoid production rate limits
  - Staging root certificates (`letsencrypt-stg-root-x1`, `letsencrypt-stg-root-x2`) are installed in the Docker image
  - This allows proper certificate validation without `-k` flag inside the container
  - Certificates are from staging CA but fully validated for testing purposes
- **Route53 settings**:
  - `wait_for_route53_sync true` - Waits for Route53 internal sync (recommended)
  - `skip_route53_sync_on_delete false` - Waits for sync on cleanup to avoid "record already exists" errors
- **No AWS region required**: Route53 is a global service

See the main [README](../README.md) for full configuration options.

## Cleanup

DNS records are automatically cleaned up by Caddy/ACME after certificate issuance. The test records (`_acme-challenge.caddydns-wildtest.*`) are temporary.

To stop the container, press `Ctrl+C`.

## Notes

- **Uses Let's Encrypt Staging**: Certificates are from staging CA to avoid production rate limits
  - Staging root certificates are installed in the Docker image for proper validation
  - Certificates are fully validated inside the container (no `-k` flag needed)
  - From host machine, you need `-k` flag since staging roots aren't in your system trust store
  - For production testing, change `acme_ca` in Caddyfile.template to production endpoint and remove staging cert installation
- Certificate files stored in `/root/.local/share/caddy/certificates` inside container
- Debug logs available in `docker-test/debug/` on host machine:
  - `caddy.log` - Caddy structured logs with DEBUG level
  - `caddy-stdout.log` - Caddy stdout/stderr output
- Container runs in interactive mode and is removed on exit (`--rm` flag)
- No ports exposed to host - DNS challenge doesn't require it
