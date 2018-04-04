package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/d2r2/go-dht"
	"log"
	"net/http"
	"time"
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
	apiUrl            string
	room              string
	basicAuthUsername string
	basicAuthPassword string
}

func (d *dispatcher) sendData(temperature, humidity float32) error {
	b := bucket{
		Room:        d.room,
		Temperature: temperature,
		Humidity:    humidity,
	}
	marshaledBucket, _ := json.Marshal(b)
	buf := bytes.NewBuffer(marshaledBucket)

	req, err := http.NewRequest(http.MethodPost, d.apiUrl+"/store", buf)
	if err != nil {
		return fmt.Errorf("failed to create http request for data payload: %v", err)
	}
	req.SetBasicAuth(d.basicAuthUsername, d.basicAuthPassword)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidAPIResponse
	}
	return nil
}

func main() {
	gpioPtr := flag.Int("gpio", 17, "gpio where the DHT22 is connected")
	apiPtr := flag.String("api", "http://35.205.150.247", "api url")
	roomPtr := flag.String("room", "bedroom", "room where the sensor is located")
	bAuthUsername := flag.String("username", "a", "api basic auth username")
	bAuthPassword := flag.String("password", "b", "api basic auth password")
	sensorReadPeriodMs := flag.Int64("period", 2000, "sensor read period in miliseconds")
	flag.Parse()

	room := *roomPtr
	if room == "" {
		room = "generic room"
	}

	disp := dispatcher{
		apiUrl:            *apiPtr,
		room:              room,
		basicAuthUsername: *bAuthUsername,
		basicAuthPassword: *bAuthPassword,
	}

	sensorReadPeriodDuration := time.Duration(*sensorReadPeriodMs) * time.Millisecond
	log.Printf("scaning sensor %d every %d seconds, api basic auth username is %s", *gpioPtr, sensorReadPeriodDuration.Seconds(), *bAuthUsername)

	for range time.Tick(sensorReadPeriodDuration) {
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
