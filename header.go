package main

type Header struct {
	Version                        uint32
	Filesize                       uint32
	CRC32                          uint32
	OpenTimestamp                  uint32
	CloseTimestamp                 uint32
	CountryCode                    uint8
	_                              [3]byte
	LanguageCode                   uint8
	TemperatureFlag                uint8
	Unknown                        uint8
	_                              uint8
	MessageOffset                  uint32
	NumberOfLongForecastTables     uint32
	LongForecastTableOffset        uint32
	NumberOfShortForecastTables    uint32
	ShortForecastTableOffset       uint32
	NumberOfWeatherConditionTables uint32
	WeatherConditionTableOffset    uint32
	NumberOfUVIndexTables          uint32
	UVIndexTableOffset             uint32
	NumberOfLaundryIndexTables     uint32
	LaundryIndexTableOffset        uint32
	NumberOfPollenCountTables      uint32
	PollenCountTableOffset         uint32
	NumberOfLocations              uint32
	LocationsTableOffset           uint32
}

func (f *Forecast) MakeHeader() {
	f.Header = Header{
		Version:                        0,
		Filesize:                       0,
		CRC32:                          0,
		OpenTimestamp:                  fixTime(int(currentTime)),
		CloseTimestamp:                 fixTime(int(currentTime)) + 90,
		CountryCode:                    f.currentCountryCode,
		LanguageCode:                   f.currentLanguageCode,
		TemperatureFlag:                f.GetTemperatureFlag(),
		Unknown:                        1,
		MessageOffset:                  0,
		NumberOfLongForecastTables:     0,
		LongForecastTableOffset:        0,
		NumberOfShortForecastTables:    0,
		ShortForecastTableOffset:       0,
		NumberOfWeatherConditionTables: 0,
		WeatherConditionTableOffset:    0,
		NumberOfUVIndexTables:          0,
		UVIndexTableOffset:             0,
		NumberOfLaundryIndexTables:     0,
		LaundryIndexTableOffset:        0,
		NumberOfPollenCountTables:      0,
		PollenCountTableOffset:         0,
		NumberOfLocations:              0,
		LocationsTableOffset:           0,
	}
}
