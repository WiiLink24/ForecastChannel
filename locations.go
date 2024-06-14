package main

import (
	"github.com/wk8/go-ordered-map/v2"
	"unicode/utf16"
)

type LocationTable struct {
	CountryCode       uint8
	RegionCode        uint8
	LocationCode      uint16
	CityTextOffset    uint32
	RegionTextOffset  uint32
	CountryTextOffset uint32
	Latitude          int16
	Longitude         int16
	LocationZoom1     uint8
	LocationZoom2     uint8
	_                 uint16
}

// Location helps us keep metadata on a specific location
type Location struct {
	CountryCode       uint8
	RegionCode        uint8
	LocationCode      uint16
	Latitude          int16
	Longitude         int16
	LocationZoom1     uint8
	LocationZoom2     uint8
	City              *City
	InternationalCity *InternationalCity
}

func (f *Forecast) PopulateLocations(cities []InternationalCity) {
	locations := orderedmap.New[string, *orderedmap.OrderedMap[string, *orderedmap.OrderedMap[string, Location]]]()

	// First we will populate the national cities. This way we can ignore any duplicates in international.
	locations.Set(f.currentCountryList.Name.English, orderedmap.New[string, *orderedmap.OrderedMap[string, Location]]())
	for _, city := range f.currentCountryList.Cities {
		// Populate province slice if it doesn't exist
		country, _ := locations.Get(f.currentCountryList.Name.English)
		if _, ok := country.Get(city.Province.English); !ok {
			country.Set(city.Province.English, orderedmap.New[string, Location]())
		}

		province, _ := country.Get(city.Province.English)
		if _, ok := province.Get(city.English); ok {
			continue
		}

		_city := city
		province.Set(city.English, Location{
			CountryCode:       f.currentCountryCode,
			RegionCode:        uint8(country.Len()) + 1,
			LocationCode:      uint16(province.Len()) + 1,
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
			City:              &_city,
			InternationalCity: &InternationalCity{},
		})
	}

	noProvince := 1
	for _, city := range cities {
		// Check if city was already populated from national cities
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		if _, ok := locations.Get(city.Country.English); !ok {
			locations.Set(city.Country.English, orderedmap.New[string, *orderedmap.OrderedMap[string, Location]]())
		}

		country, _ := locations.Get(city.Country.English)
		if _, ok := country.Get(city.Province.English); !ok {
			country.Set(city.Province.English, orderedmap.New[string, Location]())
		}

		province, _ := country.Get(city.Province.English)
		if _, ok := province.Get(city.Name.English); ok {
			continue
		}

		_city := city
		if _, ok := countryCodes[city.Country.English]; !ok && city.Province.English == "" {
			province.Set(city.Name.English, Location{
				CountryCode:       0xFE,
				RegionCode:        0xFE,
				LocationCode:      uint16(noProvince),
				Latitude:          CoordinateEncode(city.Latitude),
				Longitude:         CoordinateEncode(city.Longitude),
				LocationZoom1:     uint8(city.Zoom1),
				LocationZoom2:     uint8(city.Zoom2),
				City:              &City{},
				InternationalCity: &_city,
			})
			noProvince++
			continue
		}

		province.Set(city.Name.English, Location{
			CountryCode:       countryCodes[city.Country.English],
			RegionCode:        uint8(country.Len()) + 1,
			LocationCode:      uint16(province.Len()) + 1,
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
			City:              &City{},
			InternationalCity: &_city,
		})
	}

	f.rawLocations = locations
}

