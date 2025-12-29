package test

import (
	"net/http"
	"net/http/httptest"
)

// CreateMockAWSServer creates a shared mock server for AWS metadata service
func CreateMockAWSServer(disabled ...bool) *httptest.Server {
	isDisabled := len(disabled) > 0 && disabled[0]
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDisabled {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("IMDSv2 disabled or blocked by security policy"))
			return
		}

		token := "AQAAANhJbmV0YW1ldGFkYXRhLmFtYXpvbmF3cy5jb20vMjAyMi0xMi0yMQ=="
		if r.URL.Path == "/latest/api/token" {
			if r.Method != "PUT" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			if r.Header.Get("X-aws-ec2-metadata-token-ttl-seconds") == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing TTL header"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(token))
			return
		}

		if r.Header.Get("X-aws-ec2-metadata-token") != token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid or missing IMDSv2 token"))
			return
		}

		switch r.URL.Path {
		case "/latest/meta-data/instance-id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("i-1234567890abcdef0"))

		case "/latest/meta-data/ami-id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ami-0abcdef1234567890"))

		case "/latest/meta-data/local-ipv4":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("10.0.1.100"))

		case "/latest/meta-data/public-ipv4":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("54.123.45.67"))

		case "/latest/meta-data/ipv6":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))

		case "/latest/meta-data/hostname":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ip-10-0-1-100.us-west-2.compute.internal"))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))
}
