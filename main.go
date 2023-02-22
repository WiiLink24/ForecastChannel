package main

import (
	"ForecastChannel/accuweather"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/wii-tools/lzx/lz10"
	"hash/crc32"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type Forecast struct {
	Header                 Header
	LocationTable          []LocationTable
	LocationText           []uint16
	LongForecastTable      []LongForecastTable
	ShortForecastTable     []ShortForecastTable
	LaundryTable           []LaundryTable
	LaundryText            []uint16
	WeatherConditionsTable []WeatherConditionsTable
	WeatherConditionsText  []uint16
	UVTable                []UVTable
	UVText                 []uint16
	PollenTable            []PollenTable
	PollenText             []uint16

	currentLanguageCode uint8
	currentCountryCode  uint8
	currentCountryList  *NationalList
	rawLocations        Locations
}

var weatherMap = map[string]*accuweather.Weather{}
var currentTime = time.Now().Unix()
var weatherList *WeatherList
var mapMutex = sync.RWMutex{}

func main() {
	start := time.Now()
	// Get all important data we need
	weatherList = ParseWeatherXML()
	PopulateCountryCodes()

	// Async HTTP done safely and fast
	wg := sync.WaitGroup{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	semaphore := make(chan struct{}, 10)

	// Next retrieve international weather
	wg.Add(len(weatherList.International.Cities))
	for _, city := range weatherList.International.Cities {
		city := city
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			fmt.Println("Processing", city.Name.English)
			weather := accuweather.GetWeather(city.Longitude, city.Latitude, currentTime)
			mapMutex.Lock()
			weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)] = weather
			mapMutex.Unlock()
			fmt.Println("Finished", city.Name.English)
			<-semaphore
		}()
	}
	wg.Wait()

	// We must get the number of national cities not yet generated
	numberOfCities := 0
	for _, cities := range weatherList.National {
		for _, city := range cities.Cities {
			if weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)] == nil {
				numberOfCities++
			}
		}
	}

	wg.Add(numberOfCities)
	for _, cities := range weatherList.National {
		for _, city := range cities.Cities {
			if weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)] == nil {
				city := city
				go func() {
					defer wg.Done()
					semaphore <- struct{}{}
					fmt.Println("Processing", city.English)
					weather := accuweather.GetWeather(city.Longitude, city.Latitude, currentTime)
					mapMutex.Lock()
					weatherMap[fmt.Sprintf("%f,%f", city.Longitude, city.Latitude)] = weather
					mapMutex.Unlock()
					fmt.Println("Finished", city.English)
					<-semaphore
				}()
			}
		}
	}
	wg.Wait()

	wg.Add(len(weatherList.National))
	for _, national := range weatherList.National {
		countryCode := countryCodes[national.Name.English]
		national := national
		go func() {
			defer wg.Done()
			for _, languageCode := range GetSupportedLanguages(countryCode) {
				forecast := Forecast{}
				forecast.currentCountryList = &national
				forecast.currentCountryCode = countryCode
				forecast.currentLanguageCode = languageCode
				forecast.PopulateLocations(weatherList.International.Cities)

				buffer := new(bytes.Buffer)
				forecast.MakeHeader()
				forecast.MakeLocationTable(weatherList.International.Cities)
				forecast.MakeLocationText(weatherList.International.Cities)
				forecast.MakeLongForecastTable()
				forecast.MakeShortForecastTable(weatherList.International.Cities)
				forecast.MakeLaundryTable()
				forecast.MakeLaundryText()
				forecast.MakeWeatherConditionsTable()
				forecast.MakeWeatherConditionText()
				forecast.MakeUVTable()
				forecast.MakeUVText()
				forecast.MakePollenTable()
				forecast.MakePollenText()
				forecast.WriteAll(buffer)

				forecast.Header.Filesize = uint32(buffer.Len())
				buffer.Reset()
				forecast.WriteAll(buffer)

				crcTable := crc32.MakeTable(crc32.IEEE)
				checksum := crc32.Checksum(buffer.Bytes()[12:], crcTable)
				forecast.Header.CRC32 = checksum

				buffer.Reset()
				forecast.WriteAll(buffer)

				// Prepare to make for Wii U
				// TODO: Patch IOS to force proper UTC
				wiiUBuffer := new(bytes.Buffer)
				forecast.Header.OpenTimestamp = 0
				forecast.Header.CloseTimestamp = 0xFFFFFFFF
				forecast.WriteAll(wiiUBuffer)

				crcTable = crc32.MakeTable(crc32.IEEE)
				checksum = crc32.Checksum(wiiUBuffer.Bytes()[12:], crcTable)
				forecast.Header.CRC32 = checksum

				wiiUBuffer.Reset()
				forecast.WriteAll(wiiUBuffer)

				// Make short.bin
				wii, wiiU := forecast.MakeShortBin(weatherList.International.Cities)

				// Make the folder if it doesn't already exist
				err := os.Mkdir(fmt.Sprintf("./files/%d/%s", languageCode, ZFill(countryCode, 3)), 0755)
				if !os.IsExist(err) {
					// If the folder exists we can just continue
					checkError(err)
				}

				compressed, err := lz10.Compress(buffer.Bytes())
				checkError(err)

				compressedU, err := lz10.Compress(wiiUBuffer.Bytes())
				checkError(err)

				err = os.WriteFile(fmt.Sprintf("./files/%d/%s/forecast.bin", languageCode, ZFill(countryCode, 3)), SignFile(compressed), 0666)
				checkError(err)

				err = os.WriteFile(fmt.Sprintf("./files/%d/%s/short.bin", languageCode, ZFill(countryCode, 3)), SignFile(wii), 0666)
				checkError(err)

				err = os.WriteFile(fmt.Sprintf("./files/%d/%s/forecast.alt", languageCode, ZFill(countryCode, 3)), SignFile(compressedU), 0666)
				checkError(err)

				err = os.WriteFile(fmt.Sprintf("./files/%d/%s/short.alt", languageCode, ZFill(countryCode, 3)), SignFile(wiiU), 0666)
				checkError(err)
			}
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Forecast Channel file generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func Write(writer io.Writer, data any) {
	err := binary.Write(writer, binary.BigEndian, data)
	checkError(err)
}

func (f *Forecast) WriteAll(writer io.Writer) {
	Write(writer, f.Header)
	Write(writer, f.LocationTable)
	Write(writer, f.LocationText)
	Write(writer, f.LongForecastTable)
	Write(writer, f.ShortForecastTable)
	Write(writer, f.LaundryTable)
	Write(writer, f.LaundryText)
	Write(writer, f.WeatherConditionsTable)
	Write(writer, f.WeatherConditionsText)
	Write(writer, f.UVTable)
	Write(writer, f.UVText)
	Write(writer, f.PollenTable)
	Write(writer, f.PollenText)
}

func (f *Forecast) GetCurrentSize() uint32 {
	buffer := bytes.NewBuffer(nil)
	f.WriteAll(buffer)

	return uint32(buffer.Len())
}
