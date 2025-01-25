package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/load", loadHandler)
	http.HandleFunc("/rag_search", ragSearchHandler)
	http.HandleFunc("/chat", ragSearchHandler)
	http.HandleFunc("/semantic_search", semanticSearchHandler)
	http.HandleFunc("/upload", fileUploadHandler)
	http.HandleFunc("/knowledge-base/upload", fileUploadHandler)
	http.HandleFunc("/index.html", serveHTMLHandler)
	http.HandleFunc("/", serveHTMLHandler)

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
