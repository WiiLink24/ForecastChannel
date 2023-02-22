package accuweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	apiKey = ""
	apiURL = "https://api.accuweather.com"
)

var currentTime = time.Now().Unix()

func GetWeather(longitude float64, latitude float64, _time int64) *Weather {
	currentTime = _time
	weather := Weather{}

	// First retrieve the location code.
	queryParams := fmt.Sprintf("q=%f,%f&apikey=%s", latitude, longitude, apiKey)

	response, err := http.Get(fmt.Sprintf("%s/locations/v1/cities/geoposition/search.json?%s", apiURL, queryParams))
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		panic(err)
	}

	locationKey := jsonData["Key"].(string)
	weather.GetCurrentWeather(locationKey)
	weather.Get5DayWeather(locationKey)
	weather.Get10DayWeather(locationKey)

	return &weather
}

func (w *Weather) GetCurrentWeather(locationKey string) {
	response, err := http.Get(fmt.Sprintf("%s/currentconditions/v1/%s?apikey=%s&details=true", apiURL, locationKey, apiKey))
	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	jsonData := []any{map[string]any{}}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		panic(err)
	}

	weather := jsonData[0].(map[string]any)
	w.LocalTime = weather["LocalObservationDateTime"].(string)
	w.Current.TempFahrenheit = weather["Temperature"].(map[string]any)["Imperial"].(map[string]any)["Value"].(float64)
	w.Current.TempCelsius = weather["Temperature"].(map[string]any)["Metric"].(map[string]any)["Value"].(float64)
	w.Current.WeatherIcon = int(weather["WeatherIcon"].(float64))
	w.Current.WindDirection = weather["Wind"].(map[string]any)["Direction"].(map[string]any)["English"].(string)
	w.Current.WindImperial = weather["Wind"].(map[string]any)["Speed"].(map[string]any)["Imperial"].(map[string]any)["Value"].(float64)
	w.Current.WindMetric = weather["Wind"].(map[string]any)["Speed"].(map[string]any)["Metric"].(map[string]any)["Value"].(float64)

	one, err := strconv.ParseFloat(w.LocalTime[20:22], 32)
	if err != nil {
		panic(err)
	}

	two, err := strconv.ParseInt(w.LocalTime[23:25], 10, 32)
	if err != nil {
		panic(err)
	}

	w.Globe.Offset = int(one + float64(two/60))
	if string(w.LocalTime[19]) == "-" {
		w.Globe.Offset *= -1
	}

	w.Globe.Time = int(currentTime) + w.Globe.Offset*3600
}

func (w *Weather) Get5DayWeather(locationKey string) {
	response, err := http.Get(fmt.Sprintf("%s/forecasts/v1/daily/5day/quarters/%s?apikey=%s&details=true", apiURL, locationKey, apiKey))
	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	jsonData := []any{map[string]any{}}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		panic(err)
	}

	index := 0
	hourlyStart := 0
	isRightDay := false
	for !isRightDay {
		temp, err := strconv.ParseInt(jsonData[index].(map[string]any)["EffectiveDate"].(string)[11:13], 10, 32)
		if err != nil {
			panic(err)
		}

		hourlyStart = int(temp) / 6
		if jsonData[index].(map[string]any)["EffectiveDate"].(string)[:10] == w.LocalTime[:10] {
			isRightDay = true
		} else {
			index++
		}
	}

	// Make precipitation and icons together to save time
	i := 0
	w.Precipitation = []int{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	w.HourlyIcon = make([]int, 8)
	for _i := hourlyStart; _i < 8; _i++ {
		w.Precipitation[_i] = int(jsonData[index+i].(map[string]any)["PrecipitationProbability"].(float64))
		w.HourlyIcon[_i] = int(jsonData[index+i].(map[string]any)["Icon"].(float64))
		i++
	}
}

