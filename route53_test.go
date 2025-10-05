package route53_test

import (
	"testing"
	"time"

	"github.com/caddy-dns/route53"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	libdns_route53 "github.com/libdns/route53"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	tests := []struct {
		name      string
		caddyfile string
		expected  libdns_route53.Provider
		wantErr   bool
	}{
		{
			name: "full configuration",
			caddyfile: `route53 {
				region us-east-1
				profile production
				access_key_id AKIAFAKEAWSACCESSKEY
				secret_access_key FAKE/SECRET/KEY/FOR/TESTING/PURPOSES/ONLY
				session_token FAKESESSIONTOKEN123456789ABCDEFGHIJKLMNOP
				max_retries 10
				max_wait_dur 120
				wait_for_propagation true
				hosted_zone_id Z3M3LMPEXAMPLE
			}`,
			expected: libdns_route53.Provider{
				Region:             "us-east-1",
				Profile:            "production",
				AccessKeyId:        "AKIAFAKEAWSACCESSKEY",
				SecretAccessKey:    "FAKE/SECRET/KEY/FOR/TESTING/PURPOSES/ONLY",
				SessionToken:       "FAKESESSIONTOKEN123456789ABCDEFGHIJKLMNOP",
				MaxRetries:         10,
				MaxWaitDuration:    120 * time.Second,
				WaitForPropagation: true,
				HostedZoneID:       "Z3M3LMPEXAMPLE",
			},
		},
		{
			name: "typical static credentials config",
			caddyfile: `route53 {
				region us-west-2
				access_key_id AKIATESTINGONLY12345
				secret_access_key THIS/IS/A/FAKE/SECRET/KEY/DO/NOT/USE
			}`,
			expected: libdns_route53.Provider{
				Region:             "us-west-2",
				AccessKeyId:        "AKIATESTINGONLY12345",
				SecretAccessKey:    "THIS/IS/A/FAKE/SECRET/KEY/DO/NOT/USE",
				WaitForPropagation: true,
			},
		},
		{
			name: "profile-based config",
			caddyfile: `route53 {
				profile staging
				region eu-west-1
			}`,
			expected: libdns_route53.Provider{
				Profile:            "staging",
				Region:             "eu-west-1",
				WaitForPropagation: true,
			},
		},
		{
			name: "aws_profile alias works",
			caddyfile: `route53 {
				aws_profile production
			}`,
			expected: libdns_route53.Provider{
				Profile:            "production",
				Region:             "us-east-1",
				WaitForPropagation: true,
			},
		},
		{
			name: "token alias works",
			caddyfile: `route53 {
				token FAKE_TEST_SESSION_TOKEN_NOT_REAL
			}`,
			expected: libdns_route53.Provider{
				SessionToken:       "FAKE_TEST_SESSION_TOKEN_NOT_REAL",
				Region:             "us-east-1",
				WaitForPropagation: true,
			},
		},
		{
			name: "error on unknown directive",
			caddyfile: `route53 {
				unknown_directive value
			}`,
			wantErr: true,
		},
		{
			name: "error on extra args",
			caddyfile: `route53 unexpected {
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := caddyfile.NewTestDispenser(tt.caddyfile)
			p := &route53.Provider{Provider: &libdns_route53.Provider{}}

			err := p.UnmarshalCaddyfile(d)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check all fields match expected
			if p.Provider.Region != tt.expected.Region {
				t.Errorf("Region = %q, want %q", p.Provider.Region, tt.expected.Region)
			}
			if p.Provider.Profile != tt.expected.Profile {
				t.Errorf("Profile = %q, want %q", p.Provider.Profile, tt.expected.Profile)
			}
			if p.Provider.AccessKeyId != tt.expected.AccessKeyId {
				t.Errorf("AccessKeyId = %q, want %q", p.Provider.AccessKeyId, tt.expected.AccessKeyId)
			}
			if p.Provider.SecretAccessKey != tt.expected.SecretAccessKey {
				t.Errorf("SecretAccessKey = %q, want %q", p.Provider.SecretAccessKey, tt.expected.SecretAccessKey)
			}
			if p.Provider.SessionToken != tt.expected.SessionToken {
				t.Errorf("SessionToken = %q, want %q", p.Provider.SessionToken, tt.expected.SessionToken)
			}
			if p.Provider.MaxRetries != tt.expected.MaxRetries {
				t.Errorf("MaxRetries = %d, want %d", p.Provider.MaxRetries, tt.expected.MaxRetries)
			}
			if p.Provider.MaxWaitDuration != tt.expected.MaxWaitDuration {
				t.Errorf("MaxWaitDuration = %d, want %d", p.Provider.MaxWaitDuration, tt.expected.MaxWaitDuration)
			}
			if p.Provider.WaitForPropagation != tt.expected.WaitForPropagation {
				t.Errorf("WaitForPropagation = %v, want %v",
					p.Provider.WaitForPropagation,
					tt.expected.WaitForPropagation)
			}
			if p.Provider.HostedZoneID != tt.expected.HostedZoneID {
				t.Errorf("HostedZoneID = %q, want %q", p.Provider.HostedZoneID, tt.expected.HostedZoneID)
			}
		})
	}
}
