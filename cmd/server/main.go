package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MigrationAgent is the main agent server
type MigrationAgent struct {
	gemini *GeminiClient
	tasks  map[string]*Task
	mu     sync.RWMutex
}

// NewMigrationAgent creates a new migration pathways agent
func NewMigrationAgent() *MigrationAgent {
	return &MigrationAgent{
		gemini: NewGeminiClient(),
		tasks:  make(map[string]*Task),
	}
}

// GetAgentCard returns the agent's metadata
func (a *MigrationAgent) GetAgentCard() *AgentCard {
	return &AgentCard{
		Name:        "Migration Pathways Agent",
		Description: "An AI-powered agent that provides real-time, personalized migration pathway recommendations using Gemini LLM. Get current visa options, costs, requirements, and success probabilities based on your profile and destination country.",
		URL:         "http://localhost:8080",
		Version:     "2.0.0",
		Capabilities: Capabilities{
			Streaming:              false,
			PushNotifications:      false,
			StateTransitionHistory: false,
		},
		DefaultInputModes:  []string{"text", "text/plain"},
		DefaultOutputModes: []string{"text", "text/plain"},
		Skills: []Skill{
			{
				ID:          "get_migration_pathways",
				Name:        "Get AI-Generated Migration Pathways",
				Description: "Provides real-time, AI-generated migration pathways based on profession, destination country, origin, and budget using Google's Gemini LLM",
				Tags:        []string{"migration", "visa", "relocation", "immigration", "ai-powered", "real-time"},
				Examples: []string{
					"I'm a software engineer from Nigeria, want to move to Canada, budget $5000",
					"Data scientist looking to relocate to USA",
					"How can I migrate to Germany as a software developer?",
					"What are my options to move to Australia as an engineer with $10k budget?",
					"Nurse from India wanting to move to UK",
					"Teacher relocating from Philippines to Canada",
				},
			},
		},
	}
}

