package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	currentVersion = "1.1.1"
	checkUpdateUrl = "http://localhost:8080/checkUpdate"
)

type UpdateResponse struct {
	UpdateVersion string `json:"updateVersion"`
	DownloadLink  string `json:"downloadLink"`
}

func main() {
	log.Printf("Starting Self-Updating Application Version: %s\n", currentVersion)
	log.Printf("Application running. Press Ctrl+C to exit.")

	// Check if updated version exists
	checkUpdate()

}

// Check update version
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

	// Compare semantic versions
	if result.UpdateVersion == currentVersion {
		log.Printf("Application running latest version: %s\n", currentVersion)
		return
	}
	log.Printf("Downloading latest version...")
	log.Printf("Update version from result: %s", result.UpdateVersion)
	log.Printf("Download link: %s", result.DownloadLink)

	// Get update temp file path
	tmpFilePath, err := downloadUpdate(result.DownloadLink)

	if err != nil {
		log.Printf("Failed to download application update: %v", err)
		return
	}

	log.Printf("Successfully downloaded latest version.")

	// Clean up any previous temporary files
	defer os.Remove(tmpFilePath)

	log.Printf("Renaming executable...")

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Unable to find current executable: %v", err)
		return
	}
	backupPath := exePath + ".bak"
	err = os.Rename(exePath, backupPath)
	if err != nil {
		log.Printf("Error occurred while changing executable path: %v", err)
		return
	}
	log.Printf("Successfully backed up previous version to: %v", backupPath)

	os.Rename(tmpFilePath, exePath)
	if err != nil {
		log.Printf("Error occurred while changing tmp file path to executable path: %v", err)
		log.Printf("Rolling back executable to backup")
		_ = os.Rename(backupPath, exePath)
		return
	}
	log.Printf("Successfully replaced old binary with update!")

	log.Printf("Attempting to restart updated application...")

	// Restart application with update
	err = restartApplication(exePath)

	if err != nil {
		log.Printf("Unable to restart application at new executable path.")
		log.Printf("Rolling back executable to backup")
		_ = os.Rename(exePath, tmpFilePath+"_failed")
		_ = os.Rename(backupPath, exePath)

	}

	log.Printf("Application up-to-date. Running version: %v", result.UpdateVersion)
}

func downloadUpdate(downloadUrl string) (string, error) {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return "", fmt.Errorf("An error occurred getting updated app: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Update file could not be found: %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "myapp-update-")
	if err != nil {
		return "", fmt.Errorf("An error occurred creating tmp file: %v", err)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("An error occurred copying update to tmp file: %v", err)
	}

	err = tmpFile.Chmod(0755)
	if err != nil {
		return "", fmt.Errorf("An error occurred creating tmp file: %v", err)
	}

	// Return updated file path
	return tmpFile.Name(), nil

}

func restartApplication(filePath string) error {

	cmd := exec.Command(filePath)

	// Redirect logs of new process to current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("An error occurred while restarting the application: %v", err)
	}

	// End current process
	os.Exit(0)

	return nil
}
