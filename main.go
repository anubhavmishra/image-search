package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
)

const version = "0.0.1"

func main() {

	var httpBindAddr = "0.0.0.0"
	var httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	httpAddr := fmt.Sprintf("%s:%s", httpBindAddr, httpPort)
	log.Println("Starting image-search app...")

	var apiKey = os.Getenv("GIPHY_API_KEY")
	var aprilFools = os.Getenv("APRIL_FOOLS")
	featureFlag := false
	if aprilFools == "true" {
		featureFlag = true
	}

	mux := http.NewServeMux()
	mux.Handle("/", GiphyImageHandler(apiKey, featureFlag))
	mux.HandleFunc("/healthz", HealthCheck)

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
