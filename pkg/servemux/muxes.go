// Package servemux keeps track of ports and their corresponding http.ServeMux instances and starts
// listening and serving them.
package servemux

import (
	"fmt"
	"net/http"
	"time"
)

// Muxes keeps track of ports and their corresponding http.ServeMux instances.
type Muxes map[int]*http.ServeMux

// ForPort returns the serve mux for the specified port, creating a new one if it doesn't exist yet.
func (m Muxes) ForPort(port int) *http.ServeMux {
	mux, ok := m[port]
	if !ok {
		mux = http.NewServeMux()
		m[port] = mux
	}

	return mux
}

// ListenAndServe starts listening on the ports, serving from the muxes.  Never returns unless
// there is an error.
func (m Muxes) ListenAndServe() error {
	errChan := make(chan error, 1)
	for port, mux := range m {
		go listenAndServe(port, mux, errChan)
	}

	return <-errChan
}

func listenAndServe(port int, mux http.Handler, errChan chan<- error) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd
		Handler:           mux,
	}

	errChan <- server.ListenAndServe()
}
