package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type QuestionGenerator interface {
	GenerateQuestion(ctx context.Context, prompt string) (string, error)
	Close() error
}

type GeminiQuestionGenerator struct {
	client *genai.Client
	model  string
}

func NewGeminiQuestionGenerator(ctx context.Context, apiKey, model string) (*GeminiQuestionGenerator, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("gemini api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "gemini-1.5-flash"
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiQuestionGenerator{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiQuestionGenerator) Close() error {
	if g == nil || g.client == nil {
		return nil
	}
	return g.client.Close()
}

func (g *GeminiQuestionGenerator) GenerateQuestion(ctx context.Context, prompt string) (string, error) {
	if g == nil || g.client == nil {
		return "", fmt.Errorf("gemini client is not initialized")
	}

	model := g.client.GenerativeModel(g.model)
	temp := float32(0.4)
	model.Temperature = &temp
	maxTokens := int32(64)
	model.MaxOutputTokens = &maxTokens

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	text := extractFirstText(resp)
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("empty response from gemini")
	}
	text = strings.Split(text, "\n")[0]
	text = strings.TrimSpace(text)

	return text, nil
}

func extractFirstText(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if t, ok := part.(genai.Text); ok {
				return string(t)
			}
		}
	}
	return ""
}
