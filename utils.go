package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/exp/slices"
	"os"
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
	// Fill out the Russian supported countries first as those are specific
	if countryCode == 18 {
		return []uint8{7}
	}

	if countryCode == 1 {
		return []uint8{0, 1, 2, 3, 4, 5, 6}
	} else if 8 <= countryCode && countryCode <= 52 {
		return []uint8{1, 3, 4}
	} else if 64 <= countryCode && countryCode <= 110 {
		return []uint8{1, 2, 3, 4, 5, 6}
	}

	return []uint8{0, 1, 2, 3, 4, 5, 6}
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
	case 7:
		return names.Russian
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
	case 7:
		return city.Russian
	}

	// Impossible to reach here
	return ""
}

func SignFile(contents []byte) []byte {
	buffer := new(bytes.Buffer)

	// Get RSA key and sign
	rsaData, err := os.ReadFile("Private.pem")
	if err != nil {
		if !os.IsNotExist(err) {
			checkError(err)
		}

		// Otherwise the file does not exist. Assume this is GitHub Actions and return an empty signature.
		buffer.Write(make([]byte, 64))
		buffer.Write(make([]byte, 256))
		buffer.Write(contents)

		return buffer.Bytes()
	}

	rsaBlock, _ := pem.Decode(rsaData)

	parsedKey, err := x509.ParsePKCS8PrivateKey(rsaBlock.Bytes)
	checkError(err)

	// Hash our data then sign
	hash := sha1.New()
	_, err = hash.Write(contents)
	checkError(err)

	contentsHashSum := hash.Sum(nil)

	reader := rand.Reader
	signature, err := rsa.SignPKCS1v15(reader, parsedKey.(*rsa.PrivateKey), crypto.SHA1, contentsHashSum)
	checkError(err)

	buffer.Write(make([]byte, 64))
	buffer.Write(signature)
	buffer.Write(contents)

	return buffer.Bytes()
}
