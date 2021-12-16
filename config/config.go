package config

type Config struct {
	AppId     string `json:"appId"`     // 开发者ID
	AppSecret string `json:"appSecret"` // 开发者密钥
	Sandbox   bool   `json:"sandbox"`   // 沙盒
}
