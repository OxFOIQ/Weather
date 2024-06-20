package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float32 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	apiKey := "Your_API_Key"
	location := "tunisia"

	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", apiKey, location)
	response, err := http.Get(url)
	if err != nil {
		panic("HTTP GET request failed")
	}

	if response.StatusCode != http.StatusOK { //200 == http.StatusOK
		panic("API request failed")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic("Failed to read response body")
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	locationInfo, currentWeather, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Println("Location:", locationInfo.Name, ",", locationInfo.Country, currentWeather.TempC, currentWeather.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		message :=fmt.Sprintf("%s - %.0fC , %.0f%% , %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		if hour.ChanceOfRain < 40 {
			color.Magenta(message)
		}else {
			color.Red(message)
		}
	}

}
