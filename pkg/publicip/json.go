package publicip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
)

type HttpJsonIper struct {
	url string
	key string
}

var _ Iper = &HttpJsonIper{}

func NewHttpJsonIper(url, key string) *HttpJsonIper {
	return &HttpJsonIper{
		url: url,
		key: key,
	}
}

func (iper *HttpJsonIper) Ip(ctx context.Context) (netip.Addr, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		iper.url,
		nil,
	)
	if err != nil {
		return netip.Addr{}, fmt.Errorf(
			"cannot get ip from %s errortype: %w",
			iper.url, err,
		)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return netip.Addr{}, fmt.Errorf(
			"cannot perform body: %w",
			err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return netip.Addr{}, fmt.Errorf(
			"error received status code: %d",
			resp.StatusCode,
		)
	}

	response := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return netip.Addr{}, fmt.Errorf(
			"cannot unmarshal response: %w",
			err,
		)
	}

	valueStr, ok := response[iper.key].(string)
	if !ok {
		return netip.Addr{}, fmt.Errorf(
			"response %v doesn't have valid key",
			response,
		)
	}

	valueIp, err := netip.ParseAddr(valueStr)
	if err != nil {
		return netip.Addr{}, fmt.Errorf(
			"error malformed ip adress received: %w",
			err,
		)
	}

	return valueIp, err
}
