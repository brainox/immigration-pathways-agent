package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// GeminiClient handles communication with Gemini API
type GeminiClient struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewGeminiClient creates a new Gemini API client
func NewGeminiClient() *GeminiClient {
	// Attempt to load values from a local .env file (if present). Values
	// already set in the environment are not overwritten.
	loadDotEnv()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// Try common alternative env var names
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}

	return &GeminiClient{
		APIKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		Model:   "gemini-2.0-flash-exp", // Latest Gemini model
	}
}

// loadDotEnv reads a .env file from the current working directory and sets
// any variables that are not already present in the environment. The parser
// is minimal: lines starting with # are ignored; blank lines are skipped;
// KEY=VALUE pairs are supported (values may be quoted).
func loadDotEnv() {
	// If env already contains GEMINI_API_KEY we can skip entirely.
	if os.Getenv("GEMINI_API_KEY") != "" || os.Getenv("GOOGLE_API_KEY") != "" {
		return
	}

	f, err := os.Open(".env")
	if err != nil {
		// no .env present — nothing to do
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}

		if key == "" || val == "" {
			continue
		}

		// Only set if not already present
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	// ignore scanner.Err() errors silently — .env is optional
}

// GeminiRequest represents a request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content in a Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents a response from Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// GetMigrationPathways queries Gemini for migration pathway recommendations
func (gc *GeminiClient) GetMigrationPathways(profession, destination, origin string, budget int) (string, error) {
	if gc.APIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Construct the prompt for Gemini
	prompt := gc.buildPrompt(profession, destination, origin, budget)

	// Create request
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make API request
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", gc.BaseURL, gc.Model, gc.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract text from response
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated from API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// buildPrompt constructs the prompt for Gemini
func (gc *GeminiClient) buildPrompt(profession, destination, origin string, budget int) string {
	prompt := `You are a migration planning expert. Provide personalized migration pathway recommendations in a well-structured markdown format.

CRITICAL BEHAVIOR RULES:
- Never ask the user for additional information.
- If any profile fields are missing (profession, origin, destination, budget), proceed with best available information and reasonable assumptions.
- Output exactly one best migration option. Do not include follow-up questions.

USER PROFILE:
`
	if profession != "" {
		prompt += fmt.Sprintf("- Profession: %s\n", profession)
	}
	if origin != "" {
		prompt += fmt.Sprintf("- Current Country: %s\n", origin)
	}
	if destination != "" {
		prompt += fmt.Sprintf("- Destination Country: %s\n", destination)
	}
	if budget > 0 {
		prompt += fmt.Sprintf("- Budget: $%d USD\n", budget)
	}

	prompt += `
INSTRUCTIONS:
Research and provide the SINGLE most suitable migration pathway for this profile. If some fields are missing, infer typical constraints for 2024–2025 and proceed without asking questions. Format as follows:

# Best Migration Option: [Visa Name]

Brief overview of why this is the best option for the profile (1-2 sentences).

**Key Details:**
- Processing time: [Duration]
- Cost: [USD range]
- Success rate: [High/Medium/Low]
- Main requirements: [2-3 key points]

Next step: [Most important action to take]

IMPORTANT: Be concise. Focus on 2024-2025 requirements. Consider budget constraints. Do not ask for more details.

Generate the response now:`

	return prompt
}
