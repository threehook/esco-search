package postgres

import (
	"context"
	"fmt"

	"github.com/threehook/esco-search/csv"
	"github.com/threehook/esco-search/esco"
	"github.com/threehook/esco-search/internal/repository/ragembedding"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

const (
	promptCatagoryStore = "promptCategoryStore"
	occupationStore     = "occupationStore"
)

// Config holds the configuration necessary for establishing a postgres connection.
type Config struct {
	DefaultDBURL       string
	VectorDBURL        string
	LLModel            string
	RAGEmbeddingConfig ragembedding.Config
}

type RAGData struct {
	config              *Config
	llm                 *ollama.LLM
	embedder            embeddings.Embedder
	PromptCategoryStore *pgvector.Store
	OccupationStore     *pgvector.Store
}

func NewRAGData(config *Config) (*RAGData, error) {
	llm, err := ollama.New(ollama.WithModel(config.LLModel))
	if err != nil {
		return nil, fmt.Errorf("creating llm model: %s", err)
	}

	embedder, err := createEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("Error creating llm embedder: %s", err)
	}

	storeMap, err := createStores(context.Background(), config, embedder)
	if err != nil {
		return nil, fmt.Errorf("creating RAG store: %s", err)
	}

	return &RAGData{
		config:              config,
		llm:                 llm,
		embedder:            embedder,
		PromptCategoryStore: storeMap[promptCatagoryStore],
		OccupationStore:     storeMap[occupationStore],
	}, nil
}

func createEmbedder(llm *ollama.LLM) (embeddings.Embedder, error) {
	return embeddings.NewEmbedder(llm)
}

func (r *RAGData) CreateEmbeddingContent() error {
	skillsByOccupationMap, err := esco.CreateSkillsByOccupationMap("./data/occupationSkillRelations_nl.csv", "./data/skills_nl.csv")
	if err != nil {
		return fmt.Errorf("failed creating skills in document store: %s", err)
	}

	columns := []string{"conceptUri", "code", "iscoGroup", "preferredLabel", "altLabels", "description"}
	occupations, err := csv.ReadCSV("./data/occupations_nl.csv", columns...)
	if err != nil {
		return fmt.Errorf("failed creating occupations in document store: %s", err)
	}

	// Add skills to occupations
	esco.AddSkillsToOccupations(occupations, skillsByOccupationMap)

	educationTypesLevelsByOccupationMap, err := esco.CreateEducationTypesLevelsByOccupationMap("./data/educationTypesLevels_nl.csv")
	if err != nil {
		return fmt.Errorf("failed creating education types and levels in document store: %s", err)
	}

	// Add education types and levels to occupations
	esco.AddEducationTypesAndLevelsToOccupations(occupations, educationTypesLevelsByOccupationMap)

	docs := r.createDocuments(occupations)
	return r.storeDocuments(docs, r.OccupationStore)
}

// Create Postgres pgvector stores for prompt categories and occupations.
func createStores(ctx context.Context, config *Config, embedder embeddings.Embedder) (map[string]*pgvector.Store, error) {
	storeMap := make(map[string]*pgvector.Store, 2)

	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(config.VectorDBURL),
		pgvector.WithEmbedder(embedder),
		pgvector.WithCollectionName("promptcategories"),
		pgvector.WithCollectionTableName("occupation_collection"),
		pgvector.WithEmbeddingTableName("occupation_embedding"),
		pgvector.WithHNSWIndex(16, 64, "vector_l2_ops"),
		pgvector.WithVectorDimensions(1024),
	)
	if err != nil {
		return nil, err
	}
	storeMap[promptCatagoryStore] = &store

	store, err = pgvector.New(
		ctx,
		pgvector.WithConnectionURL(config.VectorDBURL),
		pgvector.WithEmbedder(embedder),
		pgvector.WithCollectionName("occupations"),
		pgvector.WithCollectionTableName("occupation_collection"),
		pgvector.WithEmbeddingTableName("occupation_embedding"),
		pgvector.WithHNSWIndex(16, 64, "vector_l2_ops"),
		pgvector.WithVectorDimensions(1024),
	)
	if err != nil {
		return nil, err
	}
	storeMap[occupationStore] = &store

	return storeMap, nil
}

func (r *RAGData) createDocuments(occupations [][]string) []schema.Document {
	data := make([]schema.Document, len(occupations))
	for i, occupation := range occupations {
		// Add document for the description, skills and competences
		description := fmt.Sprintf("De beschrijving van beroep '%s' is: %s\n", occupation[2], occupation[4])
		skillsCompetences := fmt.Sprintf("De benodigde vaardigheden en competenties van beroep '%s' zijn: %s\n", occupation[2], occupation[0])
		educationTypesLevels := occupation[6]
		pageContent := description + skillsCompetences + educationTypesLevels
		metadata := map[string]any{
			"code":                      occupation[5],
			"beroep":                    occupation[2],
			"groep":                     occupation[1],
			"alternatieve beroepsnamen": occupation[3],
		}
		data[i] = schema.Document{
			PageContent: pageContent,
			Metadata:    metadata,
		}
	}

	return data
}

func (r *RAGData) storeDocuments(docs []schema.Document, store *pgvector.Store) error {
	_, err := store.AddDocuments(context.Background(), docs)
	if err != nil {
		return err
	}

	return nil
}

func (r *RAGData) findDocuments(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	return r.OccupationStore.SimilaritySearch(ctx, query, numDocuments, options...)
}
