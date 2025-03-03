package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

type Connection struct {
	*pgx.Conn
}

// Create a postgres connection
func NewConnection(ctx context.Context, config *Config, defaultDB bool) (*Connection, error) {
	var connString string
	if defaultDB {
		connString = config.DefaultDBURL
	} else {
		connString = config.VectorDBURL
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(err, "creating database connection pool")
	}

	return &Connection{
		conn,
	}, nil
}

func (c *Connection) CreateDB(ctx context.Context) (bool, error) {
	// Database to check
	dbName := "vectorstore" // TODO move to .env
	query := `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);`

	var exists bool
	err := c.QueryRow(ctx, query, dbName).Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "checking database existence")
	}

	if exists {
		slog.Info("No need to create database, it already exists.")
		return false, nil
	}

	// SQL statement to create the database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s;", dbName)
	slog.Info("Creating and filling vector database. Could take more than a minute.")
	_, err = c.Exec(ctx, createDBSQL)
	if err != nil {
		return false, errors.Wrap(err, "filling vector database")
	}

	return true, nil
}

// CreateDBExtensions create database extensions.
func (c *Connection) CreateDBExtensions(ctx context.Context, extensions []string) error {
	for _, extension := range extensions {
		_, err := c.Exec(ctx, fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", extension))
		if err != nil {
			return err
		}
	}

	return nil
}

//func (c *Connection) AlterVectorDimensions(ctx context.Context, dimensions int) error {
//	alterEmbeddingsColumn := fmt.Sprintf("ALTER TABLE langchain_pg_embedding ALTER COLUMN embedding TYPE vector(%d);", dimensions)
//	_, err := c.Exec(ctx, alterEmbeddingsColumn)
//	if err != nil {
//		return errors.Wrap(err, "altering vector dimensions")
//	}
//
//	return nil
//}

func (c *Connection) CreateVectorIndex(ctx context.Context) error {
	createVectorIndex := "CREATE INDEX beroepen_vector_idx ON langchain_pg_embedding USING hnsw (embedding vector_cosine_ops);"
	_, err := c.Exec(ctx, createVectorIndex)
	if err != nil {
		return errors.Wrap(err, "creating vector index")
	}

	return nil
}

func (p *Connection) CloseConnection(ctx context.Context) {
	p.Close(ctx)
}
