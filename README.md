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

- The service requires a Gemini API key in the `GEMINI_API_KEY` environment variable. You can either export it in your shell or create a local `.env` file and load it with your preferred tool.
- Example (temporary shell export):
```bash
export GEMINI_API_KEY=your_gemini_api_key_here
```

- Or create a `.env` file containing the same line and load it with `direnv`, `dotenv`, or your shell's source command.

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Start the server:**
   ```bash
   # Build & run the server package. The server reads $PORT and defaults to 8080.
   go run ./cmd/server
   ```
   By default the server listens on port `8080`. To run on a different port:
   ```bash
   PORT=3000 go run ./cmd/server
   ```

   For Heroku deployment details see `HEROKU_DEPLOY.md`.

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
├── Procfile            # Heroku process file (optional)
├── HEROKU_DEPLOY.md    # Heroku deployment instructions (optional)
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

### Adding New Countries / Professions

Current behavior: this service builds a prompt from the user's profile and queries the Gemini LLM to generate migration pathway recommendations at request time. Because responses are generated by the model, you don't need to maintain a local static database for most use cases.

If you'd like to provide a local, static fallback or supplement (for deterministic outputs or to seed the prompt), implement a small in-memory database and consult it from `ProcessTask` before calling Gemini. Example (high-level):

```go
// define a simple pathway struct (create this in a new file, e.g. cmd/server/pathways_db.go)
type MigrationPathway struct {
   Name        string
   Description string
   Duration    string
   Cost        string
   Requirements []string
}

// a simple in-memory DB
var pathwayDB = map[string]map[string][]MigrationPathway{
   "canada": {
      "software_engineer": {{Name: "Express Entry", Description: "Federal Skilled Worker", Duration: "6-12 months", Cost: "$2k-$4k", Requirements: []string{"degree","1+ year experience"}}},
   },
}

// in ProcessTask (main.go) - consult local DB first
if local, ok := lookupLocalPathway(profile); ok {
   // build task result from `local` (skip Gemini)
} else {
   // call a.gemini.GetMigrationPathways(...)
}
```

This approach keeps the Gemini integration intact while allowing deterministic, curated pathways when needed.

### Enhancing Query Parsing

The `parseUserQuery()` function in `main.go` can be enhanced with:
- NLP libraries for better text understanding
- LLM integration for complex query parsing
- Support for more attributes (education level, years of experience, etc.)

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

Note: the agent's responses are generated by the Gemini LLM at request time. The example below is illustrative — actual output will vary depending on the prompt, model, and up-to-date information returned by Gemini.

```
User Query: "I'm a software engineer from Nigeria, want to move to Canada with $5000 budget"

### Illustrative server JSON-RPC response (what the server returns)
```json
{
  "jsonrpc": "2.0",
  "result": {
    "id": "task-id",
    "status": {
      "state": "completed",
      "message": "Migration pathways generated successfully"
    },
    "artifacts": [
      {
        "parts": [
          {
            "type": "text",
            "text": "# Best Migration Option: H-1B Visa (Employer Sponsored)\n\nThis is the most realistic and potentially fastest pathway for a data scientist to the US given the budget, profession, and country of origin, assuming they can secure a US employer.\n\n**Key Details:**\n- Processing time: Varies; lottery-dependent, premium processing (~15 days) if selected, otherwise months.\n- Cost: $1500 - $7500 (Primarily employer-paid. Applicant may incur some fees for documentation and travel.)\n- Success rate: Medium (Lottery-based system; varies annually based on demand.)\n- Main requirements: Bachelor's degree (or equivalent) in a relevant field, US employer sponsorship for a \"specialty occupation\" (data science qualifies), prevailing wage requirement met.\n\nNext step: Actively network and apply for data science positions with US-based companies willing to sponsor H-1B visas.\n"
          }
        ]
      }
    ],
    "createdAt": "2025-11-03T00:00:00Z",
    "updatedAt": "2025-11-03T00:00:00Z"
  },
  "id": 1
}
```