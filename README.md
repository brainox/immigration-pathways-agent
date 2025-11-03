# Migration Pathways AI Agent

An A2A (Agent-to-Agent) protocol-compliant AI agent powered by **Google's Gemini LLM** that provides real-time, personalized migration pathway recommendations for professionals looking to relocate internationally.

## 🎯 What It Does

This agent uses **Gemini AI** to analyze your profile (profession, destination country, budget) and generates personalized migration pathways with:
- Current visa types and descriptions (based on 2024-2025 policies)
- Real-time processing duration and costs
- Success probability estimates
- Detailed requirements tailored to your situation
- Step-by-step guidance

**Key Difference:** Unlike traditional agents with static databases, this agent **queries Gemini LLM in real-time** for each request, ensuring you get the most current and contextually relevant migration advice.

## 🏗️ System Architecture

```
┌─────────────────────┐
│    CLIENT/USER      │
└─────────┬───────────┘
          │ HTTP
┌─────────▼───────────┐
│   A2A PROTOCOL      │
│ ┌─────────────────┐ │
│ │ Agent Card      │ │
│ │ JSON-RPC 2.0    │ │
│ └─────────────────┘ │
└─────────┬───────────┘
          │
┌─────────▼───────────┐
│  MIGRATION AGENT    │
│ ┌─────────────────┐ │
│ │ Query Parser    │ │
│ │ Task Manager    │ │
│ │ Gemini LLM      │ │
│ └─────────────────┘ │
└─────────────────────┘
```

### Key Components:

1. **A2A Server** - JSON-RPC 2.0 over HTTP with task management
2. **Agent Card** - Capabilities at `/.well-known/agent.json`
3. **Gemini Integration** - Real-time AI pathway generation
4. **Query Parser** - Natural language profile extraction

### Universal Support:
- **Countries:** Any destination worldwide via Gemini LLM
- **Professions:** All occupations (tech, healthcare, education, trades, etc.)

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- **Gemini API Key** (free from Google AI Studio)
- VS Code with REST Client extension (for API testing)

### Installation

1. **Get your Gemini API Key** (free):
   - Visit https://aistudio.google.com/app/apikey
   - Click "Get API key" or "Create API key"
   - Copy your API key

2. **Set up environment**:
   ```bash
   # Copy example env file
   cp .env.example .env
   
   # Edit .env with your API key
   vim .env  # or use any editor
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Start the server:**
   ```bash
   go run cmd/server/main.go cmd/server/pathways.go cmd/server/a2a_types.go
   ```
   Server starts at `http://localhost:8080`

### Testing the API

Use the provided `.http` files in `api_tests/` directory:

1. **Get Agent Card:**
   - Open `api_tests/agent_card_test.http`
   - Click "Send Request" above the request

2. **Test Migration Query:**
   - Open `api_tests/tasks_send_test.http`
   - Click "Send Request" to try sample queries
   - Copy the returned task ID

3. **Get Results:**
   - Open `api_tests/tasks_get_test.http`
   - Replace `@taskId` with your task ID
   - Click "Send Request" to see results

Example Queries:
```
"I'm a software engineer from Nigeria, want to move to Canada"
"Data scientist looking to relocate to USA with $5000 budget"
"How can I migrate to Germany as a software developer?"
"Nurse from India wanting to move to UK"
```

## 📋 Implementation Details

### Project Structure
```
migration-pathways-agent/
├── cmd/
│   └── server/           # Main server implementation
│       ├── main.go      # A2A server + handlers
│       ├── pathways.go  # Gemini integration
│       └── a2a_types.go # Protocol types
│
├── api_tests/           # HTTP test files
│   ├── agent_card_test.http
│   ├── tasks_send_test.http
│   └── tasks_get_test.http
│
├── examples/            # Reference implementations
│   └── client/         # CLI test client
│
├── .env.example        # Environment template
└── README.md
```

### A2A Protocol Features

1. **Agent Discovery**
   - Agent Card at `/.well-known/agent.json`
   - Declares capabilities and skills

