package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/haru-24/go_chatgpt_test/model"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Could not read environment variables: %v", err)
	}
}

func requestChatGpt(reqMessage string) string {

	loadEnv()
	apiKey := os.Getenv("CHAT_GPT_APIKEY")

	message := []*model.RequestMessage{
		model.NewRequestMessage("user", reqMessage+"(この質問に100文字以内で答えてください)"),
	}
	data, err := json.Marshal(model.NewRequest("gpt-3.5-turbo", message, 150))

	if err != nil {
		fmt.Println("failue json marshal", err)
		return string(err.Error())
	}

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("create request error", err)
		return string(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// send the request and retrieve the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("create request error", err)
		return string(err.Error())
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read and parse response error:", err)
		return string(err.Error())
	}

	// var response Response
	var response model.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("failuer json unmarshal", err)
		return string(err.Error())
	}

	var return_msg string
	for _, r := range response.Choices {
		fmt.Printf("[%s]: %s\n", r.Message.Role, r.Message.Content)
		return_msg = r.Message.Content
	}
	return return_msg
}

func convertDogLang(message string)string{
	var responceMessage string
	responceMessage = strings.Replace(message, "です。", "ガウ！", -1)
	responceMessage = strings.Replace(responceMessage, "ります。", "ガウ！", -1)
	responceMessage = strings.Replace(responceMessage, "でしょう。", "ガウ！", -1)
	responceMessage = strings.Replace(responceMessage, "る。", "ガウ！", -1)
	if !strings.Contains(responceMessage, "fg"){
		responceMessage = strings.Replace(responceMessage, "る。", "ガウ！", -1)
	}
	return responceMessage
}

func main() {
	loadEnv()

	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	router.POST("/callback", func(c *gin.Context) {

		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(400)
			} else {
				c.Writer.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					replayMessage := requestChatGpt(message.Text)
					convertDogLang := convertDogLang(replayMessage)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(convertDogLang)).Do(); err != nil {
						fmt.Println(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						fmt.Println(err)
					}
				}
			}
		}

	})

	router.Run(":3000")
}
