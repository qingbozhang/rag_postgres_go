package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"log"
	"net/url"
	"os"
	"sync"
)

var Store vectorstores.VectorStore
var storeMutex = &sync.Mutex{}
var llmMutex = &sync.Mutex{}
var llm *ollama.LLM

func init() {
	var err error
	Store, err = getVectorStore()
	if err != nil {
		log.Fatal(err)
	}
	llm, err = ollama.New(ollama.WithModel(`llama3`), ollama.WithServerURL("http://localhost:11434"))
	if err != nil {
		log.Println("error connecting to ollama %v", err)
	}
}

func getVectorStore() (vectorstores.VectorStore, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	if Store == nil {
		host := os.Getenv("PG_HOST")
		if host == "" {
			log.Fatal("missing PG_HOST")
		}

		user := os.Getenv("PG_USER")
		if user == "" {
			log.Fatal("missing PG_USER")
		}

		password := os.Getenv("PG_PASSWORD")
		if password == "" {
			log.Fatal("missing PG_PASSWORD")
		}

		dbName := os.Getenv("PG_DB")
		if dbName == "" {
			log.Fatal("missing PG_DB")
		}
		var err error
		connURLFormat := "postgres://%s:%s@%s:5432/%s?sslmode=disable"

		pgConnURL := fmt.Sprintf(connURLFormat, user, url.QueryEscape(password), host, dbName)
		llm, err = getOllama()
		if err != nil {
			return nil, err
		}

		e, err := embeddings.NewEmbedder(NewOllamaEmbedder(llm))
		if err != nil {
			return nil, err
		}
		Store, err = pgvector.New(
			context.Background(),
			//pgvector.WithPreDeleteCollection(true),
			pgvector.WithConnectionURL(pgConnURL),
			pgvector.WithEmbedder(e),
		)
		if err != nil {
			log.Fatalf("error creating vector store: %v", err)
			return nil, err
		}

		fmt.Println("vector store ready")

		return Store, nil
	} else {
		return Store, nil
	}
}

func getOllama() (*ollama.LLM, error) {
	llmMutex.Lock()
	defer llmMutex.Unlock()
	if llm != nil {
		return llm, nil
	}
	var err error
	llm, err = ollama.New(ollama.WithModel(`llama3`), ollama.WithServerURL("http://localhost:11434"))
	if err != nil {
		return nil, err
	}
	return llm, nil

}
