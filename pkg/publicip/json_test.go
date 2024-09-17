package publicip_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/nicus101/godyndns-ovh/pkg/publicip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpJsonIper_invalidApiUrl(t *testing.T) {
	_, err := publicip.NewHttpJsonIper(
		"http://example.org",
		"not-relevant",
	).Ip(context.Context(nil))
	assert.Error(t, err)
}

func TestHttpJsonIper(t *testing.T) {
	tests := []struct {
		name string

		key string

		given  string
		sleep  time.Duration
		status int

		expectErr string
		expectIP  netip.Addr
	}{{
		name:     "normal",
		key:      "query",
		given:    `{"query":"127.12.1.0"}`,
		expectIP: netip.MustParseAddr("127.12.1.0"),
	}, {
		name:      "timeout",
		key:       "query",
		sleep:     time.Millisecond * 20,
		given:     `{"query":"127.12.1.0"}`,
		expectErr: "context deadline exceeded",
	}, {
		name:      "malformed-response",
		given:     `=UwU=`,
		expectErr: "invalid character",
	}, {
		name:      "status-500",
		given:     `{"query":"127.12.1.0"}`,
		status:    http.StatusInternalServerError,
		expectErr: "received status code",
	}, {
		name:      "status-200-but-no-ip",
		key:       "kici-kici",
		given:     `{}`,
		expectErr: "doesn't have valid key",
	}, {
		name:      "status-200-but-malformaed-ip",
		key:       "ip",
		given:     `{"ip":"sruwu"}`,
		expectErr: "unable to parse IP",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if test.sleep != 0 {
					time.Sleep(test.sleep)
				}
				if test.status != 0 {
					w.WriteHeader(test.status)
				}
				fmt.Fprint(w, test.given)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*10)
			defer cancel()

			ip, err := publicip.NewHttpJsonIper(
				server.URL,
				test.key,
			).Ip(ctx)
			if test.expectErr != "" {
				assert.ErrorContains(t, err, test.expectErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectIP, ip)
		})
	}
}
