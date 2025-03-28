package event

type Listener[T comparable] struct {
	message       chan T
	newClients    chan chan T
	closedClients chan chan T
	allClients    map[chan T]bool
}

//nolint:mnd
func NewListener[T comparable]() *Listener[T] {
	return &Listener[T]{
		message:       make(chan T, 10),
		newClients:    make(chan chan T, 10),
		closedClients: make(chan chan T, 10),
		allClients:    make(map[chan T]bool),
	}
}

// Serve starts the listener and waits for new clients, closed clients and messages.
//
// Blocks until the listener is closed.
func (l *Listener[T]) Serve() {
	for {
		select {
		case client := <-l.newClients:
			l.allClients[client] = true

		case client := <-l.closedClients:
			delete(l.allClients, client)
			close(client)

			// Drain the channel to avoid deadlock.
			//nolint:revive
			for range client {
			}

		case message := <-l.message:
			for client := range l.allClients {
				select {
				case client <- message:
				default:
				}
			}
		}
	}
}

func (l *Listener[T]) Notify(message T) {
	l.message <- message
}

// NewClient creates a new client and returns a channel to receive messages.
//
// Should be closed when the client is no longer needed with CloseClient.
//
//nolint:mnd
func (l *Listener[T]) NewClient() chan T {
	client := make(chan T, 100)
	l.newClients <- client

	return client
}

func (l *Listener[T]) CloseClient(client chan T) {
	l.closedClients <- client
}

func (l *Listener[T]) Close() {
	for client := range l.allClients {
		delete(l.allClients, client)
		close(client)
	}

	close(l.message)
	close(l.newClients)
	close(l.closedClients)
}
