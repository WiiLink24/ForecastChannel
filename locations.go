package main

import (
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

type Locations map[string]map[string]map[string]Location

func (f *Forecast) PopulateLocations(cities []InternationalCity) {
	locations := map[string]map[string]map[string]Location{}

	// First we will populate the national cities. This way we can ignore any duplicates in international.
	locations[f.currentCountryList.Name.English] = map[string]map[string]Location{}
	for _, city := range f.currentCountryList.Cities {
		// Populate province slice if it doesn't exist
		if locations[f.currentCountryList.Name.English][city.Province.English] == nil {
			locations[f.currentCountryList.Name.English][city.Province.English] = map[string]Location{}
		}

		if _, ok := locations[f.currentCountryList.Name.English][city.Province.English][city.English]; ok {
			continue
		}

		_city := city
		locations[f.currentCountryList.Name.English][city.Province.English][city.English] = Location{
			CountryCode:       f.currentCountryCode,
			RegionCode:        uint8(len(locations[f.currentCountryList.Name.English]) + 1),
			LocationCode:      uint16(len(locations[f.currentCountryList.Name.English][city.Province.English]) + 1),
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
			City:              &_city,
			InternationalCity: &InternationalCity{},
		}
	}

	noProvince := 1
	for _, city := range cities {
		// Check if city was already populated from national cities
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		if locations[city.Country.English] == nil {
			locations[city.Country.English] = map[string]map[string]Location{}
		}

		if locations[city.Country.English][city.Province.English] == nil {
			locations[city.Country.English][city.Province.English] = map[string]Location{}
		}

		_city := city
		if _, ok := countryCodes[city.Country.English]; !ok && city.Province.English == "" {
			locations[city.Country.English][city.Province.English][city.Name.English] = Location{
				CountryCode:       0xFE,
				RegionCode:        0xFE,
				LocationCode:      uint16(noProvince),
				City:              &City{},
				InternationalCity: &_city,
			}
			noProvince++
			continue
		}

		if _, ok := locations[city.Country.English][city.Province.English][city.Name.English]; ok {
			continue
		}

		locations[city.Country.English][city.Province.English][city.Name.English] = Location{
			CountryCode:       countryCodes[city.Country.English],
			RegionCode:        uint8(len(locations[city.Country.English]) + 1),
			LocationCode:      uint16(len(locations[city.Country.English][city.Province.English]) + 1),
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
			City:              &City{},
			InternationalCity: &_city,
		}
	}

	f.rawLocations = locations
}

func (f *Forecast) MakeLocationTable() {
	f.Header.LocationsTableOffset = f.GetCurrentSize()

	for _, province := range f.rawLocations[f.currentCountryList.Name.English] {
		for _, city := range province {
			f.LocationTable = append(f.LocationTable, LocationTable{
				CountryCode:       city.CountryCode,
				RegionCode:        city.RegionCode,
				LocationCode:      city.LocationCode,
				CityTextOffset:    0,
				RegionTextOffset:  0,
				CountryTextOffset: 0,
				Latitude:          city.Latitude,
				Longitude:         city.Longitude,
				LocationZoom1:     city.LocationZoom1,
				LocationZoom2:     city.LocationZoom2,
			})

			f.cityNames = append(f.cityNames, *city.City)
			f.internationalCities = append(f.internationalCities, *city.InternationalCity)
		}
	}

	for countryName, country := range f.rawLocations {
		if countryName == f.currentCountryList.Name.English {
			continue
		}

		for _, province := range country {
			for _, city := range province {
				f.LocationTable = append(f.LocationTable, LocationTable{
					CountryCode:       city.CountryCode,
					RegionCode:        city.RegionCode,
					LocationCode:      city.LocationCode,
					CityTextOffset:    0,
					RegionTextOffset:  0,
					CountryTextOffset: 0,
					Latitude:          city.Latitude,
					Longitude:         city.Longitude,
					LocationZoom1:     city.LocationZoom1,
					LocationZoom2:     city.LocationZoom2,
				})

				f.cityNames = append(f.cityNames, *city.City)
				f.internationalCities = append(f.internationalCities, *city.InternationalCity)
			}
		}
	}

	f.Header.NumberOfLocations = uint32(len(f.LocationTable))
}

func (f *Forecast) MakeLocationText() {
	for i, city := range f.LocationTable {
		if city.CountryCode != f.currentCountryCode {
			continue
		}

		f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetCityName(f.cityNames[i])))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}

		f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.cityNames[i].Province)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}

		f.LocationTable[i].CountryTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.currentCountryList.Name)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}
	}

	// Now do international cities
	for i, city := range f.LocationTable {
		if city.CountryCode == f.currentCountryCode {
			continue
		}

		f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Name)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}

		if f.internationalCities[i].Province.English != "" {
			f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Province)))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}
		}

		f.LocationTable[i].CountryTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(f.internationalCities[i].Country)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}
	}
}
