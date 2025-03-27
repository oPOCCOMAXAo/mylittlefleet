package docker

import "strings"

// ParseEnv parses environment variable string into name and value.
//
//nolint:mnd
func ParseEnv(env string) (string, string) {
	parts := strings.SplitN(env, "=", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}
