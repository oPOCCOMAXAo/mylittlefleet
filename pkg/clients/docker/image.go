package docker

import "strings"

// ParseImage parses image tag into image name and tag.
//
//nolint:mnd
func ParseImage(image string) (string, string) {
	parts := strings.SplitN(image, ":", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}
