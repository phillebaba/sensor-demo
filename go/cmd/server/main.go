package main

import (
	"fmt"
	"log"
	//"net"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/phillebaba/sensor-demo/go/pkg/api"
)

type config struct {
	BrokerIp     string `env:"MQTT_BROKER_IP" envDefault:"localhost"`
	BrokerPort   int    `env:"MQTT_BROKER_PORT" envDefault:"1883"`
	MessageTopic string `env:"MQTT_MESSAGE_TOPIC" envDefault:"temperature"`
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
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	// create a server instance
	s := api.Server{}

	// Subscribe to topic
	log.Println("Subscribing to topic")
	client.Subscribe(cfg.MessageTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		temperature, _ := strconv.ParseFloat(string(msg.Payload()), 64)
		if s.Stream != nil {
			if err := s.Stream.Send(&api.TemperatureResponse{Temperature: temperature}); err != nil {
				//return err
			}
		}
	})

	// grpc
	grpcServer := grpc.NewServer()
	api.RegisterTemperatureServiceServer(grpcServer, &s)

	corsOption := grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	})
	wrappedServer := grpcweb.WrapServer(grpcServer, corsOption)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
		/*if wrappedServer.IsGrpcWebRequest(req) {
			wrappedServer.ServeHTTP(resp, req)
		}

		http.DefaultServeMux.ServeHTTP(resp, req)*/
	}

	httpServer := http.Server{
		Addr:    "localhost:7777",
		Handler: http.HandlerFunc(handler),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed starting http server: %v", err)
	}
}
