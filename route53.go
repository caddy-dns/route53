package cloudflare

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
			p := &Provider{new(route53.Provider)}
			// Initialize the AWS client session
			err := p.NewSession()
			if err != nil {
				return nil
			}
			return p
		},
	}
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// route53 {
//     max_retries <int>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	repl := caddy.NewReplacer()
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "max_retries":
				if d.NextArg() {
					p.Provider.MaxRetries, _ = strconv.Atoi(repl.ReplaceAll(d.Val(), ""))
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

// Interface guard
var _ caddyfile.Unmarshaler = (*Provider)(nil)
