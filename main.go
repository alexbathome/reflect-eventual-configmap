package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type Config struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
	Enabled bool   `json:"enabled"`
}

var configPathFlag = flag.String("config", "config.json", "Path to the configuration file")

func main() {
	flag.Parse()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		var config Config
		f, err := os.Open(*configPathFlag)
		if err != nil {
			println("Failed to open config file:", err.Error())
		} else {
			jsonDecoder := json.NewDecoder(f)
			if err := jsonDecoder.Decode(&config); err != nil {
				println("Failed to decode config file:", err.Error())
			} else {
				// Use the config as needed
				println("Server will run on:", config.Address, ":", config.Port, "- Enabled:", config.Enabled)
			}
			f.Close()
		}

		<-ticker.C
	}
}
