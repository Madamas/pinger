package pinger

import (
	"net/http"
	"pinger/packages/config"
)

func NewClient(
	config config.Config,
) *http.Client {
	client := http.Client{}

	client.Timeout = config.Client.TimeoutDuration

	return &client
}
