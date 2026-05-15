package ai

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type AIAgent struct {
	llm      llms.Model
	provider *AIToolProvider
}

func NewAIAgent(llm llms.Model, provider *AIToolProvider) *AIAgent {
	return &AIAgent{
		llm:      llm,
		provider: provider,
	}
}

// Chat processes a single message through the AI agent with tool support
func (a *AIAgent) Chat(ctx context.Context, userMessage string) (string, error) {
	tools := a.provider.GetTools()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, userMessage),
	}

	// Try with tools first; fall back to plain generation on 404
	generateOpts := []llms.CallOption{llms.WithTools(tools)}

	for i := 0; i < 5; i++ {
		resp, err := a.llm.GenerateContent(ctx, content, generateOpts...)
		if err != nil {
			// 404 often means the endpoint doesn't support tool calling
			if strings.Contains(err.Error(), "404") && len(generateOpts) > 0 {
				log.Println("Warning: tools not supported by this endpoint, retrying without tools")
				generateOpts = []llms.CallOption{} // drop tools
				resp, err = a.llm.GenerateContent(ctx, content)
				if err != nil {
					return "", fmt.Errorf("failed to generate content: %w", err)
				}
			} else {
				return "", fmt.Errorf("failed to generate content: %w", err)
			}
		}
		choice := resp.Choices[0]
		// If there are tool calls, execute them
		if len(choice.ToolCalls) > 0 {
			// Add the assistant's response with tool calls to history
			parts := make([]llms.ContentPart, len(choice.ToolCalls))
			for j, tc := range choice.ToolCalls {
				parts[j] = tc
			}
			content = append(content, llms.MessageContent{
				Role:  llms.ChatMessageTypeAI,
				Parts: parts,
			})

			for _, tc := range choice.ToolCalls {
				log.Printf("Executing tool: %s with args: %s", tc.FunctionCall.Name, tc.FunctionCall.Arguments)
				result, err := a.provider.ExecuteTool(ctx, tc)
				if err != nil {
					result = fmt.Sprintf("Error executing tool: %v", err)
				}

				// Add tool result to history
				content = append(content, llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: tc.ID,
							Name:       tc.FunctionCall.Name,
							Content:    result,
						},
					},
				})
			}
			// Continue to next iteration to let LLM process the tool results
			continue
		}

		// No more tool calls, return the final response
		return choice.Content, nil
	}

	return "I'm sorry, I couldn't complete the task after several tool calls.", nil
}
