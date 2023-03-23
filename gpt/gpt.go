package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qingconglaixueit/wechatbot/config"
	"github.com/qingconglaixueit/wechatbot/pkg/logger"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Chat Completions 响应体
type GPTChatCompletionResponseBody struct {
	ID      string                        `json:"id"`
	Object  string                        `json:"object"`
	Created int                           `json:"created"`
	Model   string                        `json:"model"`
	Choices []GPTChatCompletionChoiceItem `json:"choices"`
	Usage   map[string]interface{}        `json:"usage"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type GPTChatCompletionChoiceItem struct {
	Index        int     `json:"index"`
	Mesg         Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        uint    `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

// Chat Completions请求体
type GPTChatCompletionRequestBody struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        uint      `json:"max_tokens"`
	Temperature      float64   `json:"temperature"`
	TopP             int       `json:"top_p"`
	FrequencyPenalty int       `json:"frequency_penalty"`
	PresencePenalty  int       `json:"presence_penalty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ChatCompletions(msg string) (string, error) {
	cfg := config.LoadConfig()
	messages := make([]Message, 1)
	message := Message{
		Role:    "user",
		Content: msg,
	}
	messages[0] = message
	requestBody := GPTChatCompletionRequestBody{
		Model:            cfg.Model,
		Messages:         messages,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	return RequestOpenai(requestData, cfg.ChatUrl, cfg.ApiKey)
}

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	cfg := config.LoadConfig()
	requestBody := ChatGPTRequestBody{
		Model:            cfg.Model,
		Prompt:           msg,
		MaxTokens:        cfg.MaxTokens,
		Temperature:      cfg.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	return RequestOpenai(requestData, cfg.ChatUrl, cfg.ApiKey)
}

func RequestOpenai(requestData []byte, url string, apiKey string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", errors.New(fmt.Sprintf("请求GTP出错了，gpt api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body)))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))

	gptResponseBody := &GPTChatCompletionResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Mesg.Content
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	return reply, nil
}
