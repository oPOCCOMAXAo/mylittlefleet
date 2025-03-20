package models

// ServerConfig represents public configuration for the server.
type ServerConfig struct {
	ReverseProxyEnabled bool
	NginxStatus         ServerStatus
}

type ServerStatus int

const (
	ServerStatusError ServerStatus = iota
	ServerStatusStarting
	ServerStatusRunning
	ServerStatusStopping
	ServerStatusStopped
)

func (s ServerStatus) String() string {
	switch s {
	case ServerStatusStarting:
		return "Starting"
	case ServerStatusRunning:
		return "Running"
	case ServerStatusStopping:
		return "Stopping"
	case ServerStatusStopped:
		return "Stopped"
	case ServerStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

func (s ServerStatus) BSSubClass() string {
	switch s {
	case ServerStatusStarting:
		return "success"
	case ServerStatusRunning:
		return "success"
	case ServerStatusStopping:
		return "warning"
	case ServerStatusStopped:
		return "danger"
	case ServerStatusError:
		return "danger"
	default:
		return "secondary"
	}
}
