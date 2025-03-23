package models

// ServerConfig represents public configuration for the server.
type ServerConfig struct {
	ReverseProxyEnabled bool
	NginxStatus         ContainerStatus
}
