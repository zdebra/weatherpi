package main

import (
	"flag"
	"github.com/d2r2/go-dht"
	"log"
)

func main() {
	gpioPtr := flag.Int("gpio", 17, "gpio where the DHT22 is connected")
	flag.Parse()

	log.Printf("scaning sensor %d", *gpioPtr)

	temperature, humidity, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, *gpioPtr, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		temperature, humidity, retried)
}
