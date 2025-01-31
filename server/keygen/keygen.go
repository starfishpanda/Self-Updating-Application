// keygen/main.go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Save private key
	if err := os.WriteFile("private.pem", privateKeyPEM, 0600); err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	// Encode public key to PEM format
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Save public key
	if err := os.WriteFile("public.pem", publicKeyPEM, 0644); err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

	log.Println("Successfully generated key pair")
}
