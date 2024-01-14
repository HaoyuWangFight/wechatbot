package config

import (
	"encoding/json"
	"github.com/qingconglaixueit/wechatbot/pkg/logger"
	"log"
	"os"
	"time"
)

// Configuration 项目配置
type GeminiProConfiguration struct {
	ChatUrl string `json:"chat_url"`
	// gpt apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
	// 清空会话口令
	SessionClearToken string `json:"session_clear_token"`
	// 最大历史对话轮数
	MaxContentNum int `json:"max_content_num"`
	// 最大等待回复时间
	MaxWaitTime int64 `json:"max_wait_time"`
}

var conf *GeminiProConfiguration

func LoadGeminiProConfiguration() *GeminiProConfiguration {
	once.Do(func() {
		conf = &GeminiProConfiguration{
			ChatUrl:           " https://gemini.relationshit.win/v1beta/models/gemini-pro:generateContent",
			ApiKey:            "",
			AutoPass:          false,
			SessionTimeout:    60,
			SessionClearToken: "",
			MaxContentNum:     10,
			MaxWaitTime:       120,
		}

		// 判断配置文件是否存在，存在直接JSON读取
		_, err := os.Stat("gemini-pro-config.json")
		if err == nil {
			f, err := os.Open("gemini-pro-config.json")
			if err != nil {
				log.Fatalf("open gemini-pro-config.json err: %v", err)
				return
			}
			defer f.Close()
			encoder := json.NewDecoder(f)
			err = encoder.Decode(conf)
			if err != nil {
				log.Fatalf("decode gemini-pro-config err: %v", err)
				return
			}
		}
	})

	if conf.ApiKey == "" {
		logger.Danger("gemini-pro-config err: api key required")
	}

	return conf
}
