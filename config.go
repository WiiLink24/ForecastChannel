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
	UseCloudflare      bool   `xml:"use_cloudflare"`
}

func GetConfig() Config {
	data, err := os.ReadFile("config.xml")
	checkError(err)

	var config Config
	err = xml.Unmarshal(data, &config)
	checkError(err)

	return config
}
