package main

import (
	"strconv"
	"unicode/utf16"
)

type WeatherConditionsTable struct {
	Code1      uint16
	Code2      uint16
	TextOffset uint32
}

func (f *Forecast) MakeWeatherConditionsTable() {
	f.Header.WeatherConditionTableOffset = f.GetCurrentSize()

	for _, condition := range weatherList.Conditions.Conditions {
		code1, err := strconv.ParseInt(condition.Code1, 16, 32)
		checkError(err)

		code2, err := strconv.ParseInt(condition.Code2, 16, 32)
		checkError(err)

		japaneseCode1, err := strconv.ParseInt(condition.JapaneseCode1, 16, 32)
		checkError(err)

		japaneseCode2, err := strconv.ParseInt(condition.JapaneseCode2, 16, 32)
		checkError(err)

		f.WeatherConditionsTable = append(f.WeatherConditionsTable, WeatherConditionsTable{
			Code1:      uint16(code1),
			Code2:      uint16(code2),
			TextOffset: 0,
		})

		f.WeatherConditionsTable = append(f.WeatherConditionsTable, WeatherConditionsTable{
			Code1:      uint16(japaneseCode1),
			Code2:      uint16(japaneseCode2),
			TextOffset: 0,
		})
	}

	f.Header.NumberOfWeatherConditionTables = uint32(len(f.WeatherConditionsTable))
}

func (f *Forecast) MakeWeatherConditionText() {
	i := 0
	for _, condition := range weatherList.Conditions.Conditions {
		f.WeatherConditionsTable[i].TextOffset = f.GetCurrentSize()
		f.WeatherConditionsText = append(f.WeatherConditionsText, utf16.Encode([]rune(f.GetLocalizedName(condition.Name)))...)
		f.WeatherConditionsText = append(f.WeatherConditionsText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.WeatherConditionsText = append(f.WeatherConditionsText, uint16(0))
		}

		f.WeatherConditionsTable[i+1].TextOffset = f.GetCurrentSize()
		f.WeatherConditionsText = append(f.WeatherConditionsText, utf16.Encode([]rune(f.GetLocalizedName(condition.Name)))...)
		f.WeatherConditionsText = append(f.WeatherConditionsText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.WeatherConditionsText = append(f.WeatherConditionsText, uint16(0))
		}
		i += 2
	}
}
