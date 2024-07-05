package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

func loadDocs(source string) (int, error) {

	fmt.Println("loading data from", source)
	var documentCount int
	store, err := getVectorStore()

	if err != nil {
		return documentCount, err
	}
	docs, err := getDocs(source)

	if err != nil {
		return documentCount, err
	}

	fmt.Println("no. of documents to be loaded", len(docs))

	if len(docs) > 0 && store != nil {
		_, err = store.AddDocuments(context.Background(), docs)
	}

	if err != nil {
		return documentCount, err
	}

	fmt.Println("data successfully loaded into vector store")
	documentCount = len(docs)
	return documentCount, nil
}

func getDocs(source string) ([]schema.Document, error) {
	var urlParser = &url.URL{}
	var err error
	urlParser, err = url.Parse(source)
	if err != nil {
		return nil, err
	}
	if urlParser == nil {
		urlParser = &url.URL{}
		urlParser.Scheme = "file"
	}
	if urlParser.Scheme == "" {
		urlParser.Scheme = "file"
	}
	switch urlParser.Scheme {
	case "file":
		var f *os.File
		f, err = os.Open(source)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b := bufio.NewReader(f)
		switch filepath.Ext(source) {
		case ".html":
			fmt.Println("Processing html")
			return processHTML(b)
		case ".txt":
			fmt.Println("Processing Text")
			return processText(b)
		case ".csv":
			fmt.Println("Processing csv")
			return processCSV(b)
		case ".pdf":
			fmt.Println("Processing pdf")
			return processPDF(b)
		default:
			return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(source))
		}
	case "http", "https":
		var resp *http.Response
		resp, err = http.Get(source)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")

		switch {
		case strings.Contains(contentType, "text/html"):
			log.Println("Processing HTML")
			return processHTML(resp.Body)
		case strings.Contains(contentType, "text/plain"):
			log.Println("Processing Text")
			return processText(resp.Body)
		case strings.Contains(contentType, "text/csv"):
			log.Println("Processing CSV")
			return processCSV(resp.Body)
		case strings.Contains(contentType, "application/pdf"):
			log.Println("Processing PDF")
			return processPDF(resp.Body)
		default:
			return nil, fmt.Errorf("unsupported content type: %s", contentType)
		}
	default:
		return nil, fmt.Errorf("unsupported scheme: %s for url %v", urlParser.Scheme, urlParser)
	}

}

func processHTML(body io.Reader) ([]schema.Document, error) {
	docs, err := documentloaders.NewHTML(body).LoadAndSplit(context.Background(), textsplitter.NewRecursiveCharacter())
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func processText(body io.Reader) ([]schema.Document, error) {
	docs, err := documentloaders.NewText(body).LoadAndSplit(context.Background(), textsplitter.NewRecursiveCharacter())
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func processCSV(body io.Reader) ([]schema.Document, error) {
	docs, err := documentloaders.NewCSV(body).LoadAndSplit(context.Background(), textsplitter.NewRecursiveCharacter())
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func processPDF(body io.Reader) ([]schema.Document, error) {
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, body)
	if err != nil {
		fmt.Println("error copying pdf body to buffer, err:", err)
		return nil, err
	}
	reader := bytes.NewReader(buff.Bytes())
	docs, err := documentloaders.NewPDF(reader, size).LoadAndSplit(context.Background(), textsplitter.NewRecursiveCharacter())
	if err != nil {
		fmt.Println("error loading pdf, err:", err)
		return nil, err
	}
	return docs, nil
}
