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

## Example Caddyfile

This is a very minimal Caddyfile example using the `caddy-dns/route53` plugin

```
mysite.io {
  tls myemail@mysite.io {
    dns route53 {
      max_retries 10
    }
  }
```