func (w *Weather) Get10DayWeather(locationKey string) {
	response, err := http.Get(fmt.Sprintf("%s/forecasts/v1/daily/10day/%s?apikey=%s&details=true", apiURL, locationKey, apiKey))
	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		panic(err)
	}

	day := 0
	if jsonData["DailyForecasts"].([]any)[0].(map[string]any)["Date"].(string)[:10] != w.LocalTime[:10] {
		day++
	}

	// Today Forecast
	w.Today.TempFahrenheitMin = jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Temperature"].(map[string]any)["Minimum"].(map[string]any)["Value"].(float64)
	w.Today.TempFahrenheitMax = jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Temperature"].(map[string]any)["Maximum"].(map[string]any)["Value"].(float64)
	w.Today.TempCelsiusMin = ftoC(w.Today.TempFahrenheitMin)
	w.Today.TempCelsiusMax = ftoC(w.Today.TempFahrenheitMax)
	w.Today.WeatherIcon = int(jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Day"].(map[string]any)["Icon"].(float64))

	// Tomorrow Forecast
	w.Tomorrow.TempFahrenheitMin = jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Temperature"].(map[string]any)["Minimum"].(map[string]any)["Value"].(float64)
	w.Tomorrow.TempFahrenheitMax = jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Temperature"].(map[string]any)["Maximum"].(map[string]any)["Value"].(float64)
	w.Tomorrow.TempCelsiusMin = ftoC(w.Tomorrow.TempFahrenheitMin)
	w.Tomorrow.TempCelsiusMax = ftoC(w.Tomorrow.TempFahrenheitMax)
	w.Tomorrow.WeatherIcon = int(jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Day"].(map[string]any)["Icon"].(float64))

	// UV Index
	w.UVIndex = int(jsonData["DailyForecasts"].([]any)[day].(map[string]any)["AirAndPollen"].([]any)[5].(map[string]any)["Value"].(float64))
	if w.UVIndex > 12 {
		w.UVIndex = 12
	}

	// Wind Today
	w.Wind.WindImperial = jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Speed"].(map[string]any)["Value"].(float64)
	w.Wind.WindMetric = mphToKM(jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Speed"].(map[string]any)["Value"].(float64))
	w.Wind.WindDirection = jsonData["DailyForecasts"].([]any)[day].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Direction"].(map[string]any)["English"].(string)

	// Wind Tomorrow
	w.Wind.WindImperialTomorrow = jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Speed"].(map[string]any)["Value"].(float64)
	w.Wind.WindMetricTomorrow = mphToKM(jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Speed"].(map[string]any)["Value"].(float64))
	w.Wind.WindDirectionTomorrow = jsonData["DailyForecasts"].([]any)[day+1].(map[string]any)["Day"].(map[string]any)["Wind"].(map[string]any)["Direction"].(map[string]any)["English"].(string)

	// Pollen index
	// TODO: Properly calculate pollen. Something with force type casting validation
	grass := 2
	tree := 2
	ragweed := 2
	w.Pollen = (grass + tree + ragweed) / 3

	// Complete precipitation
	for i := 8; i < 15; i++ {
		w.Precipitation[i] = int(jsonData["DailyForecasts"].([]any)[i-8+day].(map[string]any)["Day"].(map[string]any)["PrecipitationProbability"].(float64))
	}

	w.Week = make([]Week, 7)
	_i := 1
	for i := 0; i < 7; i++ {
		w.Week[i].TempFahrenheitMin = jsonData["DailyForecasts"].([]any)[day+_i].(map[string]any)["Temperature"].(map[string]any)["Minimum"].(map[string]any)["Value"].(float64)
		w.Week[i].TempFahrenheitMax = jsonData["DailyForecasts"].([]any)[day+_i].(map[string]any)["Temperature"].(map[string]any)["Maximum"].(map[string]any)["Value"].(float64)
		w.Week[i].TempCelsiusMin = ftoC(w.Week[i].TempFahrenheitMin)
		w.Week[i].TempCelsiusMax = ftoC(w.Week[i].TempFahrenheitMax)
		w.Week[i].WeatherIcon = int(jsonData["DailyForecasts"].([]any)[day+_i].(map[string]any)["Day"].(map[string]any)["Icon"].(float64))
		_i++
	}
}

func ftoC(f float64) float64 {
	return (f - 32) * 5 / 9
}

func mphToKM(mph float64) float64 {
	return mph * 1.60934
}
