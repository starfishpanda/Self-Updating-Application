package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Server response type
type UpdateResponse struct {
	UpdateVersion string `json:"updateVersion"`
	DownloadLink  string `json:"downloadLink"`
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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// log.Printf("updateHandler triggered")
	resp := UpdateResponse{
		UpdateVersion: latestVersion,
		DownloadLink:  fmt.Sprintf("%s/download?file=%s", baseUrl, binaryName),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")

	exePath, err := os.Executable()
	log.Printf("Exepath: %v", exePath)
	if err != nil {
		log.Printf("Error getting executable path: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	exeDir := filepath.Dir(exePath)
	binariesDir := filepath.Join(exeDir, "binaries")
	// Sanitize file name from client request
	sanitizeFileName := filepath.Clean(fileName)
	binaryPath := filepath.Join(binariesDir, sanitizeFileName)

	log.Printf("Requested file: %s", sanitizeFileName)
	log.Printf("Serving file from path: %s", binaryPath)

	// Check if file exists
	fileInfo, err := os.Stat(binaryPath)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}

	// Create content headers for client download
	w.Header().Set("Content-Disposition", fmt.Sprintf("application; filename=\"%s\"", binaryPath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Send the compiled update at binaryPath
	http.ServeFile(w, r, binaryPath)
}
