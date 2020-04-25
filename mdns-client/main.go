package main

import (
	"fmt"
	"github.com/micro/mdns"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	go func() {
		for {
			// Start the lookup
			mdns.Lookup("_foobar._tcp", entriesCh)
			time.Sleep(time.Second * 3)
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	close(entriesCh)

}
