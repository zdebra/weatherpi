package main

import (
	"flag"
	"fmt"
	"github.com/morus12/dht22"
	"log"
)

func main() {
	gpioPtr := flag.String("gpio", "17", "gpio where the DHT22 is connected")
	flag.Parse()

	gpio := fmt.Sprintf("GPIO_%s", *gpioPtr)
	log.Printf("reading sensor at %s", gpio)
	sensor := dht22.New(gpio)

	temperature, err := sensor.Temperature()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("temperature reading done")
	humidity, err := sensor.Humidity()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("temperature %d°C, humidity %d%%", temperature, humidity)
}
