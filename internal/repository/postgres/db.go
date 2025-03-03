// Package pgvector implements the repository interface with a postgresql vector database
package postgres

import (
	"context"
	"log"

	"github.com/cockroachdb/errors"
	"github.com/threehook/esco-search/model"
	"github.com/tmc/langchaingo/vectorstores"
)

const (
	codeMetadata   = "code"
	beroepMetadata = "beroep"
)

// Repository database struct used by service api.
type Repository struct {
	db             *RAGData
	maxDocuments   int
	scoreThreshold float32
}

// NewRepository creates a new Repository.
func NewRepository(config *Config) (*Repository, error) {
	store, err := NewRAGData(config)
	if err != nil {
		log.Fatalf("creating RAG store: %s", err.Error())
	}

	db := &Repository{
		db:             store,
		maxDocuments:   config.RAGEmbeddingConfig.MaxDocuments,
		scoreThreshold: config.RAGEmbeddingConfig.ScoreThreshold,
	}

	return db, nil
}

// ReadBeroepenFromRAGStore does a similarity search against the vector store.
func (r Repository) ReadBeroepenFromRAGStore(ctx context.Context, searchRequest string) ([]model.BeroepenMatch, error) {
	docs, err := r.db.findDocuments(ctx, searchRequest, r.maxDocuments, vectorstores.WithScoreThreshold(r.scoreThreshold))
	if err != nil {
		return nil, errors.Wrap(err, "reading beroepen form RAG store")
	}

	beroepen := make([]model.BeroepenMatch, len(docs))
	for i, doc := range docs {
		code := doc.Metadata[codeMetadata]
		beroep := doc.Metadata[beroepMetadata]
		beroepen[i] = model.BeroepenMatch{
			Code:   code.(string),
			Beroep: beroep.(string),
		}
	}

	return beroepen, nil
}
