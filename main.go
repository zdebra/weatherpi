package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/d2r2/go-dht"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

var (
	ErrInvalidAPIResponse = errors.New("invalid API response")
)

type bucket struct {
	Room        string  `json:"room"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

type dispatcher struct {
	apiUrl string
	room   string
}

func (d *dispatcher) sendData(temperature, humidity float32) error {
	b := bucket{
		Room:        d.room,
		Temperature: temperature,
		Humidity:    humidity,
	}
	marshaledBucket, _ := json.Marshal(b)
	buf := bytes.NewBuffer(marshaledBucket)
	resp, err := http.Post(d.apiUrl+"/store", "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidAPIResponse
	}
	return nil
}

func main() {
	gpioPtr := flag.Int("gpio", 17, "gpio where the DHT22 is connected")
	apiPtr := flag.String("api", "https://weatherpi:ZV3CzpCHGa3h6uC@development-195419.appspot.com", "api url")
	roomPtr := flag.String("room", "bedroom", "room where the sensor is located")
	flag.Parse()

	room := *roomPtr
	if room == "" {
		room = "generic room"
	}

	disp := dispatcher{
		apiUrl: *apiPtr,
		room:   room,
	}

	log.Printf("scaning sensor %d", *gpioPtr)

	for {
		temperature, humidity, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, *gpioPtr, false, 10)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
			temperature, humidity, retried)
		if err := disp.sendData(temperature, humidity); err != nil {
			log.Printf("failed to send data: %s", err.Error())
		}
	}
}
