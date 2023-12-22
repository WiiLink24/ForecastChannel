package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func sendWebhook() {
	if config.DiscordWebhook != "" {
		data := map[string]any{
			"username":   "Forecast Bot",
			"content":    "Weather Data has been updated!",
			"avatar_url": "https://rc24.xyz/images/logo-small.png",
			"attachments": []map[string]any{
				{
					"fallback":    "Weather Data Update",
					"color":       "#0381D7",
					"author_name": "RiiConnect24 Forecast Script",
					"author_icon": "https://rc24.xyz/images/webhooks/forecast/profile.png",
					"text":        "Weather Data has been updated!",
					"title":       "Update!",
					"fields": []map[string]any{
						{
							"title": "Script",
							"value": "Forecast Channel",
							"short": false,
						},
					},
					"thumb_url":   "https://rc24.xyz/images/webhooks/forecast/accuweather.png",
					"footer":      "RiiConnect24 Script",
					"footer_icon": "https://rc24.xyz/images/logo-small.png",
					"ts":          int(time.Now().Unix()),
				},
			},
		}

		_bytes, err := json.Marshal(data)
		checkError(err)

		_, err = http.Post(config.DiscordWebhook, "application/json", bytes.NewReader(_bytes))
		checkError(err)
	}
}
