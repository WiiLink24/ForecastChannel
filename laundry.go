package main

import "unicode/utf16"

type LaundryTable struct {
	Code       uint8
	_          [3]byte
	TextOffset uint32
}

func (f *Forecast) MakeLaundryTable() {
	f.Header.LaundryIndexTableOffset = f.GetCurrentSize()

	for _, laundry := range weatherList.Laundry {
		f.LaundryTable = append(f.LaundryTable, LaundryTable{
			Code:       uint8(laundry.Code),
			TextOffset: 0,
		})
	}

	f.Header.NumberOfLaundryIndexTables = uint32(len(f.LaundryTable))
}

func (f *Forecast) MakeLaundryText() {
	for i, laundry := range weatherList.Laundry {
		f.LaundryTable[i].TextOffset = f.GetCurrentSize()
		f.LaundryText = append(f.LaundryText, utf16.Encode([]rune(laundry.Name))...)
		f.LaundryText = append(f.LaundryText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.LaundryText = append(f.LaundryText, uint16(0))
		}
	}
}
