package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Could not read environment variables: %v", err)
	}
}

func main() {

	loadEnv()
	apiKey := os.Getenv("CHAT_GPT_APIKEY")

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, strings.NewReader(`{
		"model":"gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "自己紹介をして下さい"}]}`))
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

	fmt.Println("body", string(body))

	// var response Response
	var response interface{}
	json.Unmarshal(body, &response)
	fmt.Println("response:", response)

}
