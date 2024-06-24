# Route53 module for Caddy

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

This module supports all the credential configuration methods described in the [AWS Developer Guide](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials), such as `Environment Variables`, `Shared configuration files`, the `AWS Credentials file` located in `.aws/credentials`, and `Static Credentials`. You may also pass in static credentials directly (or via caddy's configuration).

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "route53",
        "max_retries": 10, // optional
        "profile": "real-profile", // optional
        "region": "us-east-1", // optional
        "access_key_id": "AKI...", // optional
        "secret_access_key": "wJa...", // optional
        "session_token": "TOKEN..." // optional
      }
    }
  }
}
```

or with the Caddyfile:

```
tls {
  dns route53 {
    max_retries 10 // optional
    profile "real-profile" // optional
    access_key_id "AKI..." // optional
    secret_access_key "wJa..." // optional
    session_token "TOKEN..." // optional
    region "us-east-1" // optional
  }
}
```

The following IAM policy is a minimal working example to give `libdns` permissions to manage DNS records:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": [
        "route53:ListResourceRecordSets",
        "route53:GetChange",
        "route53:ChangeResourceRecordSets"
      ],
      "Resource": [
        "arn:aws:route53:::hostedzone/ZABCD1EFGHIL",
        "arn:aws:route53:::change/*"
      ]
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": ["route53:ListHostedZonesByName", "route53:ListHostedZones"],
      "Resource": "*"
    }
  ]
}
```

For more information, refer to [libdns/route53](https://github.com/libdns/route53).
