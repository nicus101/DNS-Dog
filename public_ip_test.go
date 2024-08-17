package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIP_invalidApiUrl(t *testing.T) {
	_, err := GetIP(context.Context(nil))
	assert.Error(t, err)
}

func TestGetIP(t *testing.T) {
	tests := []struct {
		name string

		given  string
		sleep  time.Duration
		status int

		expectErr string
		expectIP  net.IP
	}{{
		name:     "normal",
		given:    `{"query":"127.12.1.0"}`,
		expectIP: net.ParseIP("127.12.1.0"),
	}, {
		name:      "timeout",
		sleep:     time.Millisecond * 20,
		given:     `{"query":"127.12.1.0"}`,
		expectErr: "context deadline exceeded",
	}, {
		name:      "malformed-response",
		given:     `=UwU=`,
		expectErr: "cannot unmarshal ip",
	}, {
		name:      "status-500",
		given:     `{"query":"127.12.1.0"}`,
		status:    http.StatusInternalServerError,
		expectErr: "received statuscode",
	}, {
		name:      "status-200-but-no-ip",
		given:     `{}`,
		expectErr: "malformed ip",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			oldIpAPiUrl := ipAPiUrl
			defer func() {
				ipAPiUrl = oldIpAPiUrl
			}()

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
			ipAPiUrl = server.URL

			ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*10)
			defer cancel()

			ip, err := GetIP(ctx)
			if test.expectErr != "" {
				assert.ErrorContains(t, err, test.expectErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectIP, ip)
		})
	}
}
