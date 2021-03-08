package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrflynn/go-aqi"
	"io/ioutil"
	"net/http"
	"os"
)

const API_BASE_URL = "https://www.purpleair.com/json?show="

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

func getData(sensorID string) ApiResponse {
	Response := ApiResponse{}

	url := API_BASE_URL + sensorID

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failure:", err)
		return Response
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &Response)
	if err != nil {
		fmt.Println("Failure:", err)
		return Response
	}

	return Response
}

func main() {
	response := getData("49489")
	stats := response.Results[0].Stats
	var statsJSON Stats
	err := json.Unmarshal([]byte(stats), &statsJSON)
	if err != nil {
		fmt.Println("Failure:", err)
		os.Exit(1)
	}

	result, err := aqi.Calculate(aqi.PM25{statsJSON.V1})
	if err != nil {
		fmt.Println("Failure:", err)
		os.Exit(1)
	}
	fmt.Printf("%.0f - %s\n", result.AQI, result.Index.Name)
}
