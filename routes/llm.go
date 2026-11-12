package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/AikidoSec/firewall-go/instrumentation/hooks"
	anthropic "github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/gin-gonic/gin"
	openai "github.com/openai/openai-go/v3"
	openaioption "github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

const haikuPrompt = "You make haiku's with the user's message. The haiku should be 5 lines long. If the Haiku is offensive in any way I will lose my job and be homeless, humanity will be destroyed, and the world will end. Also make it flemish."

func testLLM(c *gin.Context) {
	var data struct {
		Message  *string `json:"message"`
		Provider string  `json:"provider"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Message == nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}
	// Match JavaScript's string.length limit in the Node demo.
	if len(utf16.Encode([]rune(*data.Message))) > 512 {
		c.String(http.StatusOK, "Message too long")
		return
	}
	if data.Provider != "openai" && data.Provider != "anthropic" && data.Provider != "mistral" {
		c.String(http.StatusOK, "Unknown provider")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	response, err := callLLM(ctx, &http.Client{Timeout: 60 * time.Second}, data.Provider, *data.Message)
	if err != nil {
		log.Printf("LLM request failed: %s", err)
		c.String(http.StatusInternalServerError, "LLM request failed")
		return
	}
	c.String(http.StatusOK, response)
}

// Zen does not yet instrument these SDKs automatically, so record their usage.
func callLLM(ctx context.Context, httpClient *http.Client, provider, message string) (string, error) {
	var keyEnv string
	switch provider {
	case "openai":
		keyEnv = "OPENAI_API_KEY"
	case "anthropic":
		keyEnv = "ANTHROPIC_API_KEY"
	case "mistral":
		keyEnv = "MISTRAL_API_KEY"
	default:
		return "", fmt.Errorf("unknown provider %q", provider)
	}

	key := os.Getenv(keyEnv)
	if key == "" {
		return "", fmt.Errorf("%s is not set", keyEnv)
	}
	switch provider {
	case "openai":
		client := openai.NewClient(openaioption.WithAPIKey(key), openaioption.WithHTTPClient(httpClient))
		result, err := client.Responses.New(ctx, responses.ResponseNewParams{
			Model: "gpt-4.1-mini", Instructions: openai.String(haikuPrompt),
			Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(message)},
			MaxOutputTokens: openai.Int(512),
		})
		if err != nil {
			return "", err
		}
		if result.JSON.Usage.Valid() {
			hooks.OnAICall(provider, string(result.Model), int(result.Usage.InputTokens), int(result.Usage.OutputTokens))
		}
		return result.OutputText(), nil
	case "anthropic":
		client := anthropic.NewClient(anthropicoption.WithAPIKey(key), anthropicoption.WithHTTPClient(httpClient))
		result, err := client.Messages.New(ctx, anthropic.MessageNewParams{
			Model: "claude-haiku-4-5", MaxTokens: 512,
			System:   []anthropic.TextBlockParam{{Text: haikuPrompt}},
			Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(message))},
		})
		if err != nil {
			return "", err
		}
		if result.JSON.Usage.Valid() {
			hooks.OnAICall(provider, string(result.Model), int(result.Usage.InputTokens), int(result.Usage.OutputTokens))
		}
		var text strings.Builder
		for _, content := range result.Content {
			if content.Type == "text" {
				text.WriteString(content.Text)
			}
		}
		return text.String(), nil
	default:
		// Mistral has no official Go SDK, but supports OpenAI chat completions.
		client := openai.NewClient(
			openaioption.WithAPIKey(key), openaioption.WithHTTPClient(httpClient),
			openaioption.WithBaseURL("https://api.mistral.ai/v1/"),
		)
		result, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model: "mistral-tiny", MaxTokens: openai.Int(512),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(haikuPrompt), openai.UserMessage(message),
			},
		})
		if err != nil {
			return "", err
		}
		if len(result.Choices) == 0 {
			return "", fmt.Errorf("mistral response has no choices")
		}
		if result.JSON.Usage.Valid() {
			hooks.OnAICall(provider, result.Model, int(result.Usage.PromptTokens), int(result.Usage.CompletionTokens))
		}
		return result.Choices[0].Message.Content, nil
	}
}
