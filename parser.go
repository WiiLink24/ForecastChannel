package main

import (
	"encoding/xml"
	"os"
)

type WeatherList struct {
	XMLName       xml.Name          `xml:"root"`
	National      []NationalList    `xml:"country"`
	International InternationalList `xml:"international"`
	Conditions    ConditionsList    `xml:"conditions"`
	Laundry       []Laundry         `xml:"laundry"`
	Wind          []Wind            `xml:"wind"`
	UV            []UV              `xml:"uv"`
	Pollen        []Pollen          `xml:"pollen"`
}

type NationalList struct {
	Name   LocalizedNames `xml:"name"`
	Cities []City         `xml:"city"`
}

type City struct {
	XMLName   xml.Name       `xml:"city"`
	Japanese  string         `xml:"jpn,attr"`
	English   string         `xml:"eng,attr"`
	German    string         `xml:"de,attr"`
	French    string         `xml:"fr,attr"`
	Spanish   string         `xml:"es,attr"`
	Italian   string         `xml:"it,attr"`
	Dutch     string         `xml:"nl,attr"`
	Province  LocalizedNames `xml:"province"`
	Longitude float64        `xml:"longitude"`
	Latitude  float64        `xml:"latitude"`
	Zoom1     int            `xml:"zoom1"`
	Zoom2     int            `xml:"zoom2"`
}

// LocalizedNames exists because I was too lazy to fix my XML for every single country
type LocalizedNames struct {
	Japanese string `xml:"jpn,attr"`
	English  string `xml:"eng,attr"`
	German   string `xml:"de,attr"`
	French   string `xml:"fr,attr"`
	Spanish  string `xml:"es,attr"`
	Italian  string `xml:"it,attr"`
	Dutch    string `xml:"nl,attr"`
}

type InternationalList struct {
	XMLName xml.Name            `xml:"international"`
	Cities  []InternationalCity `xml:"city"`
}

type InternationalCity struct {
	XMLName   xml.Name       `xml:"city"`
	Name      LocalizedNames `xml:"name"`
	Province  LocalizedNames `xml:"province"`
	Country   LocalizedNames `xml:"country"`
	Longitude float64        `xml:"longitude"`
	Latitude  float64        `xml:"latitude"`
	Zoom1     int            `xml:"zoom1"`
	Zoom2     int            `xml:"zoom2"`
}

type ConditionsList struct {
	XMLName    xml.Name    `xml:"conditions"`
	Conditions []Condition `xml:"condition"`
}

type Condition struct {
	Code          int            `xml:"code"`
	Name          LocalizedNames `xml:"name"`
	Code1         string         `xml:"code_1"`
	Code2         string         `xml:"code_2"`
	JapaneseCode1 string         `xml:"japanese_code_1"`
	JapaneseCode2 string         `xml:"japanese_code_2"`
}

type Laundry struct {
	XMLName xml.Name `xml:"laundry"`
	Name    string   `xml:"name"`
	Code    int      `xml:"code"`
}

type Wind struct {
	XMLName xml.Name `xml:"wind"`
	Name    string   `xml:"name"`
	Code    int      `xml:"code"`
}

type UV struct {
	XMLName xml.Name       `xml:"uv"`
	Name    LocalizedNames `xml:"name"`
}

type Pollen struct {
	XMLName xml.Name `xml:"pollen"`
	Name    string   `xml:"name"`
	Code    int      `xml:"code"`
}

func ParseWeatherXML() *WeatherList {
	var weather WeatherList
	data, err := os.ReadFile("weather.xml")
	checkError(err)

	err = xml.Unmarshal(data, &weather)
	checkError(err)

	return &weather
}
