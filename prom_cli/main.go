package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MichaelS11/go-dht"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const addr string = "5000"

var (
	temp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "room_temperature_celsius",
		Help: "Current temperature of the room.",
	})

	hmd = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "room_humidity_relative",
		Help: "Current humidity of the room.",
	})
)

func main() {

	reg := prometheus.NewRegistry()

	if err := reg.Register(temp); err != nil {
		fmt.Println("room_temperature_celsius not registered:", err)
	} else {
		fmt.Println("room_temperature_celsius registered.")
	}

	if err := reg.Register(hmd); err != nil {
		fmt.Println("room_humidity_relative not registered:", err)
	} else {
		fmt.Println("room_humidity_relative registered.")
	}

	go recordTemp()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	fmt.Printf("Server listening up on port %s\n", addr)
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}

func recordTemp() {
	err := dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		panic(err)
	}

	dht, err := dht.NewDHT("GPIO17", dht.Fahrenheit, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		panic(err)
	}

	for {
		log.Println("Start detecting...")
		humidity, tempF, err := dht.ReadRetry(11)
		if err != nil {
			fmt.Println("Read error:", err)
		}

		tempC := tempF2C(tempF)

		if tempF < 0 || humidity < 40 {
			log.Println("Error: 偵測失敗")
		} else {
			log.Printf("temperature: %v humidity: %v\n", tempC, humidity)
			temp.Set(tempC)
			hmd.Set(humidity)
		}

		time.Sleep(60 * time.Second)
	}
}

func tempF2C(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}
