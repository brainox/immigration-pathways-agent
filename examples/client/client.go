package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	agentURL := "http://localhost:8080"
	
	if len(os.Args) < 2 {
		fmt.Println("Migration Pathways Agent - Test Client")
		fmt.Println("======================================")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run client.go <command> [args]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  card               - Get agent card")
		fmt.Println("  query \"<text>\"     - Send a migration query")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run client.go card")
		fmt.Println("  go run client.go query \"I'm a software engineer from Nigeria, want to move to Canada\"")
		fmt.Println("  go run client.go query \"Data scientist looking to relocate to USA with $5000 budget\"")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "card":
		getAgentCard(agentURL)
	case "query":
		if len(os.Args) < 3 {
			fmt.Println("Error: query command requires a text argument")
			fmt.Println("Example: go run client.go query \"I'm a software engineer wanting to move to Canada\"")
			os.Exit(1)
		}
		query := os.Args[2]
		sendQuery(agentURL, query)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func getAgentCard(agentURL string) {
	resp, err := http.Get(agentURL + "/.well-known/agent.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "  ")
	
	fmt.Println("🤖 Agent Card")
	fmt.Println("=============")
	fmt.Println(prettyJSON.String())
}

func sendQuery(agentURL string, query string) {
	// Create JSON-RPC request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tasks/send",
		"params": map[string]interface{}{
			"message": map[string]interface{}{
				"role": "user",
				"parts": []map[string]interface{}{
					{
						"type": "text",
						"text": query,
					},
				},
			},
		},
		"id": 1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Send request
	resp, err := http.Post(agentURL+"/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Parse response
	var rpcResponse struct {
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			ID        string `json:"id"`
			Status    struct {
				State   string `json:"state"`
				Message string `json:"message"`
			} `json:"status"`
			Artifacts []struct {
				Parts []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"artifacts"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    string `json:"data"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	// Check for errors
	if rpcResponse.Error != nil {
		fmt.Printf("❌ Error: %s (code: %d)\n", rpcResponse.Error.Message, rpcResponse.Error.Code)
		if rpcResponse.Error.Data != "" {
			fmt.Printf("   Details: %s\n", rpcResponse.Error.Data)
		}
		os.Exit(1)
	}

	// Display result
	fmt.Println("📊 Migration Pathways Result")
	fmt.Println("===========================")
	fmt.Printf("Task ID: %s\n", rpcResponse.Result.ID)
	fmt.Printf("Status: %s\n\n", rpcResponse.Result.Status.State)

	if len(rpcResponse.Result.Artifacts) > 0 {
		for _, artifact := range rpcResponse.Result.Artifacts {
			for _, part := range artifact.Parts {
				if part.Type == "text" {
					fmt.Println(part.Text)
				}
			}
		}
	}
}
