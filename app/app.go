// Package app implements the App contract
package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/threehook/esco-search/contract"
	"github.com/threehook/esco-search/model"
)

// App implements the app.App interface.
type App struct {
	repository contract.Repository
}

// NewApp returns an App implementation.
func NewApp(repository contract.Repository) *App {
	return &App{
		repository: repository,
	}
}

// SearchBeroepenViaPrompt searches for beroepen via a free format prompt and using AI RAG.
func (a *App) SearchBeroepenViaPrompt(ctx context.Context, searchRequest string) ([]model.BeroepenMatch, error) {
	beroepen, err := a.repository.ReadBeroepenFromRAGStore(ctx, searchRequest)
	if err != nil {
		return nil, errors.Wrap(err, "getting beroepen from repository")
	}

	return beroepen, nil
}
