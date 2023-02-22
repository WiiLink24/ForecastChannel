package main

import (
	"bytes"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"time"
)

type ShortHeader struct {
	Version                       uint32
	Filesize                      uint32
	CRC32                         uint32
	OpenTimestamp                 uint32
	CloseTimestamp                uint32
	CountryCode                   uint8
	_                             [3]byte
	LanguageCode                  uint8
	TemperatureFlag               uint8
	_                             uint16
	NumberOfCurrentForecastTables uint32
	CurrentForecastTableOffset    uint32
}

type CurrentForecastTable struct {
	CountryCode                  uint8
	RegionCode                   uint8
	LocationCode                 uint16
	LocalTimestamp               uint32
	GlobalTimestamp              uint32
	CurrentForecast              uint16
	_                            uint8
	CurrentTemperatureCelsius    uint8
	CurrentTemperatureFahrenheit uint8
	CurrentWindDirection         uint8
	CurrentWindSpeedMetric       uint8
	CurrentWindSpeedImperial     uint8
	_                            uint16
	Unknown                      uint16
}

func (f *Forecast) MakeShortBin(cities []InternationalCity) ([]byte, []byte) {
	header := ShortHeader{
		Version:                       0,
		Filesize:                      0,
		CRC32:                         0,
		OpenTimestamp:                 fixTime(int(time.Now().Unix())),
		CloseTimestamp:                fixTime(int(time.Now().Unix())) + 63,
		CountryCode:                   f.currentCountryCode,
		LanguageCode:                  f.currentLanguageCode,
		TemperatureFlag:               0,
		NumberOfCurrentForecastTables: 0,
		CurrentForecastTableOffset:    36,
	}
	var currentForecastTables []CurrentForecastTable

	for _, city := range f.currentCountryList.Cities {
		weather := *weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)]
		countryCode := f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].CountryCode
		currentForecastTables = append(currentForecastTables, CurrentForecastTable{
			CountryCode:                  countryCode,
			RegionCode:                   f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].RegionCode,
			LocationCode:                 f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].LocationCode,
			LocalTimestamp:               fixTime(weather.Globe.Time),
			GlobalTimestamp:              fixTime(int(time.Now().Unix())),
			CurrentForecast:              ConvertIcon(weather.Current.WeatherIcon, countryCode),
			CurrentTemperatureCelsius:    uint8(weather.Current.TempCelsius),
			CurrentTemperatureFahrenheit: uint8(weather.Current.TempFahrenheit),
			CurrentWindDirection:         GetWind(weather.Current.WindDirection),
			CurrentWindSpeedMetric:       uint8(weather.Current.WindMetric),
			CurrentWindSpeedImperial:     uint8(weather.Current.WindImperial),
			Unknown:                      0xFFFF,
		})
	}

	for _, city := range cities {
		weather := *weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)]
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		countryCode := f.rawLocations[city.Country.English][city.Province.English][city.Name.English].CountryCode
		currentForecastTables = append(currentForecastTables, CurrentForecastTable{
			CountryCode:                  countryCode,
			RegionCode:                   f.rawLocations[city.Country.English][city.Province.English][city.Name.English].RegionCode,
			LocationCode:                 f.rawLocations[city.Country.English][city.Province.English][city.Name.English].LocationCode,
			LocalTimestamp:               fixTime(weather.Globe.Time),
			GlobalTimestamp:              fixTime(int(time.Now().Unix())),
			CurrentForecast:              ConvertIcon(weather.Current.WeatherIcon, countryCode),
			CurrentTemperatureCelsius:    uint8(weather.Current.TempCelsius),
			CurrentTemperatureFahrenheit: uint8(weather.Current.TempFahrenheit),
			CurrentWindDirection:         GetWind(weather.Current.WindDirection),
			CurrentWindSpeedMetric:       uint8(weather.Current.WindMetric),
			CurrentWindSpeedImperial:     uint8(weather.Current.WindImperial),
			Unknown:                      0xFFFF,
		})
	}

	header.NumberOfCurrentForecastTables = uint32(len(currentForecastTables))

	buffer := new(bytes.Buffer)
	Write(buffer, header)
	Write(buffer, currentForecastTables)

	header.Filesize = uint32(buffer.Len())
	buffer.Reset()
	Write(buffer, header)
	Write(buffer, currentForecastTables)

	crcTable := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
	header.CRC32 = checksum

	buffer.Reset()
	Write(buffer, header)
	Write(buffer, currentForecastTables)

	compressed, err := lz10.Compress(buffer.Bytes())
	checkError(err)

	// Prepare to make for Wii U
	// TODO: Patch IOS to force proper UTC
	wiiUBuffer := new(bytes.Buffer)
	header.OpenTimestamp = 0
	header.CloseTimestamp = 0xFFFFFFFF
	Write(wiiUBuffer, header)
	Write(wiiUBuffer, currentForecastTables)

	crcTable = crc32.MakeTable(crc32.IEEE)
	checksum = crc32.Checksum(wiiUBuffer.Bytes()[12:], crcTable)
	header.CRC32 = checksum

	wiiUBuffer.Reset()
	Write(wiiUBuffer, header)
	Write(wiiUBuffer, currentForecastTables)

	compressedU, err := lz10.Compress(wiiUBuffer.Bytes())
	checkError(err)

	return compressed, compressedU
}
