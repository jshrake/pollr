package pollr

import (
	"encoding/json"
	"os"
)

type Databse struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Redis struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type Config struct {
	WebAddress  string  `json:"web_address"`
	RestAddress string  `json:"rest_address"`
	WsAddress   string  `json:"ws_address"`
	AppSecret   string  `json:"app_secret"`
	Database    Databse `json:"database"`
	Redis       Redis   `json:"redis"`
}

func NewConfig(configFile string) *Config {
	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	config := &Config{}
	decoder.Decode(&config)
	return config
}
