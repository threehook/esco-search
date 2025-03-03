package postgres

import (
	"context"
	"log"
)

func InitDatabase(config *Config) {
	ctx := context.Background()

	// Create default database connection
	conn, err := NewConnection(ctx, config, true)
	if err != nil {
		log.Fatalf("creating default database connection: %s", err)
	}

	created, err := conn.CreateDB(ctx)
	if err != nil {
		log.Fatalf("creating database: %s", err)
	}

	// Close default database connection
	conn.CloseConnection(ctx)

	// Create vector database connection
	conn, err = NewConnection(ctx, config, false)
	if err != nil {
		log.Fatalf("creating vector database connection: %s", err)
	}

	err = conn.CreateDBExtensions(ctx, []string{"vector"})
	if err != nil {
		log.Fatalf("creating database extension", err)
	}

	if created {
		// Create a RAG store
		store, err := NewRAGData(config)
		if err != nil {
			log.Fatalf("creating RAG store: %s", err.Error())
		}
		err = store.CreateEmbeddingContent()
		if err != nil {
			log.Fatalf("creating database occupation content: %s", err.Error())
		}
	}
	// Close vector database connection
	conn.CloseConnection(ctx)
}
