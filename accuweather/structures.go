package accuweather

type Weather struct {
	LocalTime     string
	Current       Current
	Today         Today
	Tomorrow      Tomorrow
	Week          []Week
	Wind          Wind
	Precipitation []int
	UVIndex       int
	Pollen        int
	HourlyIcon    []int
	Globe         Globe

	apiKey string
}

type Current struct {
	TempFahrenheit float64
	TempCelsius    float64
	WindDirection  string
	WindImperial   float64
	WindMetric     float64
	WeatherIcon    int
}

type Today struct {
	TempFahrenheitMin float64
	TempFahrenheitMax float64
	TempCelsiusMin    float64
	TempCelsiusMax    float64
	WeatherIcon       int
}

type Tomorrow struct {
	TempFahrenheitMin float64
	TempFahrenheitMax float64
	TempCelsiusMin    float64
	TempCelsiusMax    float64
	WeatherIcon       int
}

type Week struct {
	TempFahrenheitMin float64
	TempFahrenheitMax float64
	TempCelsiusMin    float64
	TempCelsiusMax    float64
	WeatherIcon       int
}

type Wind struct {
	WindDirection         string
	WindImperial          float64
	WindMetric            float64
	WindDirectionTomorrow string
	WindImperialTomorrow  float64
	WindMetricTomorrow    float64
}

type Globe struct {
	Offset int
	Time   int
}
