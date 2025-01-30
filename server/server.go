package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Server response type
type UpdateResponse struct {
	UpdateVersion string `json:"updateVersion"`
	DownloadLink  string `json:"downloadLink"`
	Checksum      string `json:"checkSum"`
}

// Server constants that populate the response
var (
	latestVersion string = "1.1.2"
	binaryName    string = "myapp-update"
	port          string = ":8080"
	baseUrl       string = "http://localhost:8080"
)

func main() {
	http.HandleFunc("/checkUpdate", updateHandler)
	http.HandleFunc("/download", downloadHandler)

	log.Printf("Server starting on %v", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server to start server on port %v: %v", port, err)
	}

}

func calculateCheckSum(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for checksum: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil

}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Calculate checksum for update binary
	binaryPath := filepath.Join("/app/server/binaries", binaryName)
	checksum, err := calculateCheckSum(binaryPath)
	if err != nil {
		log.Printf("Err calculating checksum: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// log.Printf("updateHandler triggered")
	resp := UpdateResponse{
		UpdateVersion: latestVersion,
		DownloadLink:  fmt.Sprintf("%s/download?file=%s", baseUrl, binaryName),
		Checksum:      checksum,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	binaryPath := filepath.Join("/app/server/binaries", binaryName)
	w.Header().Set("Content-Type", "application/json")

	// Send the compiled update at binaryPath
	http.ServeFile(w, r, binaryPath)
}
