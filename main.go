package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	OpenWeatherAPIKey = "00caff80d1374a0158d10c572dc19ef5"
	Name              = "six-shot"
	Port              = "https://go-http-server-production-ff45.up.railway.app"
)

type IPInfo struct {
	IP       string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func getIPInfo(visitorName string) (*IPInfo, error) {
	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	city, ok := data["city"].(string)
	if !ok || city == "" {
		return nil, fmt.Errorf("city not found in IP API response")
	}

	weatherResp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, OpenWeatherAPIKey))
	if err != nil {
		return nil, err
	}
	defer weatherResp.Body.Close()

	weatherBody, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		return nil, err
	}
	var weatherData map[string]interface{}
	if err := json.Unmarshal(weatherBody, &weatherData); err != nil {
		return nil, err
	}

	temperature, ok := weatherData["main"].(map[string]interface{})["temp"].(float64)
	if !ok {
		return nil, fmt.Errorf("temperature not found in weather API response")
	}

	if visitorName != "" {
		visitorName = strings.Trim(visitorName, "\"")
	}

	greeting := fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", visitorName, temperature, city)

	info := &IPInfo{
		IP:       data["ip"].(string),
		Location: city,
		Greeting: greeting,
	}

	return info, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	visitorName := r.URL.Query().Get("visitor_name")
	ipInfo, err := getIPInfo(visitorName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ipInfo)
}

func main() {
	http.HandleFunc("/api/hello", handler)
	fmt.Printf("Server listening on port %s\n", Port)
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		panic(err)
	}
}
