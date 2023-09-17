package model

type Request struct {
	Model     string            `json:"model"`
	Message   []*RequestMessage `json:"messages"`
	MaxTokens int               `json:"max_tokens"`
}

func NewRequest(modelID string, mesages []*RequestMessage, maxTokens int) *Request {
	return &Request{
		Model:     modelID,
		Message:   mesages,
		MaxTokens: maxTokens,
	}
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewRequestMessage(role string, content string) *RequestMessage {
	return &RequestMessage{
		Role:    role,
		Content: content,
	}
}
