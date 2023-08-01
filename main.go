// main.go

package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Get("/ws", websocket.New(handleWebSocket))

	err := app.Listen(":3000")
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func handleWebSocket(c *websocket.Conn) {
	// Close the WebSocket connection when the function returns
	defer c.Close()

	// Read the initial ID sent by the client
	var id string
	err := c.ReadJSON(&id)
	if err != nil {
		fmt.Println("Error reading JSON from client:", err)
		return
	}

	fmt.Printf("Received ID from client: %s\n", id)

	dataFromAPI2 := queryAPI(id)

	// Send the data from API2 to the client once
	err = c.WriteJSON(dataFromAPI2)
	if err != nil {
		fmt.Println("Error writing JSON to client:", err)
		return
	}

	fmt.Printf("Sent data from API2 to client: %v\n", dataFromAPI2)

	// Create a channel to signal when the client wants to close the connection
	closeSignal := make(chan struct{})

	// Start a Go routine to send data to the client every 10 seconds
	go sendDataToClient(c, id, closeSignal)

	// Wait for the close signal from the client or an error in the connection
	_, _, err = c.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message from client:", err)
	}

	// Signal the sendDataToClient Go routine to stop sending data and return
	close(closeSignal)
}

func sendDataToClient(c *websocket.Conn, id string, closeSignal <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Query your API here with the ID
			// Replace the `queryAPI` function with your actual API query implementation
			dataFromAPI := queryAPI(id)

			// Send the data to the client
			err := c.WriteJSON(dataFromAPI)
			if err != nil {
				fmt.Println("Error writing JSON to client:", err)
				return
			}

			fmt.Printf("Sent data to client: %v\n", dataFromAPI)

		case <-closeSignal:
			// The client wants to close the connection, so return from the Go routine
			fmt.Println("Client requested to close the connection")
			return
		}
	}
}

// Replace this function with your actual API query implementation
func queryAPI(id string) interface{} {
	// Simulating the API response with a simple map
	return map[string]interface{}{
		"id":        id,
		"data":      "Some data from API",
		"timestamp": time.Now().Unix(),
	}
}
