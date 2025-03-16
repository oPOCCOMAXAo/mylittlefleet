package netutils

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

// GetCurrentPublicIP returns the current public IP address of the client.
// Uses the ipquery.io API to get the IP address.
// Returns an error if the request fails.
// Returns nil if received IP address is invalid.
func GetCurrentPublicIP(
	ctx context.Context,
) (net.IP, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.ipquery.io",
		http.NoBody,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ip := net.ParseIP(string(body))

	return ip, nil
}

func IPHostOrEmpty(ip net.IP) string {
	if ip == nil {
		return ""
	}

	return ip.String()
}
