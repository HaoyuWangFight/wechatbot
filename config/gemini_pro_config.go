package config

import "time"

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
}

var conf *GeminiProConfiguration

func LoadGeminiProConfiguration() *GeminiProConfiguration {
	once.Do(func() {
		conf = &GeminiProConfiguration{}
	})
	return conf
}
