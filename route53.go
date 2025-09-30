package route53

import (
	"strconv"
	"time"

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
	p.Provider.Profile = repl.ReplaceAll(p.Provider.Profile, "")
	p.Provider.AccessKeyId = repl.ReplaceAll(p.Provider.AccessKeyId, "")
	p.Provider.SecretAccessKey = repl.ReplaceAll(p.Provider.SecretAccessKey, "")
	p.Provider.SessionToken = repl.ReplaceAll(p.Provider.SessionToken, "")
	p.Provider.Region = repl.ReplaceAll(p.Provider.Region, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	route53 {
//		region <string>
//		profile <string>
//		access_key_id <string>
//		secret_access_key <string>
//		session_token <string>
//		max_retries <int>
//		max_wait_dur <int>
//		wait_for_propagation <bool>
//		hosted_zone_id <string>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "wait_for_propagation":
				if d.NextArg() {
					if wait, err := strconv.ParseBool(d.Val()); err == nil {
						p.Provider.WaitForPropagation = wait
					}
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "max_wait_dur":
				if d.NextArg() {
					if dur, err := strconv.ParseInt(d.Val(), 10, 64); err == nil {
						p.Provider.MaxWaitDur = time.Duration(dur)
					}
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "max_retries":
				if d.NextArg() {
					p.Provider.MaxRetries, _ = strconv.Atoi(d.Val())
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "profile":
				if d.NextArg() {
					p.Provider.Profile = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "aws_profile":
				if d.NextArg() {
					p.Provider.Profile = d.Val()
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
			case "session_token":
				if d.NextArg() {
					p.Provider.SessionToken = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "token":
				if d.NextArg() {
					p.Provider.SessionToken = d.Val()
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
			case "hosted_zone_id":
				if d.NextArg() {
					p.Provider.HostedZoneID = d.Val()
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
