package main

import (
	"net/http"
	"os"
	"time"
)

var client *MonzoClient

type MonzoClient struct {
	id           string
	secret       string
	accessToken  string
	refreshToken string
	callbackCode string
	code         string
	httpClient   *http.Client
}

func init() {
	if client != nil {
		return
	}

	client = &MonzoClient{
		id:           os.Getenv("MONZO_CLIENT_ID"),
		secret:       os.Getenv("MONZO_CLIENT_SECRET"),
		accessToken:  os.Getenv("MONZO_ACCESS_TOKEN"),
		refreshToken: os.Getenv("MONZO_REFRESH_TOKEN"),
		callbackCode: "",
		code:         "",
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewClient() *MonzoClient {
	return client
}

func (c *MonzoClient) Do(req *http.Request) (*http.Response, error) {
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
