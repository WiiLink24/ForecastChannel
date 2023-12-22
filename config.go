package main

import (
	"encoding/xml"
	"os"
)

type Config struct {
	ForecastHost       string `xml:"forecast_host"`
	AccuweatherKey     string `xml:"accuweather_key"`
	CloudflareToken    string `xml:"cloudflare_token"`
	CloudflareZoneName string `xml:"cloudflare_zone_name"`
	DiscordWebhook     string `xml:"discord_webhook"`
}

func GetConfig() Config {
	data, err := os.ReadFile("config.xml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = xml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}
