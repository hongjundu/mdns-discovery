package main

import (
	"flag"
	"fmt"
	"log"
	"mdns-discovery/mdns-server/server"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	port := flag.Int("port", 7790, "http port")
	flag.Parse()

	fmt.Printf("port: %d\n", *port)

	// Setup our service export
	//host, _ := os.Hostname()
	//fmt.Printf("server host: %s\n", host)
	//
	//ip, _ := utils.IpAddress()
	//fmt.Printf("ip: %v \n", ip)
	//
	//info := []string{"Register Service", "TCP"}
	//service, _ := mdns.NewMDNSService(host, "_register._tcp", "", "", 5200, []net.IP{ip}, info)
	//
	//// Create the mDNS server, defer shutdown
	//server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	//defer server.Shutdown()

	server := server.NewHttpServer()
	go func() {
		log.Fatal(server.Run(*port))
	}()

	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
