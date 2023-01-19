package main

import (
	"fmt"
	"time"
)

type LongForecastTable struct {
	CountryCode                         uint8
	RegionCode                          uint8
	LocationCode                        uint16
	LocalTimestamp                      uint32
	GlobalTimestamp                     uint32
	Unknown                             uint32
	TodayForecast                       uint16
	Today6Hour12AMTo6AM                 uint16
	Today6Hour6AMTo12PM                 uint16
	Today6Hour12PMTo6PM                 uint16
	Today6Hour6PMTo12AM                 uint16
	TodayHighCelsius                    int8
	TodayHighDifferenceCelsius          int8
	TodayLowCelsius                     int8
	TodayLowDifferenceCelsius           int8
	TodayHighFahrenheit                 int8
	TodayHighDifferenceFahrenheit       int8
	TodayLowFahrenheit                  int8
	TodayLowDifferenceFahrenheit        int8
	Today6HourPrecipitation12AMTo6AM    uint8
	Today6HourPrecipitation6AMTo12PM    uint8
	Today6HourPrecipitation12PMTo6PM    uint8
	Today6HourPrecipitation6PMTo12AM    uint8
	TodayWindDirection                  uint8
	TodayWindSpeedMetric                uint8
	TodayWindSpeedImperial              uint8
	TodayUVIndex                        uint8
	TodayLaundryIndex                   uint8
	TodayPollenCount                    uint8
	TomorrowForecast                    uint16
	Tomorrow6Hour12AMTo6AM              uint16
	Tomorrow6Hour6AMTo12PM              uint16
	Tomorrow6Hour12PMTo6PM              uint16
	Tomorrow6Hour6PMTo12AM              uint16
	TomorrowHighCelsius                 int8
	TomorrowHighDifferenceCelsius       int8
	TomorrowLowCelsius                  int8
	TomorrowLowDifferenceCelsius        int8
	TomorrowHighFahrenheit              int8
	TomorrowHighDifferenceFahrenheit    int8
	TomorrowLowFahrenheit               int8
	TomorrowLowDifferenceFahrenheit     int8
	Tomorrow6HourPrecipitation12AMTo6AM uint8
	Tomorrow6HourPrecipitation6AMTo12PM uint8
	Tomorrow6HourPrecipitation12PMTo6PM uint8
	Tomorrow6HourPrecipitation6PMTo12AM uint8
	TomorrowWindDirection               uint8
	TomorrowWindSpeedMetric             uint8
	TomorrowWindSpeedImperial           uint8
	TomorrowUVIndex                     uint8
	TomorrowLaundryIndex                uint8
	TomorrowPollenCount                 uint8
	FiveDayForecastDay1                 uint16
	FiveDayForecastDay1HighCelsius      int8
	FiveDayForecastDay1LowCelsius       int8
	FiveDayForecastDay1HighFahrenheit   int8
	FiveDayForecastDay1LowFahrenheit    int8
	FiveDayForecastDay1Precipitation    int8
	_                                   uint8
	FiveDayForecastDay2                 uint16
	FiveDayForecastDay2HighCelsius      int8
	FiveDayForecastDay2LowCelsius       int8
	FiveDayForecastDay2HighFahrenheit   int8
	FiveDayForecastDay2LowFahrenheit    int8
	FiveDayForecastDay2Precipitation    int8
	_                                   uint8
	FiveDayForecastDay3                 uint16
	FiveDayForecastDay3HighCelsius      int8
	FiveDayForecastDay3LowCelsius       int8
	FiveDayForecastDay3HighFahrenheit   int8
	FiveDayForecastDay3LowFahrenheit    int8
	FiveDayForecastDay3Precipitation    int8
	_                                   uint8
	FiveDayForecastDay4                 uint16
	FiveDayForecastDay4HighCelsius      int8
	FiveDayForecastDay4LowCelsius       int8
	FiveDayForecastDay4HighFahrenheit   int8
	FiveDayForecastDay4LowFahrenheit    int8
	FiveDayForecastDay4Precipitation    int8
	_                                   uint8
	FiveDayForecastDay5                 uint16
	FiveDayForecastDay5HighCelsius      int8
	FiveDayForecastDay5LowCelsius       int8
	FiveDayForecastDay5HighFahrenheit   int8
	FiveDayForecastDay5LowFahrenheit    int8
	FiveDayForecastDay5Precipitation    int8
	_                                   uint8
	FiveDayForecastDay6                 uint16
	FiveDayForecastDay6HighCelsius      int8
	FiveDayForecastDay6LowCelsius       int8
	FiveDayForecastDay6HighFahrenheit   int8
	FiveDayForecastDay6LowFahrenheit    int8
	FiveDayForecastDay6Precipitation    int8
	_                                   uint8
	FiveDayForecastDay7                 uint16
	FiveDayForecastDay7HighCelsius      int8
	FiveDayForecastDay7LowCelsius       int8
	FiveDayForecastDay7HighFahrenheit   int8
	FiveDayForecastDay7LowFahrenheit    int8
	FiveDayForecastDay7Precipitation    int8
	_                                   uint8
}

