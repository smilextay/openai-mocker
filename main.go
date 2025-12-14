package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var stopReason string
var port int

func init() {
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()
}

func main() {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	stopReason = "stop"
	server := gin.Default()
	server.Use(CORS())
	server.POST("/openai/v1/chat/completions", OpenAIResponse)
	// Gemini API 原生路由
	// 使用通配符匹配 /v1/models/{model}:generateContent 或 /v1/models/{model}:streamGenerateContent
	server.POST("/v1/models/*path", GeminiResponse)
	server.POST("/v1beta/models/:model/generateContent", GeminiResponse)
	server.POST("/v1beta/models/:model/streamGenerateContent", GeminiResponse)
	log.Printf("Starting server on port %d", port)
	log.Fatal(server.Run(":" + strconv.Itoa(port)))
}
