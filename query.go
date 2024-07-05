package main

import (
	"context"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/vectorstores"
)

type RagSearchResponse struct {
	FinalAnswer string `json:"final_answer"`
}

type SemanticSearchResult struct {
	PageContent string  `json:"page_content"`
	Score       float64 `json:"score"`
}

type SemanticSearchResponse struct {
	Results []SemanticSearchResult `json:"results"`
}
type LoadRequest struct {
	Source string `json:"source"`
}

type SearchRequest struct {
	Query      string `json:"query"`
	MaxResults int    `json:"maxResults"`
}

func ragSearch(question string, numOfResults int) (RagSearchResponse, error) {
	store, err := getVectorStore()
	if err != nil {
		return RagSearchResponse{}, err
	}

	llm, err = getOllama()
	if err != nil {
		return RagSearchResponse{}, err
	}

	result, err := chains.Run(
		context.Background(),
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, numOfResults),
		),
		question,
		chains.WithMaxTokens(5120),
	)
	if err != nil {
		return RagSearchResponse{}, err
	}

	return RagSearchResponse{FinalAnswer: result}, nil
}

func semanticSearch(searchQuery string, maxResults int) (SemanticSearchResponse, error) {
	store, err := getVectorStore()
	if err != nil {
		return SemanticSearchResponse{}, err
	}

	searchResults, err := store.SimilaritySearch(context.Background(), searchQuery, maxResults)
	if err != nil {
		return SemanticSearchResponse{}, err
	}

	var results []SemanticSearchResult
	for _, doc := range searchResults {
		results = append(results, SemanticSearchResult{
			PageContent: doc.PageContent,
			Score:       float64(doc.Score),
		})
	}

	return SemanticSearchResponse{Results: results}, nil
}
