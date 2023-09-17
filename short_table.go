package main

import (
	"fmt"
)

type ShortForecastTable struct {
	CountryCode                         uint8
	RegionCode                          uint8
	LocationCode                        uint16
	LocalTimestamp                      uint32
	GlobalTimestamp                     uint32
	_                                   uint32
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
	Unknown                             [3]byte
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
	TodayUVIndex                        uint8
	TodayLaundryIndex                   uint8
	TodayPollenCount                    uint8
}

func (f *Forecast) MakeShortForecastTable(cities []InternationalCity) {
	f.Header.ShortForecastTableOffset = f.GetCurrentSize()

	for _, city := range cities {
		weather := *weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)]
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		country, _ := f.rawLocations.Get(city.Country.English)
		province, _ := country.Get(city.Province.English)
		currentCity, _ := province.Get(city.Name.English)
		countryCode := currentCity.CountryCode
		f.ShortForecastTable = append(f.ShortForecastTable, ShortForecastTable{
			CountryCode:                         countryCode,
			RegionCode:                          currentCity.RegionCode,
			LocationCode:                        currentCity.LocationCode,
			LocalTimestamp:                      fixTime(weather.Globe.Time),
			GlobalTimestamp:                     fixTime(int(currentTime)),
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
			Unknown:                             [3]byte{0xFF, 0xFF, 0xFF},
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
			TodayUVIndex:                        uint8(weather.UVIndex),
			TodayLaundryIndex:                   255,
			TodayPollenCount:                    uint8(weather.Pollen),
		})
	}

	f.Header.NumberOfShortForecastTables = uint32(len(f.ShortForecastTable))
}
