# Telegram AI Agent Integration

**Description**: This feature introduces a Telegram bot powered by an AI Agent that can execute various administrative and DevOps operations in Devopin using conversational text. The AI uses the `langchaingo` framework to execute function calling/tools.

## Objective
Implement an AI agent accessible via a Telegram worker/forwarder. The agent will respond to commands by using predefined tools, allowing users to query metrics, manage workers, manage Docker containers, and analyze project logs.

## Requirements
- Support multiple AI Providers: Gemini, OpenAI, HuggingFace (Free), and Custom Deployed AI (OpenAI-compatible endpoints like Ollama, vLLM, etc.). The user can select their preferred model via the application settings.
- Implement the agent as a background worker process.
- Tools required for the agent:
  - System Metrics: View CPU, Disk, and Memory usage.
  - Data Analysis: Analyze system usages.
  - Worker Management: List workers, Start, Stop, Restart workers.
  - Docker Management: List containers, Start, Stop, Restart containers.
  - Log Analyzer: Fetch project logs, summarize and analyze errors/logs.
- Library: `github.com/tmc/langchaingo`.

---

## Step-by-Step Implementation Guide

### 1. Update Settings Model and Migration
Update the existing application settings to store AI configuration and Telegram Bot Token.

**Action**: Add the following fields to the `Setting` struct in the database models.
```go
// models/setting.go
type Setting struct {
    // ...existing fields
    TelegramBotToken string `gorm:"type:varchar(255)" json:"telegram_bot_token"`
    AIProvider       string `gorm:"type:varchar(50);default:'gemini'" json:"ai_provider"` // 'gemini', 'openai', 'huggingface', 'custom'
    AIModelName      string `gorm:"type:varchar(100)" json:"ai_model_name"` // e.g. 'gemini-1.5-pro', 'gpt-4', or custom model name
    AIApiKey         string `gorm:"type:varchar(255)" json:"ai_api_key"`
    AIBaseURL        string `gorm:"type:varchar(255)" json:"ai_base_url"` // For custom user-deployed AI endpoints
}
```
*Note: Perform GORM auto-migration on server start. Ensure frontend state is updated to handle these inputs in `pages/settings.vue`.*

### 2. Initialize LangChainGo LLM Client Factory
Create a factory service that generates the appropriate LangChain LLM instance based on the current settings.

**Action**: Create `services/ai/factory.go`.
```go
// Example skeleton code for LLM Factory
package ai

import (
    "context"
    "errors"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
    "github.com/tmc/langchaingo/llms/googleai"
    "github.com/tmc/langchaingo/llms/huggingface"
)

func NewLLMClient(ctx context.Context, provider, apiKey, modelName, baseURL string) (llms.Model, error) {
    switch provider {
    case "openai":
        return openai.New(openai.WithToken(apiKey), openai.WithModel(modelName))
    case "gemini":
        return googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(modelName))
    case "huggingface":
        return huggingface.New(huggingface.WithToken(apiKey), huggingface.WithModel(modelName))
    case "custom":
        // For user-deployed models that are OpenAI-API compatible (e.g., Ollama, vLLM, LocalAI)
        return openai.New(
            openai.WithToken(apiKey), 
            openai.WithModel(modelName), 
            openai.WithBaseURL(baseURL),
        )
    default:
        return nil, errors.New("unsupported AI provider")
    }
}
```

### 3. Implement AI Agent Tools (Function Calling)
Create standard LangChainGo tools that the LLM can use. Each tool should wrap existing service methods (e.g., Docker Service, Worker Service).

**Action**: Create tools in `services/ai/tools/`.
- `ToolSystemMetrics`: Returns latest `ThresholdLog` data (CPU, Mem, Disk).
- `ToolWorkerManager`: Takes an action (`list`, `start`, `stop`, `restart`) and target worker name.
- `ToolDockerManager`: Takes an action (`list`, `start`, `stop`, `restart`) and target container ID.
- `ToolLogAnalyzer`: Fetches logs from the database for a specific project and returns a summarized string.

*Best Practice: When fetching logs for the analyzer, avoid fetching millions of rows. Use `LIMIT 500` or time-based filters. Avoid N+1 query structures by eagerly loading relationships.*

### 4. Create the Telegram Bot Worker Service
Create a worker that initializes a Telegram Bot using `gopkg.in/telebot.v3` (or similar) and listens to messages.

**Action**: Implement `cmd/cli/telegram_worker.go` or a background goroutine in the server.
- The worker listens for messages `On(tb.OnText, ...)`
- For every message, it passes the text to the LangChain Agent.
- The Agent evaluates the tools, runs them, and produces a final answer.
- The worker sends the generated answer back to the user via Telegram.

### 5. Add Frontend Settings Interface
Allow users to configure the Telegram Bot Token and AI configurations in the UI.

**Action**: In `frontend/src/pages/settings.vue`:
- Add a new Card or Expansion Panel for "AI & Telegram Configuration".
- Add fields: Provider (Select), Model Name (Input), API Key (Input type password), Telegram Bot Token (Input).
- Add field: `Base URL` (Input, only visible if Provider is 'custom').
- Ensure saving calls the existing `settings` API.

---

## Important Constraints & Best Practices
- **No bloated queries**: Ensure that metrics analysis tools fetch aggregated data (using SQL aggregates or Materialized views we set up earlier) instead of fetching thousands of raw rows and processing in memory.
- **Security**: The Telegram Bot should check if the `ChatID` or `Username` is authorized. Do not allow public access to server controls. Add a whitelist feature in settings or hardcode checking if needed.
- **Wait for command**: As per workflow guidelines, DO NOT implement these steps automatically until explicitly instructed to proceed by the user.
