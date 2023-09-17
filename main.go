package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/haru-24/go_chatgpt_test/model"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Could not read environment variables: %v", err)
	}
}

func inputKeybord() string {
	var inputMessage string
	print("[質問を入力してください。]:")
	fmt.Scan(&inputMessage)

	return inputMessage
}

func requestChatGpt(reqMessage string) {

	loadEnv()
	apiKey := os.Getenv("CHAT_GPT_APIKEY")

	message := []*model.RequestMessage{
		model.NewRequestMessage("user", reqMessage),
	}
	data, err := json.Marshal(model.NewRequest("gpt-3.5-turbo", message, 40))

	if err != nil {
		fmt.Println("failue json marshal", err)
		return
	}

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("create request error", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// send the request and retrieve the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("create request error", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read and parse response error:", err)
		return
	}

	// var response Response
	var response model.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("failuer json unmarshal", err)
		return
	}

	for _, r := range response.Choices {
		fmt.Printf("[%s]: %s\n", r.Message.Role, r.Message.Content)
	}

}

func main() {
	for {
		message := inputKeybord()
		requestChatGpt(message)
	}
}
