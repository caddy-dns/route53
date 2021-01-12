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

// Provision implements the Provisioner interface to initialize the AWS Client sessions
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.AWSProfile = repl.ReplaceAll(p.Provider.AWSProfile, "")

	// Initialize the AWS client session
	return p.NewSession()
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// route53 {
//     max_retries <int>
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
