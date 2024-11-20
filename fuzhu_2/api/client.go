package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"fuzhu_2/types"
)

// APIClient AI API客户端
type APIClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewAPIClient 创建新的API客户端
func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		client:  &http.Client{},
		apiKey:  apiKey,
		baseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
	}
}

// ProcessText 处理单条文本
func (c *APIClient) ProcessText(input string) string {
	startTime := time.Now()
	log.Printf("开始处理输入文本: %s", truncateString(input, 50))

	requestBody := types.RequestBody{
		Model: "qwen-plus",
		Messages: []types.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: input},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("处理失败: %v", err)
		return ""
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ 请求创建失败: %v", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("❌ API请求失败: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ API返回非200状态码: %d", resp.StatusCode)
		return ""
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("响应读取失败: %v", err)
		return ""
	}

	var chatCompletion types.ChatCompletion
	err = json.Unmarshal(bodyText, &chatCompletion)
	if err != nil {
		log.Printf("JSON解析失败: %v", err)
		return ""
	}

	output := chatCompletion.Choices[0].Message.Content
	log.Printf("✅ 处理完成，耗时: %v, 输出长度: %d字符", time.Since(startTime), len(output))
	return output
}

// truncateString 辅助函数：如果字符串超过指定长度，截断并添加省略号
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
