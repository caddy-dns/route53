package route53

import (
	"strconv"

	"github.com/libdns/route53"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *route53.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.route53",
		New: func() caddy.Module {
			return &Provider{new(route53.Provider)}
		},
	}
}

// Provision implements the Provisioner interface to initialize the AWS Client
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AWSProfile = repl.ReplaceAll(p.Provider.AWSProfile, "")
	p.Provider.AccessKeyId = repl.ReplaceAll(p.Provider.AccessKeyId, "")
	p.Provider.SecretAccessKey = repl.ReplaceAll(p.Provider.SecretAccessKey, "")
	p.Provider.Token = repl.ReplaceAll(p.Provider.Token, "")
	p.Provider.Region = repl.ReplaceAll(p.Provider.Region, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// route53 {
//     max_retries <int>
//     aws_profile <string>
//     access_key_id <string>
//     secret_access_key <string>
//	   token <string>
//     region <string>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "max_retries":
				if d.NextArg() {
					p.Provider.MaxRetries, _ = strconv.Atoi(d.Val())
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "aws_profile":
				if d.NextArg() {
					p.Provider.AWSProfile = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "access_key_id":
				if d.NextArg() {
					p.Provider.AccessKeyId = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "secret_access_key":
				if d.NextArg() {
					p.Provider.SecretAccessKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "token":
				if d.NextArg() {
					p.Provider.Token = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "region":
				if d.NextArg() {
					p.Provider.Region = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
