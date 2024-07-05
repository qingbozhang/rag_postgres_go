package main

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
)

type ollamaEmbedder struct {
	llm *ollama.LLM
}

func NewOllamaEmbedder(llm *ollama.LLM) *ollamaEmbedder {
	return &ollamaEmbedder{llm: llm}
}

func (o *ollamaEmbedder) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	return o.llm.CreateEmbedding(ctx, texts)
}
