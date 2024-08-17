package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

var ipAPiUrl = "http://ip-api.com/json/"

func GetIP(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipAPiUrl, nil)
	if err != nil {
		return net.IP{}, fmt.Errorf("cannot get ip from %s errortype: %w", ipAPiUrl, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return net.IP{}, fmt.Errorf("cannot perform body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return net.IP{}, fmt.Errorf("error received statuscode: %d", resp.StatusCode)
	}

	response := struct {
		Query net.IP `json:"query"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return net.IP{}, fmt.Errorf(
			"cannot unmarshal ip: %w", err,
		)
	}

	if response.Query == nil {
		return net.IP{}, fmt.Errorf("error malformed ip adress received")
	}

	return response.Query, nil
}
