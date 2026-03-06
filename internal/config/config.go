package config

import (
	"encoding/json"
	"os"
)

type Selectors struct {
	BuyButton      string `json:"buy_button"`
	CheckoutButton string `json:"checkout_button"`
}

type Config struct {
	ProductURL     string    `json:"product_url"`
	WorkerCount    int       `json:"worker_count"`
	Headless       bool      `json:"headless"`
	Debug          bool      `json:"debug"`
	TimeoutSeconds int       `json:"timeout_seconds"`
	Selectors      Selectors `json:"selectors"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
