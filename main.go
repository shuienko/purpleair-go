package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mrflynn/go-aqi"
	"github.com/patrickmn/go-cache"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	ApiBaseURL      = "https://www.purpleair.com/json?show="
	GetAQIText      = "AQI Індекс 🔄"
	CacheExpiration = time.Minute * 2
	CacheCleanup    = CacheExpiration * 10
)

var c *cache.Cache
var BotToken, SensorID string

type SensorData struct {
	AQI         float64
	AQIDesc     string
	AQIColor    string
	Temperature string
	Humidity    string
	Pressure    string
	Uptime      string
}

// Initialize SensorData based on API response
func (s *SensorData) New(api ApiResponse) {
	// Get 'normal' stats from ApiResponse
	stats := getStats(api)

	// Set all values
	s.AQI = calcAQI(stats)
	s.Temperature = FtoC(api.Results[0].TempF)
	s.Humidity = api.Results[0].Humidity
	s.Pressure = api.Results[0].Pressure
	s.Uptime = api.Results[0].Uptime

	switch {
	case s.AQI >= 0 && s.AQI <= 50:
		s.AQIColor = "🟢"
		s.AQIDesc = "Добре"
	case s.AQI >= 51 && s.AQI <= 100:
		s.AQIColor = "⚪️"
		s.AQIDesc = "Прийнятно"
	case s.AQI >= 101 && s.AQI <= 150:
		s.AQIColor = "🟡"
		s.AQIDesc = "Ризик для людей з респіраторними хворобами"
	case s.AQI >= 151 && s.AQI <= 200:
		s.AQIColor = "🟠"
		s.AQIDesc = "Погано"
	case s.AQI >= 201 && s.AQI <= 300:
		s.AQIColor = "🟠"
		s.AQIDesc = "Дуже погано"
	case s.AQI >= 301 && s.AQI <= 400:
		s.AQIColor = "🔴"
		s.AQIDesc = "Небезпечно"
	case s.AQI >= 401 && s.AQI <= 500:
		s.AQIColor = "🆘"
		s.AQIDesc = "ДУЖЕ НЕБЕЗПЕЧНО"
	}
}

// Returns text message for telegram bot with AQI info
func (s *SensorData) AqiInfo() string {
	return fmt.Sprintf("%s *%.0f* | %s\n", s.AQIColor, s.AQI, s.AQIDesc)
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
		Stats                        string  `json:"Stats,omitempty"` // Stats object
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
		log.Println("Error while making HTTP Request:", err)
		return apiResponse
	}

	// Read HTTP response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Can't read from response body:", err)
	}

	// Convert HTTP response to Go object.
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Println("Can't parse JSON from response body:", err)
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
		log.Println("Can't convert Stats string to JSON:", err)
	}

	return statsJSON
}

// Calculate AQI index based ob PM2.5 value from 'Stats' object.
func calcAQI(s Stats) float64 {
	result, err := aqi.Calculate(aqi.PM25{Concentration: s.V1})
	if err != nil {
		log.Println("Can't calculate AQI based on PM2.5 value:", err)
	}

	return result.AQI
}

// Converts Fahrenheits to Celsius
func FtoC(f string) string {
	tempF, err := strconv.ParseFloat(f, 64)
	if err != nil {
		log.Println("Can't parse temperature:", err)
	}

	tempC := (tempF - 32) * 5 / 9

	return fmt.Sprintf("%f", tempC)
}

// Compose response message using cache
func GetSensorData() string {
	var cacheData interface{}
	var data SensorData

	// Use cache if not empty
	// If cache is empty the make HTTP call and set cache
	cacheData, _ = c.Get("sensorData")
	if cacheData != nil {
		data = cacheData.(SensorData)
	} else {
		data.New(makeAPICall())
		c.Set("sensorData", data, CacheExpiration)
	}

	return data.AqiInfo()
}

// ##### INIT #####
func init() {
	// Create Cache
	c = cache.New(CacheExpiration, CacheCleanup)

	// Get environment variables and check errors
	BotToken = os.Getenv("PURPLEAIR_BOT_TOKEN")
	if len(BotToken) == 0 {
		log.Fatal("PURPLEAIR_BOT_TOKEN environment variable is not set. Exit.")
	}
	SensorID = os.Getenv("PURPLEAIR_BOT_SENSOR_ID")
	if len(SensorID) == 0 {
		log.Fatal("PURPLEAIR_BOT_SENSOR_ID environment variable is not set. Exit.")
	}
}

func main() {
	// Create new bot entity
	b, err := tb.NewBot(tb.Settings{
		Token:  BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("Can't create new bot object:", err)
		return
	}

	// Set Reply keyboard
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnGetAQI := menu.Text(GetAQIText)
	menu.Reply(
		menu.Row(btnGetAQI),
	)

	// Add send options
	options := &tb.SendOptions{
		ParseMode:   "Markdown",
		ReplyMarkup: menu,
	}

	// Handle /start command
	b.Handle("/start", func(m *tb.Message) {
		_, err = b.Send(m.Sender, GetSensorData(), options)
		if err != nil {
			log.Println("Failed to respond on /start command:", err)
		}
	})

	// Handle button
	b.Handle(&btnGetAQI, func(m *tb.Message) {
		_, err = b.Send(m.Sender, GetSensorData(), options)
		if err != nil {
			log.Println("Failed to respond on btnGetAQI:", err)
		}
	})

	b.Start()
}
