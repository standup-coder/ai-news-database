package mocks

import "news4coder/internal/llm"

// LLMClientMock is a mock implementation of llm.LLMClient for testing
type LLMClientMock struct {
	SimpleChatFunc   func(prompt string, maxTokens int) (string, error)
	ChatFunc         func(messages []llm.Message, maxTokens int) (string, error)
	GetEmbeddingFunc func(text string) ([]float64, error)
}

func (m *LLMClientMock) SimpleChat(prompt string, maxTokens int) (string, error) {
	if m.SimpleChatFunc != nil {
		return m.SimpleChatFunc(prompt, maxTokens)
	}
	return "", nil
}

func (m *LLMClientMock) Chat(messages []llm.Message, maxTokens int) (string, error) {
	if m.ChatFunc != nil {
		return m.ChatFunc(messages, maxTokens)
	}
	return "", nil
}

func (m *LLMClientMock) GetEmbedding(text string) ([]float64, error) {
	if m.GetEmbeddingFunc != nil {
		return m.GetEmbeddingFunc(text)
	}
	return nil, nil
}