func (f *Forecast) MakeLocationTable() {
	f.Header.LocationsTableOffset = f.GetCurrentSize()

	currentCountry, _ := f.rawLocations.Get(f.currentCountryList.Name.English)
	for province := currentCountry.Oldest(); province != nil; province = province.Next() {
		for city := province.Value.Oldest(); city != nil; city = city.Next() {
			f.LocationTable = append(f.LocationTable, LocationTable{
				CountryCode:       city.Value.CountryCode,
				RegionCode:        city.Value.RegionCode,
				LocationCode:      city.Value.LocationCode,
				CityTextOffset:    0,
				RegionTextOffset:  0,
				CountryTextOffset: 0,
				Latitude:          city.Value.Latitude,
				Longitude:         city.Value.Longitude,
				LocationZoom1:     city.Value.LocationZoom1,
				LocationZoom2:     city.Value.LocationZoom2,
			})

			f.cityNames = append(f.cityNames, *city.Value.City)
			f.internationalCities = append(f.internationalCities, *city.Value.InternationalCity)
		}
	}

	for country := f.rawLocations.Oldest(); country != nil; country = country.Next() {
		if country.Key == f.currentCountryList.Name.English {
			continue
		}

		for province := country.Value.Oldest(); province != nil; province = province.Next() {
			for city := province.Value.Oldest(); city != nil; city = city.Next() {
				f.LocationTable = append(f.LocationTable, LocationTable{
					CountryCode:       city.Value.CountryCode,
					RegionCode:        city.Value.RegionCode,
					LocationCode:      city.Value.LocationCode,
					CityTextOffset:    0,
					RegionTextOffset:  0,
					CountryTextOffset: 0,
					Latitude:          city.Value.Latitude,
					Longitude:         city.Value.Longitude,
					LocationZoom1:     city.Value.LocationZoom1,
					LocationZoom2:     city.Value.LocationZoom2,
				})

				f.cityNames = append(f.cityNames, *city.Value.City)
				f.internationalCities = append(f.internationalCities, *city.Value.InternationalCity)
			}
		}
	}

	f.Header.NumberOfLocations = uint32(len(f.LocationTable))
}

func (f *Forecast) MakeLocationText() {
	// Map with text as a key and offset as value. This ensures we only write text once.
	writtenText := make(map[string]uint32)

	for i, city := range f.LocationTable {
		if city.CountryCode != f.currentCountryCode {
			continue
		}

		if value, ok := writtenText[f.GetCityName(f.cityNames[i])]; ok {
			f.LocationTable[i].CityTextOffset = value
		} else {
			f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetCityName(f.cityNames[i])))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}

			writtenText[f.GetCityName(f.cityNames[i])] = f.LocationTable[i].CityTextOffset
		}

		if value, ok := writtenText[f.GetLocalizedName(f.cityNames[i].Province)]; ok {
			f.LocationTable[i].RegionTextOffset = value
		} else {
			f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.cityNames[i].Province)))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}

			writtenText[f.GetLocalizedName(f.cityNames[i].Province)] = f.LocationTable[i].RegionTextOffset
		}

		if value, ok := writtenText[f.GetLocalizedName(f.currentCountryList.Name)]; ok {
			f.LocationTable[i].CountryTextOffset = value
		} else {
			f.LocationTable[i].CountryTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.currentCountryList.Name)))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}

			writtenText[f.GetLocalizedName(f.currentCountryList.Name)] = f.LocationTable[i].CountryTextOffset
		}
	}

	// Now do international cities
	for i, city := range f.LocationTable {
		if city.CountryCode == f.currentCountryCode {
			continue
		}

		if value, ok := writtenText[f.GetLocalizedName(f.internationalCities[i].Name)]; ok {
			f.LocationTable[i].CityTextOffset = value
		} else {
			f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Name)))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}

			writtenText[f.GetLocalizedName(f.internationalCities[i].Name)] = f.LocationTable[i].CityTextOffset
		}

		if f.internationalCities[i].Province.English != "" {
			if value, ok := writtenText[f.GetLocalizedName(f.internationalCities[i].Province)]; ok {
				f.LocationTable[i].RegionTextOffset = value
			} else {
				f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
				f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Province)))...)
				f.LocationText = append(f.LocationText, uint16(0))
				for f.GetCurrentSize()&3 != 0 {
					f.LocationText = append(f.LocationText, uint16(0))
				}

				writtenText[f.GetLocalizedName(f.internationalCities[i].Province)] = f.LocationTable[i].RegionTextOffset
			}
		}

		if f.internationalCities[i].Country.English != "" {
			if value, ok := writtenText[f.GetLocalizedName(f.internationalCities[i].Country)]; ok {
				f.LocationTable[i].CountryTextOffset = value
			} else {
				f.LocationTable[i].CountryTextOffset = f.GetCurrentSize()
				f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Country)))...)
				f.LocationText = append(f.LocationText, uint16(0))
				for f.GetCurrentSize()&3 != 0 {
					f.LocationText = append(f.LocationText, uint16(0))
				}

				writtenText[f.GetLocalizedName(f.internationalCities[i].Country)] = f.LocationTable[i].CountryTextOffset
			}
		}
	}
}
