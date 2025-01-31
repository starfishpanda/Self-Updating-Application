package main

import (
	"crypto"
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
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	currentVersion = "1.1.1"
	checkUpdateUrl = "http://localhost:8080/checkUpdate"
	publicKey      *rsa.PublicKey
)

type UpdateResponse struct {
	UpdateVersion string `json:"updateVersion"`
	DownloadLink  string `json:"downloadLink"`
	Checksum      string `json:"checksum"`
	Signature     string `json:"signature"`
}

func main() {
	// Setup signal handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Printf("Self-Updating Application is running version: %s\n", currentVersion)
	// Check if updated version exists ever 4 seconds
	go func() {
		for {
			checkUpdate()
			time.Sleep(4 * time.Second)
		}
	}()

	<-stop
	log.Println("Received interrupt signal. Shutting down gracefully...")
	os.Exit(0)
}

func init() {
	publicKeyPath := filepath.Join("/app/client", "public.pem")

	// Read the public key file
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("Failed to read public key from %s: %v", publicKeyPath, err)
	}

	// Decode the PEM formatted key
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		log.Fatal("Failed to decode PEM block containing public key")
	}

	// Parse the key into our publicKey variable
	var err2 error
	publicKey, err2 = x509.ParsePKCS1PublicKey(block.Bytes)
	if err2 != nil {
		log.Fatalf("Failed to parse public key: %v", err2)
	}

	log.Printf("Successfully loaded public key from %s", publicKeyPath)

}

func verifySignature(filepath string, signature string) (bool, error) {
	// Calculate the hash of the downloaded binary
	hash := sha256.New()
	file, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("failed to open file for verification: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("failed to calculate hash: %v", err)
	}

	// Decode the hex-encoded signature to bytes
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Verify the signature using public key
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), signatureBytes)
	return err == nil, err
}

func checkUpdate() {
	log.Printf("Checking for updates...")

	// GET check for updates
	resp, err := http.Get(checkUpdateUrl)
	if err != nil {
		log.Printf("Update check failed: %v", err)
		return
	}
	// Close TCP connection when function exits
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Update check returned status code: %v", resp.StatusCode)
		return
	}

	var result UpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("There was an error decoding JSON: %v", err)
		return
	}

	// Compare semantic versions of current and Update
	if result.UpdateVersion == currentVersion {
		log.Printf("Application running latest version: %s\n", currentVersion)
		return
	}
	log.Printf("New version available: %s", result.UpdateVersion)

	// Download update
	tmpPath, err := downloadUpdate(result.DownloadLink)
	if err != nil {
		log.Printf("Failed to download update: %v", err)
		return
	}

	defer os.Remove(tmpPath)
	log.Printf("Successfully downloaded update.")

	isValidSignature, err := verifySignature(tmpPath, result.Signature)
	if err != nil {
		log.Printf("Failed to verify signature: %v", err)
		return
	}

	if !isValidSignature {
		log.Printf("Invalid signature: binary may be tampered with")
		return
	}
	log.Printf("Successfully verified binary signature")

	// Verify checksum of update
	isValidChecksum, err := verifyChecksum(tmpPath, result.Checksum)
	if err != nil {
		log.Printf("Failed to verifyChecksum: %v", err)
		return
	}

	if !isValidChecksum {
		log.Printf("checksum verification failed: file may be corrupted or tampered with.")
		return
	}
	log.Printf("Successfully verified checksum of update.")

	newBinaryPath := filepath.Join("/app/client", "myapp-update")
	if err := os.Rename(tmpPath, newBinaryPath); err != nil {
		log.Printf("Failed to move update to binaries: %v", err)
		return
	}

	// Ensure update binary is executable
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		log.Printf("failed to set executable permissions: %v", err)
		return
	}
	cmd := exec.Command(newBinaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start new version: %v", err)
		return
	}

	log.Printf("Successfully started new version: %s", result.UpdateVersion)
	os.Exit(0)
}

func downloadUpdate(downloadLink string) (string, error) {
	resp, err := http.Get(downloadLink)
	if err != nil {
		return "", fmt.Errorf("download failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status code: %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "myapp-update-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write update: %v", err)
	}
	tmpFile.Close()

	return tmpFile.Name(), nil

}

func verifyChecksum(filepath string, expectedChecksum string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("failed to open file for verification: %v", err)
	}

	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("failed to calculate checksum: %v", err)
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	return actualChecksum == expectedChecksum, nil
}
