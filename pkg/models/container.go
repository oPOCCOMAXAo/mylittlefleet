package models

type ContainerStatus int

const (
	CSError ContainerStatus = iota
	CSStarting
	CSRunning
	CSStopping
	CSStopped
)

func (s ContainerStatus) String() string {
	switch s {
	case CSStarting:
		return "Starting"
	case CSRunning:
		return "Running"
	case CSStopping:
		return "Stopping"
	case CSStopped:
		return "Stopped"
	case CSError:
		return "Error"
	default:
		return "Unknown"
	}
}

func (s ContainerStatus) BSSubClass() string {
	switch s {
	case CSStarting:
		return "success"
	case CSRunning:
		return "success"
	case CSStopping:
		return "warning"
	case CSStopped:
		return "danger"
	case CSError:
		return "danger"
	default:
		return "secondary"
	}
}
