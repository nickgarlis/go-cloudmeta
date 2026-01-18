package ipdetect

import (
	"net"
	"net/http"
	"time"
)

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: 300 * time.Millisecond,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 300 * time.Millisecond,
			}).DialContext,
		},
	}
}
