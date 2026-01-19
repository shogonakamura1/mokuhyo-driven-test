package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type QuestionSelector interface {
	Select(ctx context.Context, prompt string, candidates []string) (string, error)
	Close() error
}

type GeminiQuestionSelector struct {
	client *genai.Client
	model  string
}

func NewGeminiQuestionSelector(ctx context.Context, apiKey, model string) (*GeminiQuestionSelector, error) {
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

	return &GeminiQuestionSelector{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiQuestionSelector) Close() error {
	if g == nil || g.client == nil {
		return nil
	}
	return g.client.Close()
}

func (g *GeminiQuestionSelector) Select(ctx context.Context, prompt string, candidates []string) (string, error) {
	if g == nil || g.client == nil {
		return "", fmt.Errorf("gemini client is not initialized")
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("candidates are empty")
	}

	model := g.client.GenerativeModel(g.model)
	temp := float32(0.2)
	model.Temperature = &temp
	maxTokens := int32(32)
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

	candidateSet := make(map[string]struct{}, len(candidates))
	for _, c := range candidates {
		candidateSet[c] = struct{}{}
	}
	if _, ok := candidateSet[text]; !ok {
		return "", fmt.Errorf("unexpected response: %s", text)
	}
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
