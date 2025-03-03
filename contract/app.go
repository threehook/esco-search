// Package contract defines the contract that an app client can use service
package contract

import (
	"context"

	"github.com/threehook/esco-search/model"
)

const (
	ErrorKey string = "error"
)

// App provides the interface to communicate with external components.
type App interface {
	// SearchBeroepenViaPrompt searches for beroepen via a free format prompt and using AI RAG.
	SearchBeroepenViaPrompt(ctx context.Context, searchRequest string) ([]model.BeroepenMatch, error)
}
