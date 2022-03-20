package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//structs for busstop

type Route struct {
	Id int `json:"id"`
}

type Forecast struct {
	Forecast_seconds float64 `json:"forecast_seconds"`
	Route            Route   `json:"route"`
}

type Busstop struct {
	Forecast []Forecast `json:"forecast"`
	Id       int        `json:"id"`
	Invalid  bool       `json:"invalid"`
}

//structs for busline
type stats struct {
	Bearing   int     `json:"bearing"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Avg_speed float64 `json:"avg_speed"`
}

type Vehicle struct {
	Stats stats `json:"stats"`
}

type Busline struct {
	Id      int       `json:"id"`
	Vehicle []Vehicle `json:"vehicles"`
	Invalid bool      `json:"invalid"`
}

func main() {
	handleRequests()
}

//helper function to check if id is within array of valid ids
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//match path with func
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/busstop/", returnbusstop)
	http.HandleFunc("/busline/", returnbusline)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//respond to /busstop/
func returnbusstop(w http.ResponseWriter, r *http.Request) {
	busstopList := []string{
		"378204", "383050", "378202", "383049", "382998", "378237", "378233", "378230",
		"378229", "378228", "378227", "382995", "378224", "378226", "383010", "383009",
		"383006", "383004", "378234", "383003", "378222", "383048", "378203", "382999",
		"378225", "383014", "383013", "383011", "377906", "383018", "383015", "378207",
	}
	fmt.Println("Endpoint Hit: returnbusstop")
	keys, err := r.URL.Query()["id"]
	if !err {
		log.Fatal(err)
	}
	if !contains(busstopList, keys[0]) {
		fmt.Println("invalid busstop")
		temperr := Busstop{}
		temperr.Invalid = true
		json.NewEncoder(w).Encode(temperr)
		return
	}
	json.NewEncoder(w).Encode(getbusstop(keys[0]))
}

//respond to /busline/
func returnbusline(w http.ResponseWriter, r *http.Request) {
	buslineList := []string{"44478", "44479", "44480", "44481"}
	fmt.Println("Endpoint Hit: returnbusline")
	keys, err := r.URL.Query()["id"]
	if !err {
		log.Fatal(err)
	}
	if !contains(buslineList, keys[0]) {
		fmt.Println("invalid busline")
		temperr := Busline{}
		temperr.Invalid = true
		json.NewEncoder(w).Encode(temperr)
		return
	}
	json.NewEncoder(w).Encode(getbusline(keys[0]))
}

//get information on arrival timing for busstop
func getbusstop(id string) Busstop {
	url := "https://dummy.uwave.sg/busstop/" + id

	uWavebusstop := http.Client{
		Timeout: time.Second * 1, // Timeout after 1 second
	}

	resp, err := uWavebusstop.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	busstop1 := Busstop{}
	json.Unmarshal(body, &busstop1)
	busstop1.Invalid = false
	return busstop1
}

//get information on location of buses and avg speed for busline
func getbusline(id string) Busline {
	url := "https://dummy.uwave.sg/busline/" + id

	uWavebusstop := http.Client{
		Timeout: time.Second * 1, // Timeout after 1 second
	}

	resp, err := uWavebusstop.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	busline1 := Busline{}
	busline1.Invalid = false
	json.Unmarshal(body, &busline1)
	return busline1
}
