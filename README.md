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
        "profile": "real-profile", // defaults to $AWS_PROFILE
        "region": "us-east-1", // required if $AWS_REGION is not defined
        "access_key_id": "AKI...", // required if $AWS_ACCESS_KEY_ID is not defined
        "secret_access_key": "wJa...", // required if $AWS_SECRET_ACCESS_KEY is not defined
        "session_token": "TOKEN...", // defaults to $AWS_SESSION_TOKEN (optional)
        "max_wait_dur": 60, // propagation wait duration in seconds (optional)
        "wait_for_propagation": false, // wait for records to propagate (optional)
        "hosted_zone_id": "ZABCD1EFGHIL" // AWS hosted zone ID to update (optional)
      }
    }
  }
}
```

or with the Caddyfile:

```caddy
tls {
  dns route53 {
    max_retries 10 // optional
    profile "real-profile" // defaults to $AWS_PROFILE
    access_key_id "AKI..." // required if $AWS_ACCESS_KEY_ID is not defined
    secret_access_key "wJa..." // required if $AWS_SECRET_ACCESS_KEY is not defined
    session_token "TOKEN..." // defaults to $AWS_SESSION_TOKEN (optional)
    region "us-east-1" // required if $AWS_REGION is not defined
    max_wait_dur 60, // propagation wait duration in seconds (optional)
    wait_for_propagation false // wait for records to propagate (optional)
    hosted_zone_id ZABCD1EFGHIL // AWS hosted zone ID to update (optional)
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
