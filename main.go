package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrflynn/go-aqi"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	ApiBaseURL = "https://www.purpleair.com/json?show="
	SensorID   = "49489"
)

type SensorData struct {
	AQI         float64
	AQIName     string
	Temperature string
	Humidity    string
	Pressure    string
	Uptime      string
}

func (s *SensorData) Init(api ApiResponse) {
	// Get 'normal' stats from ApiResponse
	stats := getStats(api)

	// Set all values
	s.AQI, s.AQIName = calcAQI(stats)
	s.Temperature = FtoC(api.Results[0].TempF)
	s.Humidity = api.Results[0].Humidity
	s.Pressure = api.Results[0].Pressure
	s.Uptime = api.Results[0].Uptime
}

// Prints out Sensor Data
func (s *SensorData) Print() {
	fmt.Printf("Purple Air Sensor #%s data:\n", SensorID)
	fmt.Printf("- AQI: %.0f (%s)\n", s.AQI, s.AQIName)
	fmt.Printf("- Temperature: %s\n", s.Temperature)
	fmt.Printf("- Humidity: %s\n", s.Humidity)
	fmt.Printf("- Pressure: %s\n", s.Pressure)
}

type Stats struct {
	V                 float64 `json:"v"`
	V1                float64 `json:"v1"`
	V2                float64 `json:"v2"`
	V3                float64 `json:"v3"`
	V4                float64 `json:"v4"`
	V5                float64 `json:"v5"`
	V6                float64 `json:"v6"`
	Pm                float64 `json:"pm"`
	LastModified      int64   `json:"lastModified"`
	TimeSinceModified int     `json:"timeSinceModified"`
}

type ApiResponse struct {
	MapVersion       string `json:"mapVersion"`
	BaseVersion      string `json:"baseVersion"`
	MapVersionString string `json:"mapVersionString"`
	Results          []struct {
		ID                           int     `json:"ID"`
		Label                        string  `json:"Label"`
		DEVICELOCATIONTYPE           string  `json:"DEVICE_LOCATIONTYPE,omitempty"`
		THINGSPEAKPRIMARYID          string  `json:"THINGSPEAK_PRIMARY_ID"`
		THINGSPEAKPRIMARYIDREADKEY   string  `json:"THINGSPEAK_PRIMARY_ID_READ_KEY"`
		THINGSPEAKSECONDARYID        string  `json:"THINGSPEAK_SECONDARY_ID"`
		THINGSPEAKSECONDARYIDREADKEY string  `json:"THINGSPEAK_SECONDARY_ID_READ_KEY"`
		Lat                          float64 `json:"Lat"`
		Lon                          float64 `json:"Lon"`
		PM25Value                    string  `json:"PM2_5Value,omitempty"`
		LastSeen                     int     `json:"LastSeen"`
		Type                         string  `json:"Type,omitempty"`
		Hidden                       string  `json:"Hidden"`
		DEVICEBRIGHTNESS             string  `json:"DEVICE_BRIGHTNESS,omitempty"`
		DEVICEHARDWAREDISCOVERED     string  `json:"DEVICE_HARDWAREDISCOVERED,omitempty"`
		Version                      string  `json:"Version,omitempty"`
		LastUpdateCheck              int     `json:"LastUpdateCheck,omitempty"`
		Created                      int     `json:"Created"`
		Uptime                       string  `json:"Uptime,omitempty"`
		RSSI                         string  `json:"RSSI,omitempty"`
		Adc                          string  `json:"Adc,omitempty"`
		P03Um                        string  `json:"p_0_3_um,omitempty"`
		P05Um                        string  `json:"p_0_5_um,omitempty"`
		P10Um                        string  `json:"p_1_0_um,omitempty"`
		P25Um                        string  `json:"p_2_5_um,omitempty"`
		P50Um                        string  `json:"p_5_0_um,omitempty"`
		P100Um                       string  `json:"p_10_0_um,omitempty"`
		Pm10Cf1                      string  `json:"pm1_0_cf_1,omitempty"`
		Pm25Cf1                      string  `json:"pm2_5_cf_1,omitempty"`
		Pm100Cf1                     string  `json:"pm10_0_cf_1,omitempty"`
		Pm10Atm                      string  `json:"pm1_0_atm,omitempty"`
		Pm25Atm                      string  `json:"pm2_5_atm,omitempty"`
		Pm100Atm                     string  `json:"pm10_0_atm,omitempty"`
		IsOwner                      int     `json:"isOwner"`
		Humidity                     string  `json:"humidity,omitempty"`
		TempF                        string  `json:"temp_f,omitempty"`
		Pressure                     string  `json:"pressure,omitempty"`
		AGE                          int     `json:"AGE"`
		Stats                        string  `json:"Stats,omitempty"`
		ParentID                     int     `json:"ParentID,omitempty"`
	} `json:"results"`
}

// Makes actual HTTP call to 'ApiBaseURL' endpoint. Returns 'ApiResponse' object
func makeAPICall() ApiResponse {
	// Object for holding API response JSON
	apiResponse := ApiResponse{}

	url := ApiBaseURL + SensorID

	// Make HTTP call
	response, err := http.Get(url)
	if err != nil {
		log.Panic("Error while making HTTP Request:", err)
		return apiResponse
	}

	// Read HTTP response body
	body, _ := ioutil.ReadAll(response.Body)

	// Convert HTTP response to Go object.
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatal("Can't parse JSON from response body:", err)
		return apiResponse
	}

	return apiResponse
}

// Converts 'Stats' field from input 'ApiResponse' object from string type to 'Stats' object type
func getStats(a ApiResponse) Stats {
	var statsJSON Stats
	statsString := a.Results[0].Stats

	err := json.Unmarshal([]byte(statsString), &statsJSON)
	if err != nil {
		log.Fatal("Can't convert Stats string to JSON:", err)
	}

	return statsJSON
}

// Calculate AQI index based ob PM2.5 value from 'Stats' object.
// Output is: 'AQI index' and 'Index description'
func calcAQI(s Stats) (float64, string) {
	result, err := aqi.Calculate(aqi.PM25{Concentration: s.V1})
	if err != nil {
		log.Fatal("Can't calculate AQI based on PM2.5 value:", err)
	}
	return result.AQI, result.Index.Name
}

// Converts Fahrenheits to Celsius
func FtoC(f string) string {
	tempF, err := strconv.ParseFloat(f, 64)
	if err != nil {
		log.Fatal("Can't parse temperature:", err)
	}

	tempC := (tempF - 32) * 5 / 9

	return fmt.Sprintf("%f", tempC)
}

func main() {
	var data SensorData
	data.Init(makeAPICall())
	data.Print()
}
