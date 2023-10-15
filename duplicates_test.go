package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"os"
	"testing"
)

func TestDuplicates(t *testing.T) {
	// Every country supports English so we will read from that
	PopulateCountryCodes()
	for c, u := range countryCodes {
		data, err := os.ReadFile(fmt.Sprintf("files/1/%s/forecast.bin", ZFill(u, 3)))
		if err != nil {
			fmt.Println(fmt.Sprintf("Error in Country %s", c))
			t.Error(err)
		}

		decompressed, err := lz10.Decompress(data[320:])
		if err != nil {
			fmt.Println(fmt.Sprintf("Error in Country %s", c))
			t.Error(err)
		}

		var header Header
		err = binary.Read(bytes.NewReader(decompressed), binary.BigEndian, &header)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error in Country %s", c))
			t.Error(err)
		}

		shortTableOffset := header.ShortForecastTableOffset
		for i := 0; i < int(header.NumberOfShortForecastTables); i++ {
			locationCode := binary.BigEndian.Uint32(decompressed[shortTableOffset:])
			locationTableOffset := header.LocationsTableOffset
			count := 0
			for ii := 0; ii < int(header.NumberOfLocations); ii++ {
				currentLocationCode := binary.BigEndian.Uint32(decompressed[locationTableOffset:])
				if currentLocationCode == locationCode {
					count++
				}

				locationTableOffset += 24
			}

			if count != 1 {
				fmt.Println(fmt.Sprintf("Error in Country %s", c))
				fmt.Println(fmt.Sprintf("Duplicate Detected. Count: %d, Location Code: %d", count, locationCode))
			}

			shortTableOffset += 72
		}
	}
}
