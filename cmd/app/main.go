package main

import (
	"flag"
	"log"
	"my_app/config"
	"my_app/internal/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config/config.yml", "Config path")
	flag.Parse()
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
