package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"news4coder/internal/config"
	"time"
)

// Client LLM HTTP 客户端
type Client struct {
	cfg    *config.LLMConfig
	client *http.Client
}

// NewClient 创建 LLM 客户端
func NewClient(cfg *config.LLMConfig) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Chat 发送对话请求
func (c *Client) Chat(messages []Message, maxTokens int) (string, error) {
	if c.cfg.APIKey == "" {
		return "", fmt.Errorf("LLM API Key 未配置，请编辑 ~/.news4coder/config.json")
	}

	reqBody := ChatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.3,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.cfg.BaseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 LLM 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API 错误 (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("LLM API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM 返回空结果")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// SimpleChat 快速单轮对话
func (c *Client) SimpleChat(prompt string, maxTokens int) (string, error) {
	return c.Chat([]Message{
		{Role: "system", Content: "You are a helpful assistant for software engineers."},
		{Role: "user", Content: prompt},
	}, maxTokens)
}

// GetEmbedding 获取文本向量
func (c *Client) GetEmbedding(text string) ([]float64, error) {
	if c.cfg.APIKey == "" {
		return nil, fmt.Errorf("LLM API Key 未配置")
	}

	reqBody := EmbeddingRequest{
		Model: c.cfg.EmbeddingModel,
		Input: text,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.cfg.BaseURL+"/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 embedding 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API 错误 (status %d): %s", resp.StatusCode, string(body))
	}

	var embedResp EmbeddingResponse
	if err := json.Unmarshal(body, &embedResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if embedResp.Error != nil {
		return nil, fmt.Errorf("Embedding API 错误: %s", embedResp.Error.Message)
	}

	if len(embedResp.Data) == 0 {
		return nil, fmt.Errorf("Embedding 返回空结果")
	}

	return embedResp.Data[0].Embedding, nil
}

// CosineSimilarity 计算余弦相似度
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}


