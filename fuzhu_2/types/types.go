package types

// ChatCompletion 定义API响应的数据结构
type ChatCompletion struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"` // API返回的文本内容
		} `json:"message"`
	} `json:"choices"`
}

// Message 定义聊天消息的数据结构
type Message struct {
	Role    string `json:"role"`    // 消息角色（system/user）
	Content string `json:"content"` // 消息内容
}

// RequestBody 定义发送给API的请求体结构
type RequestBody struct {
	Model    string    `json:"model"`    // 使用的AI模型
	Messages []Message `json:"messages"` // 对话消息列表
}

// Result 定义处理结果的数据结构
type Result struct {
	RowIndex int    // Excel中的行索引
	Input    string // 输入文本
	Output   string // AI处理后的输出文本
}
