package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	UploadDirectory = "./uploads"
)

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file type
	if !strings.HasSuffix(header.Filename, ".csv") && !strings.HasSuffix(header.Filename, ".txt") && !strings.HasSuffix(header.Filename, ".pdf") && !strings.HasSuffix(header.Filename, ".html") {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// Save the file
	err = CreateDirectoryIfNotExist("./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(UploadDirectory, header.Filename)
	var out *os.File
	out, err = os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = loadDocs(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("File Uploaded.Data successfully processed."))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateDirectoryIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