2. **Task Management (JSON-RPC 2.0)**
   - `tasks/send` - Submit migration queries
   - `tasks/get` - Retrieve results
   - Task state tracking and history

3. **Gemini Integration**
   - Real-time AI query processing
   - Current immigration policy knowledge
   - Contextual response generation

4. **Query Processing**
   - Natural language parsing
   - Profile extraction (profession, origin, destination)
   - Budget awareness
   - Success probability calculation

### 3. Message Format
- ✅ Supports text-based input/output
- ✅ Role-based messaging (user/agent)
- ✅ Structured response artifacts

## 🔌 API Examples

### Get Agent Card
```bash
curl http://localhost:8080/.well-known/agent.json | jq .
```

### Send a Task (JSON-RPC)
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tasks/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [{
          "type": "text",
          "text": "I am a software engineer from Kenya wanting to move to Canada with $10,000"
        }]
      }
    },
    "id": 1
  }' | jq .
```

### Get Task Status
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tasks/get",
    "params": {
      "id": "task-id-here"
    },
    "id": 2
  }' | jq .
```

## 🛠️ Extending the Agent

### Adding New Countries

Edit `pathways.go` and add entries to the `NewPathwayDatabase()` function:

```go
db.Pathways["france"] = make(map[string][]MigrationPathway)
db.Pathways["france"]["software_engineer"] = []MigrationPathway{
    {
        Name:        "French Tech Visa",
        Description: "Fast-track visa for tech workers",
        Duration:    "2-3 months",
        Cost:        "$200 - $500",
        // ... add requirements
    },
}
```

### Adding New Professions

Simply add new profession keys to existing countries:

```go
db.Pathways["canada"]["doctor"] = []MigrationPathway{
    // Add medical profession pathways
}
```

### Enhancing Query Parsing

The `parseUserQuery()` function in `main.go` can be enhanced with:
- NLP libraries for better text understanding
- LLM integration for complex query parsing
- Support for more attributes (education level, years of experience, etc.)

## 📊 Project Structure

```
migration-pathways-agent/
├── main.go          # A2A server implementation
├── pathways.go      # Migration pathways database
├── a2a_types.go     # A2A protocol types
├── client.go        # Test client
├── go.mod           # Go module file
└── README.md        # This file
```

## 🌐 A2A Protocol Resources

- [A2A Protocol Documentation](https://a2a-protocol.org/latest/)
- [A2A Specification](https://a2a-protocol.org/latest/specification/)
- [A2A GitHub Repository](https://github.com/a2aproject/A2A)

## 🔒 Security Considerations

This is a demonstration agent. For production use:

1. Add authentication (OAuth2, API keys)
2. Implement rate limiting
3. Add input validation and sanitization
4. Use HTTPS/TLS
5. Implement proper error handling
6. Add logging and monitoring

## 🤝 Contributing

This agent can be enhanced with:
- More countries and visa types
- Real-time data from immigration APIs
- Cost calculators with currency conversion
- Document checklist generation
- Timeline tracking features
- Multi-agent collaboration (e.g., consulting a visa expert agent)

## Example Interaction

```
User Query: "I'm a software engineer from Nigeria, want to move to Canada with $5000 budget"

Agent Response:
# Migration Pathways for Software Engineer to Canada

Based on your profile from Nigeria with a budget of $5000, here are your migration options:

## 1. Express Entry (Federal Skilled Worker)
**Visa Type:** Permanent Residence

**Description:** Point-based immigration system for skilled workers

**Duration:** 6-12 months
**Cost:** $2,300 - $3,500
**Success Probability:** High (if CRS score > 470)

**Requirements:**
- Bachelor's degree or higher
- At least 1 year of work experience
- Language test (IELTS/CELPIP)
- Educational Credential Assessment (ECA)
- Minimum 67 points out of 100

## 2. Provincial Nominee Program (PNP)
[... additional pathways ...]

**Next Steps:**
1. Review each pathway and assess which aligns best with your situation
2. Check official immigration websites for the most current requirements
3. Prepare required documents
4. Consider consulting with an immigration lawyer
```