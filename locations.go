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
	CountryCode  uint8
	RegionCode   uint8
	LocationCode uint16
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

		locations[f.currentCountryList.Name.English][city.Province.English][city.English] = Location{
			CountryCode:  f.currentCountryCode,
			RegionCode:   uint8(len(locations[f.currentCountryList.Name.English]) + 1),
			LocationCode: uint16(len(locations[f.currentCountryList.Name.English][city.Province.English]) + 1),
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

		if _, ok := countryCodes[city.Country.English]; !ok && city.Province.English == "" {
			locations[city.Country.English][city.Province.English][city.Name.English] = Location{
				CountryCode:  0xFE,
				RegionCode:   0xFE,
				LocationCode: uint16(noProvince),
			}
			noProvince++
			continue
		}

		if _, ok := locations[city.Country.English][city.Province.English][city.Name.English]; ok {
			continue
		}

		locations[city.Country.English][city.Province.English][city.Name.English] = Location{
			CountryCode:  countryCodes[city.Country.English],
			RegionCode:   uint8(len(locations[city.Country.English]) + 1),
			LocationCode: uint16(len(locations[city.Country.English][city.Province.English]) + 1),
		}
	}

	f.rawLocations = locations
}

func (f *Forecast) MakeLocationTable(cities []InternationalCity) {
	f.Header.LocationsTableOffset = f.GetCurrentSize()

	for _, city := range f.currentCountryList.Cities {
		f.LocationTable = append(f.LocationTable, LocationTable{
			CountryCode:       f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].CountryCode,
			RegionCode:        f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].RegionCode,
			LocationCode:      f.rawLocations[f.currentCountryList.Name.English][city.Province.English][city.English].LocationCode,
			CityTextOffset:    0,
			RegionTextOffset:  0,
			CountryTextOffset: 0,
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
		})
	}

	for _, city := range cities {
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		f.LocationTable = append(f.LocationTable, LocationTable{
			CountryCode:       f.rawLocations[city.Country.English][city.Province.English][city.Name.English].CountryCode,
			RegionCode:        f.rawLocations[city.Country.English][city.Province.English][city.Name.English].RegionCode,
			LocationCode:      f.rawLocations[city.Country.English][city.Province.English][city.Name.English].LocationCode,
			CityTextOffset:    0,
			RegionTextOffset:  0,
			CountryTextOffset: 0,
			Latitude:          CoordinateEncode(city.Latitude),
			Longitude:         CoordinateEncode(city.Longitude),
			LocationZoom1:     uint8(city.Zoom1),
			LocationZoom2:     uint8(city.Zoom2),
		})
	}

	f.Header.NumberOfLocations = uint32(len(f.LocationTable))
}

func (f *Forecast) MakeLocationText(cities []InternationalCity) {
	i := 0
	for _, city := range f.currentCountryList.Cities {
		f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetCityName(city)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}

		f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(city.Province)))...)
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
		i++
	}

	for _, city := range cities {
		if city.Country.English == f.currentCountryList.Name.English {
			continue
		}

		f.LocationTable[i].CityTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(city.Name)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}

		if city.Province.English != "" {
			f.LocationTable[i].RegionTextOffset = f.GetCurrentSize()
			f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(city.Province)))...)
			f.LocationText = append(f.LocationText, uint16(0))
			for f.GetCurrentSize()&3 != 0 {
				f.LocationText = append(f.LocationText, uint16(0))
			}
		}

		f.LocationTable[i].CountryTextOffset = f.GetCurrentSize()
		f.LocationText = append(f.LocationText, utf16.Encode([]rune(f.GetLocalizedName(city.Country)))...)
		f.LocationText = append(f.LocationText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LocationText = append(f.LocationText, uint16(0))
		}
		i++
	}
}
