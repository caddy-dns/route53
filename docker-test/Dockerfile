FROM caddy:2.11.2-builder-alpine

# NOTE: Build context is the parent directory (caddy-dns-route53/), not docker-test/
# This is set via: docker build -f docker-test/Dockerfile <parent-dir>

# Install additional dependencies
RUN apk add --no-cache \
    bash \
    curl \
    gettext \
    ca-certificates \
    tini

# Install Let's Encrypt staging root certificate for testing
RUN curl -sL https://letsencrypt.org/certs/staging/letsencrypt-stg-root-x1.pem -o /usr/local/share/ca-certificates/letsencrypt-stg-root-x1.crt \
    && curl -sL https://letsencrypt.org/certs/staging/letsencrypt-stg-root-x2.pem -o /usr/local/share/ca-certificates/letsencrypt-stg-root-x2.crt \
    && update-ca-certificates

# Copy the caddy-dns-route53 module source (entire parent directory)
WORKDIR /build
RUN xcaddy build && rm ./caddy # warm up cache in docker layer

COPY ./ /build/caddy-dns-route53/

RUN xcaddy build \
    --with github.com/caddy-dns/route53=/build/caddy-dns-route53 \
    && mv caddy /usr/bin/caddy \
    && caddy version

# Clean up build artifacts (keep module source for reference)
RUN rm -rf /tmp/buildenv_* /root/.cache/go-build

# Create working directory for runtime
WORKDIR /test

# Copy test scripts (from docker-test/ subdirectory in build context)
COPY docker-test/inside-docker.sh /test/
COPY docker-test/Caddyfile.template /test/
RUN chmod +x /test/inside-docker.sh

# Use tini as entrypoint to properly handle signals (Ctrl+C)
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/test/inside-docker.sh"]
