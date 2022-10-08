package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func schedulle(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

type Status struct {
	Status Stats `json:"status"`
}

type Stats struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func writeDataJSON(t time.Time) {
	var status Status
	var min int = 1
	var max int = 100

	rand.Seed(time.Now().UnixNano())

	status.Status.Water = rand.Intn(max-min) + min
	status.Status.Wind = rand.Intn(max-min) + min

	file, _ := os.Create("status.json")
	defer file.Close()

	byteValue, _ := json.Marshal(status)

	file.Write(byteValue)
}

func getDataJSON() []string {
	var water, wind string
	var checkWater, checkWind int
	var status Status

	file, _ := os.Open("status.json")
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &status)

	checkWater = status.Status.Water
	checkWind = status.Status.Wind

	if checkWater < 5 {
		water = "Aman"
	} else if checkWater >= 6 && checkWater <= 8 {
		water = "Siaga"
	} else {
		water = "Bahaya"
	}

	if checkWind < 6 {
		wind = "Aman"
	} else if checkWind >= 7 && checkWind <= 15 {
		wind = "Siaga"
	} else {
		wind = "Bahaya"
	}

	return []string{water, wind, strconv.Itoa(checkWater), strconv.Itoa(checkWind)}
}

func setupWebServer() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/status", func(c *gin.Context) {
		data := getDataJSON()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"water":      data[0],
			"wind":       data[1],
			"waterValue": data[2],
			"windValue":  data[3],
		})
	})
	r.Run(":8080")
}

func main() {
	go schedulle(15000*time.Millisecond, writeDataJSON)
	setupWebServer()
}
