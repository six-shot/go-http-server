package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	OpenWeatherAPIKey = "00caff80d1374a0158d10c572dc19ef5"
	Name              = "six-shot"
	Port              = "8000"
)

type IPInfo struct {
	IP       string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func getClientIP(r *http.Request) string {
	// Check the X-Forwarded-For header first, as it might be set by a proxy
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}

	// Fallback to using the remote address
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}

func getIPInfo(visitorName string, clientIP string) (*IPInfo, error) {
	ipapiURL := fmt.Sprintf("https://ipapi.co/%s/json/", clientIP)
	resp, err := http.Get(ipapiURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching IP info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ipapi.co responded with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading IP API response: %w", err)
	}

	log.Printf("IP API Response: %s", body) // Log the response body for debugging

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling IP API response: %w", err)
	}

	log.Printf("Unmarshaled IP Data: %#v", data)

	city, ok := data["city"].(string)
	if !ok || city == "" {
		return nil, fmt.Errorf("city not found in IP API response")
	}

	latitude, ok := data["latitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("latitude not found in IP API response")
	}

	longitude, ok := data["longitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("longitude not found in IP API response")
	}

	weatherURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", latitude, longitude, OpenWeatherAPIKey)
	log.Printf("Weather API Request URL: %s", weatherURL) // Log the request URL

	weatherResp, err := http.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather info: %w", err)
	}
	defer weatherResp.Body.Close()

	if weatherResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(weatherResp.Body)
		return nil, fmt.Errorf("openweathermap.org responded with status %d: %s", weatherResp.StatusCode, string(body))
	}

	weatherBody, err := io.ReadAll(weatherResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading weather API response: %w", err)
	}

	log.Printf("Weather API Response: %s", weatherBody) // Log the response body for debugging

	var weatherData map[string]interface{}
	if err := json.Unmarshal(weatherBody, &weatherData); err != nil {
		return nil, fmt.Errorf("error unmarshaling weather API response: %w", err)
	}

	log.Printf("Unmarshaled Weather Data: %#v", weatherData)

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
	clientIP := getClientIP(r)
	ipInfo, err := getIPInfo(visitorName, clientIP)
	if err != nil {
		log.Printf("Error getting IP info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print the JSON
	if err := encoder.Encode(ipInfo); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/hello", handler)
	fmt.Printf("Server listening on port %s\n", Port)
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
