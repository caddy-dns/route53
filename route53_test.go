package route53_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/caddy-dns/route53"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
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
				Route53MaxWait:     120 * time.Second,
				WaitForRoute53Sync: true,
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
				WaitForRoute53Sync: true,
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
				WaitForRoute53Sync: true,
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
				WaitForRoute53Sync: true,
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
				WaitForRoute53Sync: true,
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
			if p.Provider.Route53MaxWait != tt.expected.Route53MaxWait {
				t.Errorf("Route53MaxWait = %d, want %d", p.Provider.Route53MaxWait, tt.expected.Route53MaxWait)
			}
			if p.Provider.WaitForRoute53Sync != tt.expected.WaitForRoute53Sync {
				t.Errorf("WaitForRoute53Sync = %v, want %v",
					p.Provider.WaitForRoute53Sync,
					tt.expected.WaitForRoute53Sync)
			}
			if p.Provider.HostedZoneID != tt.expected.HostedZoneID {
				t.Errorf("HostedZoneID = %q, want %q", p.Provider.HostedZoneID, tt.expected.HostedZoneID)
			}
		})
	}
}

// TestAppendRecordsUsesUpsert verifies that AppendRecords sends an UPSERT
// (not CREATE) to Route53, preventing "already exists" errors (issue #67).
func TestAppendRecordsUsesUpsert(t *testing.T) {
	var capturedBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		capturedBody = string(body)

		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, `<?xml version="1.0"?>
<ChangeResourceRecordSetsResponse xmlns="https://route53.amazonaws.com/doc/2013-04-01/">
  <ChangeInfo>
    <Id>/change/C1PA6795UKMFR9</Id>
    <Status>INSYNC</Status>
    <SubmittedAt>2024-01-01T00:00:00.000Z</SubmittedAt>
  </ChangeInfo>
</ChangeResourceRecordSetsResponse>`)
	}))
	defer server.Close()

	t.Setenv("AWS_ENDPOINT_URL", server.URL)

	p := &route53.Provider{Provider: &libdns_route53.Provider{
		HostedZoneID:    "ZFAKEZONE",
		AccessKeyId:     "AKIAFAKEKEY123456789",
		SecretAccessKey: "fakesecretkey1234567890abcdefghijklmnop",
		Region:          "us-east-1",
	}}

	record, err := libdns.RR{
		Name: "_acme-challenge",
		Type: "TXT",
		Data: "test-validation-token",
		TTL:  300 * time.Second,
	}.Parse()
	if err != nil {
		t.Fatalf("failed to parse record: %v", err)
	}

	_, err = p.AppendRecords(context.Background(), "example.com.", []libdns.Record{record})
	if err != nil {
		t.Fatalf("AppendRecords() error = %v", err)
	}

	if !strings.Contains(capturedBody, "UPSERT") {
		t.Errorf("expected UPSERT action in request body, got:\n%s", capturedBody)
	}
	if strings.Contains(capturedBody, "CREATE") {
		t.Errorf("expected no CREATE action in request body, got:\n%s", capturedBody)
	}
}
