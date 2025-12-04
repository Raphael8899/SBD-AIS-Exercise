package main

import (
	"exc8/client"
	"exc8/server"
	"fmt"
	"time"
)

func main() {
	go func() {
		// start server
		if err := server.StartGrpcServer(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	time.Sleep(1 * time.Second)

	// start client
	c, err := client.NewGrpcClient()
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	if err := c.Run(); err != nil {
		fmt.Printf("Client error: %v\n", err)
	}

	fmt.Println("Orders complete!")
}
