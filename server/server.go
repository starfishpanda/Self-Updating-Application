package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
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
	Signature     string `json:"signature"`
}

// Server constants that populate the response
var (
	latestVersion string = "1.1.2"
	binaryName    string = "myapp-update"
	port          string = ":8080"
	baseUrl       string = "http://localhost:8080"
	privateKey    *rsa.PrivateKey
)

func main() {
	http.HandleFunc("/checkUpdate", updateHandler)
	http.HandleFunc("/download", downloadHandler)

	log.Printf("Server starting on %v", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server on port %v: %v", port, err)
	}

}

// Loads private key generated in Dockerfile from keygen.go
func init() {
	privateKeyPath := filepath.Join("/app/server", "private.pem")

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key from %s: %v", privateKeyPath, err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		log.Fatal("failed to decode PEM block containing private key")
	}

	var err2 error
	privateKey, err2 = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err2 != nil {
		log.Fatalf("Failed to parse private key: %v", err2)
	}

	log.Printf("Successfully loaded private key from %s", privateKeyPath)
}

func signBinary(filepath string) (string, error) {
	hash := sha256.New()
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for signing: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %v", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("failed to sign binary: %v", err)
	}

	return hex.EncodeToString(signature), nil
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

	// Create signature from binary
	signature, err := signBinary(binaryPath)
	if err != nil {
		log.Printf("Error signing binary: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := UpdateResponse{
		UpdateVersion: latestVersion,
		DownloadLink:  fmt.Sprintf("%s/download?file=%s", baseUrl, binaryName),
		Checksum:      checksum,
		Signature:     signature,
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
