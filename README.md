Route53 module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records in Route53 Hosted zones.

## Caddy module name

```
dns.providers.route53
```


## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/route53) for important information about credentials and an IAM policy example.

## Building

To compile this Caddy module, follow the steps describe at the [Caddy Build from Source](https://github.com/caddyserver/caddy#build-from-source) instructions and import the `github.com/caddy-dns/route53` plugin

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "route53",
        "max_retries": 10,
        "aws_profile": "real-profile"
      }
    }
  }
}
```

or with the Caddyfile:

```
tls {
  dns route53 {
    max_retries 10,
    aws_profile "real-profile"
  }
}
```
