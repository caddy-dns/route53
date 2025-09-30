# Route53 module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy).
It lets Caddy read and manipulate DNS records hosted by Route53, to obtain TLS certificates
with the [DNS challenge](https://caddyserver.com/docs/automatic-https#dns-challenge).

> [!NOTE]
> Caddy 2.10 upgraded to libdns 1.0, which breaks compatibility with older DNS providers.
> To use Caddy 2.10 or newer, install **version 1.6** or later.
> For earlier Caddy versions, use a corresponding older module release.  


## Caddy module name

```
dns.providers.route53
```

See [the associated README in the libdns package](https://github.com/libdns/route53) for important information about credentials and an IAM policy example.

## Building

To compile this Caddy module, you can use [xcaddy](https://github.com/caddyserver/xcaddy) the following way:

```bash
$ xcaddy build \
    --with github.com/caddy-dns/route53@version
```

Replace `version` with the desired version.

For more advanced cases follow the steps described at the [Caddy Build from Source](https://github.com/caddyserver/caddy#build-from-source) instructions and import the `github.com/caddy-dns/route53` plugin

## Authenticating

This module uses [libdns/route53](https://github.com/libdns/route53) which uses the [AWS SDK for Go v2](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/welcome.html) and supports all its authentication methods:
- Static credentials (directly in Caddy config)
- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`)
- EC2 instance roles
- Shared configuration files (`~/.aws/config`, `~/.aws/credentials`)

The following AWS IAM policy is a minimal working example to give `libdns` permissions to manage DNS records of hosted zone `ZABCD1EFGHIL` (replace with your own):

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

## Configuration

To use this module for the ACME DNS challenge, configure the ACME issuer in your Caddyfile like so:

```caddy
tls {
  dns route53 {
    access_key_id "AKI..."        # required unless using EC2 Roles
                                  #   or $AWS_ACCESS_KEY_ID is defined

    secret_access_key "wJa..."    # required unless using EC2 Roles
                                  #   or $AWS_SECRET_ACCESS_KEY is defined

    region "us-east-1"            # optional in 1.6
                                  #   defaults to $AWS_REGION
                                  #   or us-east-1 (Route53 default)

    wait_for_propagation false    # optional, defaults to true in libdns/route53 1.6
                                  #   note that this is internal AWS propagation,
                                  #   not external DNS

    max_wait_dur 60               # optional, defaults to 60 in libdns/route53 1.6
    max_retries 5                 # optional, defaults to 5 in libdns/route53 1.6
    profile "real-profile"        # optional, rarely needed, defaults to $AWS_PROFILE
    session_token "TOKEN..."      # optional, rarely needed, defaults to $AWS_SESSION_TOKEN

    hosted_zone_id ZABCD1EFGHIL   # optional
                                  # hosted_zone_id would be discovered from AWS otherwise
  }
}
```

> [!NOTE]
> As of 2025, the `region` option is rarely needs to be changed because most AWS Route53 regions use [the same endpoints](https://docs.aws.amazon.com/general/latest/gr/r53.html) as `us-east-1`. It is only required for AWS GovCloud and the China Beijing and Ningxia regions.

### JSON configuration example (see above for comments):

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "route53",
        "region": "us-east-1",
        "access_key_id": "AKI...",
        "secret_access_key": "wJa...",
        "max_retries": 10,
        "profile": "real-profile",
        "session_token": "TOKEN...",
        "max_wait_dur": 60,
        "wait_for_propagation": false,
        "hosted_zone_id": "ZABCD1EFGHIL"
      }
    }
  }
}
```

When using AWS EC2 instance roles, a minimal Caddy configuration may look like this:

```caddy
*.caddyexample.example.com {
  tls {
      dns route53 {
      }
  }
}
```

When using AWS access keys, the configuration becomes:

```caddy
*.caddyexample.example.com {
  tls {
      dns route53 {
        access_key_id "AKI..."
        secret_access_key "wJa..."
      }
  }
}
```

## More information

This module is extremely compact and primarily does configuration - the actual Route53 calls are made by [libdns/route53](https://github.com/libdns/route53). Refer to that project for more information.

