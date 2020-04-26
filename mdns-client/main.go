package main

import (
	"encoding/json"
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
			if bytes, err := json.Marshal(entry); err != nil {
				fmt.Printf("New Entry Discovered: %s\n", string(bytes))
			} else {
				fmt.Printf("New Entry Discovered: %+v\n", entry)
			}
		}
	}()

	go func() {
		for {
			// Start the lookup
			fmt.Println("\nStart Discovering....")
			mdns.Lookup("register_device", entriesCh)
			time.Sleep(time.Second * 3)
		}

	}()

	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	close(entriesCh)

}