// ProcessTask handles incoming tasks
func (a *MigrationAgent) ProcessTask(taskID string, message Message) (*Task, error) {
	// Generate a message ID
	messageID := uuid.New().String()

	// Create task
	task := &Task{
		ID:   taskID,
		Kind: "task",
		Status: TaskStatus{
			State:     "working",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store task
	a.mu.Lock()
	a.tasks[taskID] = task
	a.mu.Unlock()

	// Extract text from message
	var userQuery string
	for _, part := range message.Parts {
		if part.Type == "text" {
			userQuery += part.Text + " "
		}
	}
	userQuery = strings.TrimSpace(userQuery)

	// Parse user query to extract: profession, destination, origin, budget
	profile := a.parseUserQuery(userQuery)

	// Query Gemini LLM for migration pathways
	responseText, err := a.gemini.GetMigrationPathways(
		profile.Profession,
		profile.Destination,
		profile.Origin,
		profile.Budget,
	)

	if err != nil {
		// Update task with error
		task.Status = TaskStatus{
			State:     "failed",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Message: &StatusMessage{
				Kind: "message",
				Role: "agent",
				Parts: []Part{
					{
						Kind: "text",
						Text: fmt.Sprintf("Failed to generate pathways: %v", err),
					},
				},
				MessageID: messageID,
				TaskID:    taskID,
			},
		}
		task.UpdatedAt = time.Now()

		a.mu.Lock()
		a.tasks[taskID] = task
		a.mu.Unlock()

		return task, err
	}

	// Generate artifact ID
	artifactID := uuid.New().String()

	// Update task with result
	task.Status = TaskStatus{
		State:     "completed",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message: &StatusMessage{
			Kind: "message",
			Role: "agent",
			Parts: []Part{
				{
					Kind: "text",
					Text: responseText,
				},
			},
			MessageID: messageID,
			TaskID:    taskID,
		},
	}
	task.Artifacts = []Artifact{
		{
			ArtifactID: artifactID,
			Name:       "Migration Pathway Recommendation",
			Parts: []Part{
				{
					Kind: "text",
					Text: responseText,
				},
			},
		},
	}
	task.UpdatedAt = time.Now()

	// Update stored task
	a.mu.Lock()
	a.tasks[taskID] = task
	a.mu.Unlock()

	return task, nil
}

// UserProfile represents parsed user information
type UserProfile struct {
	Profession  string
	Destination string
	Budget      int
	Origin      string
}

// parseUserQuery extracts information from user's natural language query
func (a *MigrationAgent) parseUserQuery(query string) UserProfile {
	queryLower := strings.ToLower(query)
	profile := UserProfile{
		Budget: 0, // 0 means no budget specified
	}

	// Extract profession (common ones)
	professions := map[string]string{
		"software engineer": "software engineer",
		"data scientist":    "data scientist",
		"engineer":          "engineer",
		"developer":         "developer",
		"programmer":        "programmer",
		"doctor":            "doctor",
		"nurse":             "nurse",
		"teacher":           "teacher",
		"accountant":        "accountant",
		"designer":          "designer",
		"manager":           "manager",
		"analyst":           "analyst",
		"consultant":        "consultant",
	}

	for key, value := range professions {
		if strings.Contains(queryLower, key) {
			profile.Profession = value
			break
		}
	}

	// Extract destination country
	countries := map[string]string{
		"canada":         "Canada",
		"usa":            "USA",
		"united states":  "USA",
		"america":        "USA",
		"uk":             "UK",
		"united kingdom": "UK",
		"britain":        "UK",
		"germany":        "Germany",
		"australia":      "Australia",
		"france":         "France",
		"netherlands":    "Netherlands",
		"sweden":         "Sweden",
		"norway":         "Norway",
		"denmark":        "Denmark",
		"switzerland":    "Switzerland",
		"new zealand":    "New Zealand",
		"singapore":      "Singapore",
		"japan":          "Japan",
		"south korea":    "South Korea",
		"dubai":          "UAE",
		"uae":            "UAE",
	}

	for key, value := range countries {
		if strings.Contains(queryLower, key) {
			profile.Destination = value
			break
		}
	}

	// Extract origin country
	origins := map[string]string{
		"nigeria":      "Nigeria",
		"ghana":        "Ghana",
		"kenya":        "Kenya",
		"south africa": "South Africa",
		"ethiopia":     "Ethiopia",
		"egypt":        "Egypt",
		"morocco":      "Morocco",
		"tanzania":     "Tanzania",
		"uganda":       "Uganda",
		"india":        "India",
		"pakistan":     "Pakistan",
		"bangladesh":   "Bangladesh",
		"philippines":  "Philippines",
		"china":        "China",
		"brazil":       "Brazil",
		"mexico":       "Mexico",
		"argentina":    "Argentina",
	}

	for key, value := range origins {
		if strings.Contains(queryLower, key) {
			profile.Origin = value
			break
		}
	}

	// Extract budget - look for dollar amounts
	// Patterns: $5000, $5,000, 5000 dollars, 5k
	if idx := strings.Index(queryLower, "$"); idx != -1 {
		// Extract number after $
		numStr := ""
		for i := idx + 1; i < len(queryLower); i++ {
			c := queryLower[i]
			if (c >= '0' && c <= '9') || c == ',' || c == '.' {
				if c != ',' { // ignore commas
					numStr += string(c)
				}
			} else {
				break
			}
		}

		// Handle 'k' suffix (e.g., $5k = $5000)
		if idx+len(numStr)+1 < len(queryLower) && queryLower[idx+len(numStr)+1] == 'k' {
			if budget, err := fmt.Sscanf(numStr, "%d", &profile.Budget); err == nil && budget == 1 {
				profile.Budget *= 1000
			}
		} else {
			fmt.Sscanf(numStr, "%d", &profile.Budget)
		}
	}

	return profile
}

// GetTask retrieves a task by ID
func (a *MigrationAgent) GetTask(taskID string) (*Task, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	task, exists := a.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// ServeHTTP handles HTTP requests
func (a *MigrationAgent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle agent card endpoint
	if r.URL.Path == "/.well-known/agent.json" && r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a.GetAgentCard())
		return
	}

	// Handle JSON-RPC requests
	if r.Method == "POST" {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			a.sendError(w, nil, -32700, "Parse error", req.ID)
			return
		}

		switch req.Method {
		case "tasks/send":
			a.handleTasksSend(w, req)
		case "tasks/get":
			a.handleTasksGet(w, req)
		default:
			a.sendError(w, nil, -32601, "Method not found", req.ID)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleTasksSend processes tasks/send RPC method
func (a *MigrationAgent) handleTasksSend(w http.ResponseWriter, req JSONRPCRequest) {
	// Parse params
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		a.sendError(w, err, -32602, "Invalid params", req.ID)
		return
	}

	var params TaskSendParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		a.sendError(w, err, -32602, "Invalid params", req.ID)
		return
	}

	// Generate task ID if not provided
	taskID := params.ID
	if taskID == "" {
		taskID = uuid.New().String()
	}

	// Process task
	task, err := a.ProcessTask(taskID, params.Message)
	if err != nil {
		a.sendError(w, err, -32603, "Internal error", req.ID)
		return
	}

	// Send response
	a.sendSuccess(w, task, req.ID)
}

// handleTasksGet processes tasks/get RPC method
func (a *MigrationAgent) handleTasksGet(w http.ResponseWriter, req JSONRPCRequest) {
	// Parse params
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		a.sendError(w, err, -32602, "Invalid params", req.ID)
		return
	}

	var params TaskIDParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		a.sendError(w, err, -32602, "Invalid params", req.ID)
		return
	}

	// Get task
	task, err := a.GetTask(params.ID)
	if err != nil {
		a.sendError(w, err, -32602, err.Error(), req.ID)
		return
	}

	// Send response
	a.sendSuccess(w, task, req.ID)
}

// sendSuccess sends a successful JSON-RPC response
func (a *MigrationAgent) sendSuccess(w http.ResponseWriter, result interface{}, id interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error JSON-RPC response
func (a *MigrationAgent) sendError(w http.ResponseWriter, err error, code int, message string, id interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    nil,
		},
		ID: id,
	}
	if err != nil {
		response.Error.Data = err.Error()
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	agent := NewMigrationAgent()

	// Check if API key is set
	if agent.gemini.APIKey == "" {
		log.Println("⚠️  WARNING: GEMINI_API_KEY environment variable not set!")
		log.Println("   Please set it with: export GEMINI_API_KEY=your-api-key")
		log.Println("   Get your key at: https://aistudio.google.com/app/apikey")
		log.Println("")
	}

	http.HandleFunc("/", agent.ServeHTTP)

	// Heroku (and other platforms) provide the port via the PORT env var.
	// Fall back to 8080 for local development.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("🚀 Migration Pathways Agent (AI-Powered) starting on %s", addr)
	log.Printf("📋 Agent Card available at: http://localhost:%s/.well-known/agent.json", port)
	log.Printf("🔗 A2A endpoint: http://localhost:%s/", port)
	log.Printf("🤖 Using Gemini LLM for real-time migration pathway generation")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
