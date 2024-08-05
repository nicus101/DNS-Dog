package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type IP struct {
	Query string
}

func GetIP() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}
	var ip IP
	json.Unmarshal(body, &ip)
	return ip.Query
}