func (f *Forecast) MakeLongForecastTable() {
	f.Header.LongForecastTableOffset = f.GetCurrentSize()

	for _, city := range f.currentCountryList.Cities {
		weather := *weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)]
		countryCode := f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].CountryCode
		f.LongForecastTable = append(f.LongForecastTable, LongForecastTable{
			CountryCode:                         countryCode,
			RegionCode:                          f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].RegionCode,
			LocationCode:                        f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].LocationCode,
			LocalTimestamp:                      fixTime(weather.Globe.Time),
			GlobalTimestamp:                     fixTime(int(time.Now().Unix())),
			TodayForecast:                       ConvertIcon(weather.Today.WeatherIcon, countryCode),
			Today6Hour12AMTo6AM:                 ConvertIcon(weather.HourlyIcon[0], countryCode),
			Today6Hour6AMTo12PM:                 ConvertIcon(weather.HourlyIcon[1], countryCode),
			Today6Hour12PMTo6PM:                 ConvertIcon(weather.HourlyIcon[2], countryCode),
			Today6Hour6PMTo12AM:                 ConvertIcon(weather.HourlyIcon[3], countryCode),
			TodayHighCelsius:                    int8(weather.Today.TempCelsiusMax),
			TodayHighDifferenceCelsius:          -128,
			TodayLowCelsius:                     int8(weather.Today.TempCelsiusMin),
			TodayLowDifferenceCelsius:           -128,
			TodayHighFahrenheit:                 int8(weather.Today.TempFahrenheitMax),
			TodayHighDifferenceFahrenheit:       -128,
			TodayLowFahrenheit:                  int8(weather.Today.TempFahrenheitMin),
			TodayLowDifferenceFahrenheit:        -128,
			Today6HourPrecipitation12AMTo6AM:    uint8(weather.Precipitation[0]),
			Today6HourPrecipitation6AMTo12PM:    uint8(weather.Precipitation[1]),
			Today6HourPrecipitation12PMTo6PM:    uint8(weather.Precipitation[2]),
			Today6HourPrecipitation6PMTo12AM:    uint8(weather.Precipitation[3]),
			TodayWindDirection:                  GetWind(weather.Wind.WindDirection),
			TodayWindSpeedMetric:                uint8(weather.Wind.WindMetric),
			TodayWindSpeedImperial:              uint8(weather.Wind.WindImperial),
			TodayUVIndex:                        uint8(weather.UVIndex),
			TodayLaundryIndex:                   231,
			TodayPollenCount:                    uint8(weather.Pollen),
			TomorrowForecast:                    ConvertIcon(weather.Tomorrow.WeatherIcon, countryCode),
			Tomorrow6Hour12AMTo6AM:              ConvertIcon(weather.HourlyIcon[4], countryCode),
			Tomorrow6Hour6AMTo12PM:              ConvertIcon(weather.HourlyIcon[5], countryCode),
			Tomorrow6Hour12PMTo6PM:              ConvertIcon(weather.HourlyIcon[6], countryCode),
			Tomorrow6Hour6PMTo12AM:              ConvertIcon(weather.HourlyIcon[7], countryCode),
			TomorrowHighCelsius:                 int8(weather.Tomorrow.TempCelsiusMax),
			TomorrowHighDifferenceCelsius:       -128,
			TomorrowLowCelsius:                  int8(weather.Tomorrow.TempCelsiusMin),
			TomorrowLowDifferenceCelsius:        -128,
			TomorrowHighFahrenheit:              int8(weather.Tomorrow.TempFahrenheitMax),
			TomorrowHighDifferenceFahrenheit:    -128,
			TomorrowLowFahrenheit:               int8(weather.Tomorrow.TempFahrenheitMin),
			TomorrowLowDifferenceFahrenheit:     -128,
			Tomorrow6HourPrecipitation12AMTo6AM: uint8(weather.Precipitation[4]),
			Tomorrow6HourPrecipitation6AMTo12PM: uint8(weather.Precipitation[5]),
			Tomorrow6HourPrecipitation12PMTo6PM: uint8(weather.Precipitation[6]),
			Tomorrow6HourPrecipitation6PMTo12AM: uint8(weather.Precipitation[7]),
			TomorrowWindDirection:               GetWind(weather.Wind.WindDirectionTomorrow),
			TomorrowWindSpeedMetric:             uint8(weather.Wind.WindMetricTomorrow),
			TomorrowWindSpeedImperial:           uint8(weather.Wind.WindImperialTomorrow),
			TomorrowUVIndex:                     uint8(weather.UVIndex),
			TomorrowLaundryIndex:                231,
			TomorrowPollenCount:                 uint8(weather.Pollen),
			FiveDayForecastDay1:                 ConvertIcon(weather.Week[0].WeatherIcon, countryCode),
			FiveDayForecastDay1HighCelsius:      int8(weather.Week[0].TempCelsiusMax),
			FiveDayForecastDay1LowCelsius:       int8(weather.Week[0].TempCelsiusMin),
			FiveDayForecastDay1HighFahrenheit:   int8(weather.Week[0].TempFahrenheitMax),
			FiveDayForecastDay1LowFahrenheit:    int8(weather.Week[0].TempFahrenheitMin),
			FiveDayForecastDay1Precipitation:    int8(weather.Precipitation[8]),
			FiveDayForecastDay2:                 ConvertIcon(weather.Week[1].WeatherIcon, countryCode),
			FiveDayForecastDay2HighCelsius:      int8(weather.Week[1].TempCelsiusMax),
			FiveDayForecastDay2LowCelsius:       int8(weather.Week[1].TempCelsiusMin),
			FiveDayForecastDay2HighFahrenheit:   int8(weather.Week[1].TempFahrenheitMax),
			FiveDayForecastDay2LowFahrenheit:    int8(weather.Week[1].TempFahrenheitMin),
			FiveDayForecastDay2Precipitation:    int8(weather.Precipitation[9]),
			FiveDayForecastDay3:                 ConvertIcon(weather.Week[2].WeatherIcon, countryCode),
			FiveDayForecastDay3HighCelsius:      int8(weather.Week[2].TempCelsiusMax),
			FiveDayForecastDay3LowCelsius:       int8(weather.Week[2].TempCelsiusMin),
			FiveDayForecastDay3HighFahrenheit:   int8(weather.Week[2].TempFahrenheitMax),
			FiveDayForecastDay3LowFahrenheit:    int8(weather.Week[2].TempFahrenheitMin),
			FiveDayForecastDay3Precipitation:    int8(weather.Precipitation[10]),
			FiveDayForecastDay4:                 ConvertIcon(weather.Week[3].WeatherIcon, countryCode),
			FiveDayForecastDay4HighCelsius:      int8(weather.Week[3].TempCelsiusMax),
			FiveDayForecastDay4LowCelsius:       int8(weather.Week[3].TempCelsiusMin),
			FiveDayForecastDay4HighFahrenheit:   int8(weather.Week[3].TempFahrenheitMax),
			FiveDayForecastDay4LowFahrenheit:    int8(weather.Week[3].TempFahrenheitMin),
			FiveDayForecastDay4Precipitation:    int8(weather.Precipitation[11]),
			FiveDayForecastDay5:                 ConvertIcon(weather.Week[4].WeatherIcon, countryCode),
			FiveDayForecastDay5HighCelsius:      int8(weather.Week[4].TempCelsiusMax),
			FiveDayForecastDay5LowCelsius:       int8(weather.Week[4].TempCelsiusMin),
			FiveDayForecastDay5HighFahrenheit:   int8(weather.Week[4].TempFahrenheitMax),
			FiveDayForecastDay5LowFahrenheit:    int8(weather.Week[4].TempFahrenheitMin),
			FiveDayForecastDay5Precipitation:    int8(weather.Precipitation[12]),
			FiveDayForecastDay6:                 ConvertIcon(weather.Week[5].WeatherIcon, countryCode),
			FiveDayForecastDay6HighCelsius:      int8(weather.Week[5].TempCelsiusMax),
			FiveDayForecastDay6LowCelsius:       int8(weather.Week[5].TempCelsiusMin),
			FiveDayForecastDay6HighFahrenheit:   int8(weather.Week[5].TempFahrenheitMax),
			FiveDayForecastDay6LowFahrenheit:    int8(weather.Week[5].TempFahrenheitMin),
			FiveDayForecastDay6Precipitation:    int8(weather.Precipitation[13]),
			FiveDayForecastDay7:                 ConvertIcon(weather.Week[6].WeatherIcon, countryCode),
			FiveDayForecastDay7HighCelsius:      int8(weather.Week[6].TempCelsiusMax),
			FiveDayForecastDay7LowCelsius:       int8(weather.Week[6].TempCelsiusMin),
			FiveDayForecastDay7HighFahrenheit:   int8(weather.Week[6].TempFahrenheitMax),
			FiveDayForecastDay7LowFahrenheit:    int8(weather.Week[6].TempFahrenheitMin),
			FiveDayForecastDay7Precipitation:    int8(weather.Precipitation[14]),
		})
	}

	f.Header.NumberOfLongForecastTables = uint32(len(f.LongForecastTable))
}
