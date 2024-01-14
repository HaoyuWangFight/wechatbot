package gemini

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
	"net/url"
	"time"
)

type Text struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Text `json:"parts"`
	Role  string `json:"role"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type PromptFeedback struct {
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type GeminiProResponseBody struct {
	Candidates     []Candidate    `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback"`
}

type GeminiProRequestBody struct {
	Contents []Content `json:"contents"`
}

type Chatter struct {
	HistoryContents []Content
	LastRefreshTime int64
}

var Chatters map[string]*Chatter

func ChatCompletions(msg, user string) (string, error) {
	cfg := config.LoadGeminiProConfiguration()

	if _, ok := Chatters[user]; !ok {
		Chatters[user] = &Chatter{HistoryContents: make([]Content, 0)}
	}

	Chatters[user].HistoryContents = append(Chatters[user].HistoryContents, Content{
		Parts: []Text{{Text: msg}},
		Role:  "user",
	})
	Chatters[user].LastRefreshTime = time.Now().UnixMilli()
	if len(Chatters[user].HistoryContents) > cfg.MaxContentNum {
		Chatters[user].HistoryContents = Chatters[user].HistoryContents[len(
			Chatters[user].HistoryContents)-cfg.MaxContentNum:]
	}

	requestBody := GeminiProRequestBody{
		Contents: Chatters[user].HistoryContents,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gemini pro json string : %v", string(requestData)))
	return RequestGeminiPro(requestData, user)
}

func RequestGeminiPro(requestData []byte, user string) (string, error) {
	cfg := config.LoadGeminiProConfiguration()

	queryParams := url.Values{}
	queryParams.Set("key", cfg.ApiKey)
	req, err := http.NewRequest("POST", cfg.ChatUrl+"?"+queryParams.Encode(),
		bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", errors.New(fmt.Sprintf("请求Gemini pro出错了，gpt api status code not equals 200,code is %d ,"+
			"details:  %v ", response.StatusCode, string(body)))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gemini pro json string : %v", string(body)))

	responseBody := &GeminiProResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(responseBody.Candidates) > 0 && len(responseBody.Candidates[0].Content.Parts) > 0 {
		reply = responseBody.Candidates[0].Content.Parts[0].Text
		Chatters[user].HistoryContents = append(Chatters[user].HistoryContents,
			responseBody.Candidates[0].Content)
		Chatters[user].LastRefreshTime = time.Now().UnixMilli()
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	return reply, nil
}
