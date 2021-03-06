package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jacobsa/go-serial/serial"
)

type config struct {
	Debug        bool   `env:"DEBUG" envDefault:"false"`
	BrokerIp     string `env:"MQTT_BROKER_IP" envDefault:"localhost"`
	BrokerPort   int    `env:"MQTT_BROKER_PORT" envDefault:"1883"`
	MessageTopic string `env:"MQTT_MESSAGE_TOPIC" envDefault:"temperature"`
	PortName     string `env:"DEVICE_PORT_NAME" envDefault:"/dev/ttyACM0"`
	BaudRate     uint   `env:"DEVICE_BAUD_RATE" envDefault:"9600"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	log.Println("%+v", cfg)

	// Setup mqtt connection
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.BrokerIp, cfg.BrokerPort))
	client := mqtt.NewClient(opts)
	log.Println("Connecting to mqtt broker")
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mqtt broker")

	port := openPort(cfg.PortName, cfg.BaudRate, cfg.Debug)
	defer port.Close()

	// Start reading serial data
	buffer := []byte{}
	readBuffer := make([]byte, 128)
	port.Read(readBuffer) // Clear buffer before beginning
	for {
		n, err := port.Read(readBuffer)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			log.Println("\nEOF")
			break
		}

		buffer = append(buffer, readBuffer[:n]...)
		results, lastIndex := scan(buffer)
		buffer = buffer[lastIndex:]

		for _, result := range results {
			client.Publish(cfg.MessageTopic, 0, false, result)
		}
	}
}

// Opens serial port and returns io object that can be read from.
func openPort(portName string, baudRate uint, debug bool) io.ReadWriteCloser {
	if debug {
		log.Println("Using debug serial device")
		return DebugSerialDevice{}
	}

	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	log.Println("Opening serial port")
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	log.Println("Opened serial port")

	return port
}

// Parse read buffer for float values encapsulated with "<" and ">".
func scan(buffer []byte) ([]string, int) {
	results := []string{}
	start := false
	result := []byte{}
	lastIndex := 0

	for i, value := range buffer {
		stringValue := string(value)
		if stringValue == "<" {
			start = true
		} else if stringValue == ">" && start == true {
			start = false
			results = append(results, string(result))

			result = []byte{}
			lastIndex = i + 1
		} else if start {
			result = append(result, value)
		}
	}

	return results, lastIndex
}

type DebugSerialDevice struct {
}

func (d DebugSerialDevice) Read(p []byte) (n int, err error) {
	time.Sleep(1000 * time.Millisecond)

	min := -25
	max := 50
	number := rand.Intn(max-min+1) + min
	b := []byte("<" + strconv.Itoa(number) + ">")
	for i, value := range b {
		p[i] = value
	}

	return len(b), nil
}

func (d DebugSerialDevice) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (d DebugSerialDevice) Close() error {
	return nil
}
