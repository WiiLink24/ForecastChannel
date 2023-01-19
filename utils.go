package main

import (
	"golang.org/x/exp/slices"
	"strconv"
)

// fixTime adjusts the timestamp to coincide with the Wii's UTC timestamp.
func fixTime(value int) uint32 {
	return uint32((value - 946684800) / 60)
}

func ConvertIcon(icon int, countryCode uint8) uint16 {
	code := "FFFF"
	for _, condition := range weatherList.Conditions.Conditions {
		if condition.Code == icon {
			if countryCode == 1 {
				code = condition.JapaneseCode1
			} else {
				code = condition.Code1
			}
		}
	}
	value, err := strconv.ParseInt(code, 16, 32)
	checkError(err)

	return uint16(value)
}

func GetWind(value string) uint8 {
	for _, wind := range weatherList.Wind {
		if wind.Name == value {
			return uint8(wind.Code)
		}
	}

	return 0xFF
}

func CoordinateEncode(value float64) int16 {
	value /= 0.0054931640625
	return int16(value)
}

func ZFill(value uint8, size int) string {
	str := strconv.FormatInt(int64(value), 10)
	temp := ""

	for i := 0; i < size-len(str); i++ {
		temp += "0"
	}

	return temp + str
}

func (f *Forecast) IsJapan() bool {
	return f.currentCountryCode == 1
}

func (f *Forecast) GetTemperatureFlag() uint8 {
	if f.currentCountryCode == 1 {
		return 0
	} else if slices.Contains([]uint8{8, 9, 12, 14, 17, 19, 37, 43, 48, 49, 51}, f.currentCountryCode) {
		return 1
	} else {
		return 2
	}
}

func GetSupportedLanguages(countryCode uint8) []uint8 {
	if countryCode == 1 {
		return []uint8{0}
	} else if 8 <= countryCode && countryCode <= 52 {
		return []uint8{1, 3, 4}
	} else if 64 <= countryCode && countryCode <= 110 {
		return []uint8{1, 2, 3, 4, 5, 6}
	}

	return []uint8{1}
}

func (f *Forecast) GetLocalizedName(names LocalizedNames) string {
	switch f.currentLanguageCode {
	case 0:
		return names.Japanese
	case 1:
		return names.English
	case 2:
		return names.German
	case 3:
		return names.French
	case 4:
		return names.Spanish
	case 5:
		return names.Italian
	case 6:
		return names.Dutch
	}

	// Impossible to reach here
	return ""
}

func (f *Forecast) GetCityName(city City) string {
	switch f.currentLanguageCode {
	case 0:
		return city.Japanese
	case 1:
		return city.English
	case 2:
		return city.German
	case 3:
		return city.French
	case 4:
		return city.Spanish
	case 5:
		return city.Italian
	case 6:
		return city.Dutch
	}

	// Impossible to reach here
	return ""
}
