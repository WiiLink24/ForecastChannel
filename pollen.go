package main

import "unicode/utf16"

type PollenTable struct {
	Code       uint8
	_          [3]byte
	TextOffset uint32
}

func (f *Forecast) MakePollenTable() {
	f.Header.PollenCountTableOffset = f.GetCurrentSize()

	for _, pollen := range weatherList.Pollen {
		f.PollenTable = append(f.PollenTable, PollenTable{
			Code:       uint8(pollen.Code),
			TextOffset: 0,
		})
	}

	f.Header.NumberOfPollenCountTables = uint32(len(f.PollenTable))
}

func (f *Forecast) MakePollenText() {
	for i, pollen := range weatherList.Pollen {
		f.PollenTable[i].TextOffset = f.GetCurrentSize()
		f.PollenText = append(f.PollenText, utf16.Encode([]rune(pollen.Name))...)
		f.PollenText = append(f.PollenText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.PollenText = append(f.PollenText, uint16(0))
		}
	}

}
