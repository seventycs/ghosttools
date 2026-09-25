package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    BotToken   string `json:"bot_token"`
    ServerID   string `json:"server_id"`
    CategoryID string `json:"category_id"`
    WebhookURL string `json:"webhook_url"`
}

var Cfg Config

func Load() {
    f, err := os.Open("config.json")
    if err != nil {
        return
    }
    defer f.Close()
    json.NewDecoder(f).Decode(&Cfg)
}

func Save() {
    f, _ := os.Create("config.json")
    defer f.Close()
    json.NewEncoder(f).Encode(Cfg)
}