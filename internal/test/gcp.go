package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// CreateMockGCPServer creates a minimal mock server for GCP metadata service
func CreateMockGCPServer(disabled ...bool) *httptest.Server {
	isDisabled := len(disabled) > 0 && disabled[0]

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for required Metadata-Flavor header
		if r.Header.Get("Metadata-Flavor") != "Google" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Missing Metadata-Flavor header"))
			return
		}

		if isDisabled {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Metadata service temporarily unavailable"))
			return
		}

		switch r.URL.Path {
		case "/computeMetadata/v1/project/project-id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("my-test-project"))

		case "/computeMetadata/v1/instance/id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1234567890123456789"))

		case "/computeMetadata/v1/instance/network-interfaces/0/ip":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("10.128.0.5"))

		case "/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("34.123.45.67"))

		case "/computeMetadata/v1/instance/network-interfaces/0/ipv6s":
			w.WriteHeader(http.StatusOK)
			ips := []string{"2001:db8:85a3::8a2e:370:7334", "2001:db8:85a3::8a2e:370:7335", "2001:db8:85a3::8a2e:370:7336"}
			w.Write([]byte(strings.Join(ips, "\n")))
		case "/computeMetadata/v1/instance/hostname":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test-instance-1.c.my-test-project.internal"))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))
}
