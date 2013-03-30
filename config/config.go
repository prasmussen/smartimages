package config

import (
    "os"
    "encoding/json"
)

const (
    DefaultConfig = "config.json"
)

type Config struct {
    Listen string
    LogFile string
    ImageDir string
}

func Load() (*Config, error) {
    // Open config for reading
    f, err := os.Open(DefaultConfig)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // Unmarshal config data
    cfg := &Config{}
    if err := json.NewDecoder(f).Decode(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}
