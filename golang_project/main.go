package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const apiKey = "41e7400148ca49e397d192030252311"

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC float64 `json:"temp_c"`

		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`

		Humidity     int     `json:"humidity"`
		WindKph      float64 `json:"wind_kph"`
		FeelsLikeC   float64 `json:"feelslike_c"`
		IsDay        int     `json:"is_day"`
		WindDir      string  `json:"wind_dir"`
		Cloud        int     `json:"cloud"`
		Gustkph      float64 `json:"gust_kph"`
		PressureMb   float64 `json:"pressure_mb"`
		VisibilityKm float64 `json:"vis_km"`
	} `json:"current"`
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func main() {
	fmt.Println("Welcome to the GO prognoza 🌞 CLI app koristeci")
	fmt.Println("-----------------------------")

	city := getUserInput("Unesi ime grada: ")

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	fmt.Println("Api URL: ", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error u pravljenju zahteva", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get podatke prognoze (HTTP %d)\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Greska u ucitavanju odgovora: ", err)
		os.Exit(1)
	}

	var weather WeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		fmt.Println("X Error parsing JSON", err)
		os.Exit(1)
	}

	fmt.Println("\n ccc JSON data parsec uspesno!")
	displayWeather(weather)

}

func displayWeather(w WeatherResponse) {
	location := fmt.Sprintf("%s, %s", w.Location.Name, w.Location.Country)
	condition := w.Current.Condition.Text
	temp := w.Current.TempC
	feelsLike := w.Current.FeelsLikeC
	humidity := w.Current.WindKph
	isDay := w.Current.IsDay

	emoji := weatherEmoji(condition, isDay)

	fmt.Println("Location: ", location)
	fmt.Printf("Temperature: %.1f (feels like %.1f)%s\n", temp, feelsLike, emoji)
	fmt.Println("Feels: ", feelsLike)
	fmt.Println("Humidity: ", humidity)
	fmt.Println("Emoji: ", emoji)
	// fmt.Println("")
	// fmt.Println()
	// fmt.Println()

	if isDay == 1 {
		fmt.Println("Sad je daytime !!")
	} else {
		fmt.Println("Sada je nocno vreme!!")
	}

}

func weatherEmoji(condition string, isDay int) any {
	condition = strings.ToLower(condition)

	switch {
	case strings.Contains(condition, "sunny"), strings.Contains(condition, "clear"):
		if isDay == 1 {
			return "🌞"
		}
		return "Mesec"

	case strings.Contains(condition, "rain"), strings.Contains(condition, "drizzler"):
		return "kisa oblak"

	case strings.Contains(condition, "cloud"):
		return "oblak"

	case strings.Contains(condition, "snow"):
		return "sneg"

	case strings.Contains(condition, "thunder"):
		return "munja"

	case strings.Contains(condition, "fog"):
		return "magla"
	}

	return ""
}
