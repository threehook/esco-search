package contract

import (
	"context"

	"github.com/threehook/esco-search/model"
)

// Repository defines the repository interface which is used to interact with a storage layer for matches.
type Repository interface {
	// ReadBeroepenFromRAGStore does a similarity search against the vector store.
	ReadBeroepenFromRAGStore(ctx context.Context, searchRequest string) ([]model.BeroepenMatch, error)
}
