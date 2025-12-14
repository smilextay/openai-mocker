package main

type ChatRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Stream   bool                 `json:"stream"`
	Model    string               `json:"model"`
	Messages []ChatRequestMessage `json:"messages"`
}

type Completion struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionsStreamResponseChoiceDelta struct {
	Content string `json:"content"`
}

type ChatCompletionsStreamResponseChoice struct {
	Delta        ChatCompletionsStreamResponseChoiceDelta `json:"delta"`
	FinishReason *string                                  `json:"finish_reason"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created int64                                 `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}

// Gemini API 原生格式
type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason,omitempty"`
}

type GeminiGenerateContentResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiStreamChunk struct {
	Candidates []GeminiCandidate `json:"candidates"`
}
