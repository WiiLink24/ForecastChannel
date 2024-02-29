package main

import (
	"encoding/xml"
	"os"
)

type Config struct {
	AccuweatherKey    string `xml:"accuweather_key"`
	UseS3             bool   `xml:"use_s3"`
	S3BucketName      string `xml:"s3_bucket_name"`
	S3AccountID       string `xml:"s3_account_id"`
	S3ConnectionURL   string `xml:"s3_connection_url"`
	S3AccessIDKey     string `xml:"s3_access_key_id"`
	S3SecretAccessKey string `xml:"s3_secret_access_key"`
}

func GetConfig() Config {
	data, err := os.ReadFile("config.xml")
	checkError(err)

	var config Config
	err = xml.Unmarshal(data, &config)
	checkError(err)

	return config
}
