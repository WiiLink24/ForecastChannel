package main

import "unicode/utf16"

type UVTable struct {
	Code       uint8
	_          [3]byte
	TextOffset uint32
}

func (f *Forecast) MakeUVTable() {
	f.Header.UVIndexTableOffset = f.GetCurrentSize()

	for i, _ := range weatherList.UV {
		f.UVTable = append(f.UVTable, UVTable{
			Code:       uint8(i),
			TextOffset: 0,
		})
	}

	f.Header.NumberOfUVIndexTables = uint32(len(f.UVTable))
}

func (f *Forecast) MakeUVText() {
	for i, uv := range weatherList.UV {
		f.UVTable[i].TextOffset = f.GetCurrentSize()
		f.UVText = append(f.UVText, utf16.Encode([]rune(f.GetLocalizedName(uv.Name)))...)

		f.UVText = append(f.UVText, uint16(0))
		for f.GetCurrentSize()&3 != 0 {
			f.UVText = append(f.UVText, uint16(0))
		}
	}

}
