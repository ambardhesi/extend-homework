package main

import (
	"log"
	"os"

	"github.com/ambardhesi/extend-homework/internal/server"
)

func main() {
	config := server.Config{
		Port:           8080,
		LogDir:         "logs",
		CertFilePath:   "certs/svr-cert.pem",
		KeyFilePath:    "certs/svr-key.pem",
		CaCertFilePath: "certs/ca-cert.pem",
		TestMode:       false,
	}

	s, err := server.NewServer(config)
	if err != nil {
		log.Printf("Failed to start server %v\n", err)
		os.Exit(1)
	}

	s.Start()
}