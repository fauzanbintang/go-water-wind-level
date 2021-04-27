package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

type DataStatus struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

func init() {
	go AutoReloadJSON()
}

func main() {
	http.HandleFunc("/", AutoReloadWeb)

	http.ListenAndServe(":8080", nil)
}

func AutoReloadJSON() {
	for {
		min := 1
		max := 21
		wind := rand.Intn(max-min) + min
		water := rand.Intn(max-min) + min

		data := DataStatus{}
		data.Status.Water = water
		data.Status.Wind = wind

		dataMarshal, err := json.Marshal(data)
		if err != nil {
			fmt.Println("error marshal: ", err)
		}

		_ = ioutil.WriteFile("data.json", dataMarshal, 0644)

		time.Sleep(time.Second * 2)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	jsonFile, errJson := os.Open("data.json")
	if errJson != nil {
		fmt.Println("Error read file json: ", errJson)
	}

	defer jsonFile.Close()

	byteValue, errByte := ioutil.ReadAll(jsonFile)
	if errByte != nil {
		fmt.Println("Error read file json: ", errByte)
	}

	var dataStatus DataStatus
	json.Unmarshal(byteValue, &dataStatus)

	waterVal := dataStatus.Status.Water
	windVal := dataStatus.Status.Wind

	var (
		waterStatus string
		windStatus  string
	)

	switch {
	case waterVal < 5:
		waterStatus = "aman"
	case waterVal <= 8:
		waterStatus = "siaga"
	case waterVal > 8:
		waterStatus = "bahaya"
	}

	switch {
	case windVal < 6:
		windStatus = "aman"
	case windVal <= 15:
		windStatus = "siaga"
	case windVal > 15:
		windStatus = "bahaya"
	}

	var data = map[string]string{
		"waterStatus": waterStatus,
		"waterVal":    strconv.Itoa(waterVal),
		"windStatus":  windStatus,
		"windVal":     strconv.Itoa(windVal),
	}

	t, err := template.ParseFiles("main.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Execute(w, data)
	return
}
