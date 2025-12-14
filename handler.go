package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func OpenAIResponse(c *gin.Context) {
	var chatRequest ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prompt := "This is a mock server."
	if len(chatRequest.Messages) != 0 {
		prompt = chatRequest.Messages[len(chatRequest.Messages)-1].Content
	}
	response := prompt2response(prompt)

	if chatRequest.Stream {
		setEventStreamHeaders(c)
		dataChan := make(chan string)
		stopChan := make(chan bool)
		streamResponse := ChatCompletionsStreamResponse{
			Id:      "chat_cmp_l-i_will_always_love_you",
			Object:  "chat.completion.chunk",
			Created: 1689411338,
			Model:   "gpt-3.5-turbo",
		}
		streamResponseChoice := ChatCompletionsStreamResponseChoice{}
		go func() {
			for i, s := range response {
				streamResponseChoice.Delta.Content = string(s)
				if i == len(response)-1 {
					streamResponseChoice.FinishReason = &stopReason
				}
				streamResponse.Choices = []ChatCompletionsStreamResponseChoice{streamResponseChoice}
				jsonStr, _ := json.Marshal(streamResponse)
				dataChan <- string(jsonStr)
			}
			stopChan <- true
		}()

		c.Stream(func(w io.Writer) bool {
			select {
			case data := <-dataChan:
				c.Render(-1, CustomEvent{Data: "data: " + data})
				return true
			case <-stopChan:
				c.Render(-1, CustomEvent{Data: "data: [DONE]"})
				return false
			}
		})
	} else {
		c.JSON(http.StatusOK, Completion{
			Id:      "chat_cmp_l-7f8Qxn9XkoGsVcl0RVGltZpPeqMAG",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "gpt-3.5-turbo",
			Choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    "assistant",
						Content: prompt,
					},
					FinishReason: "length",
				},
			},
			Usage: Usage{
				PromptTokens:     9,
				CompletionTokens: 1,
				TotalTokens:      10,
			},
		})
	}
}

func GeminiResponse(c *gin.Context) {
	var geminiRequest GeminiRequest
	if err := c.ShouldBindJSON(&geminiRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户输入
	prompt := "This is a mock Gemini server."
	if len(geminiRequest.Contents) != 0 {
		lastContent := geminiRequest.Contents[len(geminiRequest.Contents)-1]
		if len(lastContent.Parts) != 0 {
			prompt = lastContent.Parts[len(lastContent.Parts)-1].Text
		}
	}
	response := prompt2response(prompt)

	// 从 URL 路径中获取模型名称和操作类型
	modelName := c.Param("model")
	path := c.Param("path")

	// 处理 /v1/models/{model}:generateContent 格式
	if path != "" {
		// path 格式为 /{model}:generateContent 或 /{model}:streamGenerateContent
		parts := strings.Split(strings.TrimPrefix(path, "/"), ":")
		if len(parts) > 0 {
			modelName = parts[0]
		}
	}

	if modelName == "" {
		modelName = "gemini-pro"
	}

	// 检查是否是流式请求
	// 通过路径判断：包含 streamGenerateContent 则为流式
	stream := strings.Contains(c.Request.URL.Path, "streamGenerateContent") ||
		c.Query("stream") == "true" ||
		c.GetHeader("Accept") == "text/event-stream"

	if stream {
		setEventStreamHeaders(c)
		dataChan := make(chan string)
		stopChan := make(chan bool)
		go func() {
			for i, s := range response {
				candidate := GeminiCandidate{
					Content: GeminiContent{
						Role: "model",
						Parts: []GeminiPart{
							{Text: string(s)},
						},
					},
				}
				if i == len(response)-1 {
					finishReason := "STOP"
					candidate.FinishReason = finishReason
				}
				chunk := GeminiStreamChunk{
					Candidates: []GeminiCandidate{candidate},
				}
				jsonStr, _ := json.Marshal(chunk)
				dataChan <- string(jsonStr)
			}
			stopChan <- true
		}()

		c.Stream(func(w io.Writer) bool {
			select {
			case data := <-dataChan:
				c.Render(-1, CustomEvent{Data: "data: " + data})
				return true
			case <-stopChan:
				c.Render(-1, CustomEvent{Data: "data: [DONE]"})
				return false
			}
		})
	} else {
		c.JSON(http.StatusOK, GeminiGenerateContentResponse{
			Candidates: []GeminiCandidate{
				{
					Content: GeminiContent{
						Role: "model",
						Parts: []GeminiPart{
							{Text: response},
						},
					},
					FinishReason: "STOP",
				},
			},
		})
	}
}
