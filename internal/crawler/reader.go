package crawler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// JinaReader 使用 Jina AI Reader 提取网页正文
type JinaReader struct {
	client *http.Client
}

// NewJinaReader 创建 Jina Reader
func NewJinaReader() *JinaReader {
	return &JinaReader{
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

// Fetch 读取指定 URL 的 Markdown 内容
func (r *JinaReader) Fetch(url string) (string, error) {
	jinaURL := "https://r.jina.ai/http://" + strings.TrimPrefix(url, "http://")
	if strings.HasPrefix(url, "https://") {
		jinaURL = "https://r.jina.ai/http://" + strings.TrimPrefix(url, "https://")
		// Actually jina supports https like this:
		jinaURL = "https://r.jina.ai/http://" + strings.TrimPrefix(url, "https://")
		// Wait, the correct format for jina reader is:
		// https://r.jina.ai/http://URL or https://r.jina.ai/http://URL
		// For https sites, it's actually https://r.jina.ai/http://URL (replacing https:// with empty then prefixing)
		// But more reliably: https://r.jina.ai/http://<url-without-scheme>
		// Let's use a simpler approach:
	}

	// Correct Jina AI Reader URL construction
	cleanURL := strings.TrimPrefix(url, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	jinaURL = "https://r.jina.ai/http://" + cleanURL

	req, err := http.NewRequest("GET", jinaURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Jina Reader 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取 Jina Reader 响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Jina Reader 返回状态码 %d", resp.StatusCode)
	}

	content := string(body)
	// Jina sometimes returns a title line starting with "Title: "
	// Keep it as is, enricher can handle it
	return content, nil
}

// FetchWithFallback 先尝试 Jina Reader，失败返回空字符串
func (r *JinaReader) FetchWithFallback(url string) string {
	content, err := r.Fetch(url)
	if err != nil {
		return ""
	}
	return content
}
