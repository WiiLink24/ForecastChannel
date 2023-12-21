package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"os"
	"testing"
	"unicode/utf16"
)

func TestDuplicates(t *testing.T) {
	// Every country supports English so we will read from that
	PopulateCountryCodes()
	for c, u := range countryCodes {
		data, err := os.ReadFile(fmt.Sprintf("files/1/%s/forecast.bin", ZFill(u, 3)))
		if data == nil {
			fmt.Println(fmt.Sprintf("Skipping Country %s", c))
			continue
		}
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
			var locationOffsets []uint32
			for ii := 0; ii < int(header.NumberOfLocations); ii++ {
				currentLocationCode := binary.BigEndian.Uint32(decompressed[locationTableOffset:])
				if currentLocationCode == locationCode {
					count++
					locationOffsets = append(locationOffsets, locationTableOffset)
				}

				locationTableOffset += 24
			}

			if count != 1 {
				fmt.Println(fmt.Sprintf("Error in Country %s", c))
				fmt.Println(fmt.Sprintf("Duplicate Detected. Count: %d, Location Code: %d", count, locationCode))

				for _, offset := range locationOffsets {
					theNameOffset := binary.BigEndian.Uint32(decompressed[offset+4:])
					var name []uint16
					for {
						// Find the name of the city I hope
						currBytes := binary.BigEndian.Uint16(decompressed[theNameOffset:])
						if currBytes == 0 {
							break
						}

						name = append(name, currBytes)
						theNameOffset += 2
					}

					fmt.Println("City Name: ", string(utf16.Decode(name)))
				}
			}

			shortTableOffset += 72
		}
	}
}
