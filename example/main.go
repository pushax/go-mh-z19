package main

import (
	"fmt"
	"github.com/pushax/go-mh-z19"
	"log"
	"time"
)

func main() {
	sensor, err := mhz19.New("/dev/ttyAMA0")
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer func() {
		_ = sensor.Close()
	}()

	fmt.Println("sensor is ready")

	for {
		co2, err := sensor.ReadCO2()
		if err != nil {
			log.Printf("reading error: %v", err)
		} else {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] CO2: %d ppm\n", timestamp, co2)

			var status string
			switch {
			case co2 < 400:
				status = "Very low"
			case co2 < 600:
				status = "Very good"
			case co2 < 1000:
				status = "Good"
			case co2 < 1500:
				status = "Ok"
			case co2 < 2500:
				status = "Bad"
			default:
				status = "Very bad"
			}
			fmt.Printf("air quality: %s\n\n", status)
		}

		time.Sleep(10 * time.Second)
	}
}
