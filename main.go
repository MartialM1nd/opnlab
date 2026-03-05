package main

import (
	"fmt"
	"log"

	"github.com/MartialM1nd/opnlab/internal/providers"
	"github.com/MartialM1nd/opnlab/internal/server"
)

func main() {
	srv := server.New()

	srv.RegisterProvider(providers.NewSystemProvider())

	// Start server
	addr := ":8080"
	log.Printf("Starting opnlab server on %s", addr)
	if err := srv.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Println("opnlab server started")
}
