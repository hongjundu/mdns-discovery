package server

import (
	"encoding/json"
	"fmt"
	"github.com/hongjundu/go-rest-api-helper.v1"
	"github.com/micro/mdns"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type HttpServer struct {
	*mdnsServerCache
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		&mdnsServerCache{
			mdnsServers: map[string]*mdns.Server{},
			mutex:       &sync.RWMutex{},
		},
	}
}

func (server *HttpServer) Run(port int) error {
	http.HandleFunc("/register", server.register)
	http.HandleFunc("/deregister", server.deregister)
	http.HandleFunc("/", server.notFound)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

func (server *HttpServer) register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[HttpServer] register")

	var response interface{}
	var err error
	if strings.Compare(strings.ToUpper(r.Method), "POST") != 0 && strings.Compare(strings.ToUpper(r.Method), "PUT") != 0 {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err = apihelper.NewError(http.StatusMethodNotAllowed, "Method not allowed")
	} else {
		var entry serviceEntry
		if err = apihelper.ReadJsonRequestBody(r.Body, &entry); err == nil {
			if len(entry.UUID) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"uuid\"")
			}
			if len(entry.Service) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"service\"")
			}
			if len(entry.ServiceName) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"serviceName\"")
			}
			if len(entry.IP) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"ip\"")
			}
			if entry.Port == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"port\"")
			}
			if len(entry.Version) == 0 {
				// fine...
			}
			if len(entry.Protocol) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"protocol\"")
			}
			if len(entry.ProbeUrl) == 0 {
				// fine..
			}

			info := []string{entry.ServiceName, entry.Version, entry.Protocol, entry.ProbeUrl}

			if err == nil {
				host, _ := os.Hostname()
				service, _ := mdns.NewMDNSService(host, entry.Service, "", "", entry.Port, []net.IP{net.ParseIP(entry.IP)}, info)
				mdnsServer, _ := mdns.NewServer(&mdns.Config{Zone: service})

				server.add(entry.UUID, mdnsServer)

				response = apihelper.NewOKResponse(entry)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}

	if err != nil {
		response = apihelper.NewErrorResponse(err)
	}

	if bytes, err := json.Marshal(response); err == nil {
		if _, e := w.Write(bytes); e != nil {
			fmt.Printf("ERROR: %v", e)
		} else {
			fmt.Printf("%v\n", string(bytes))
		}
	}

}

func (server *HttpServer) deregister(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[HttpServer] deregister")

	var response interface{}
	var err error
	if strings.Compare(strings.ToUpper(r.Method), "POST") != 0 && strings.Compare(strings.ToUpper(r.Method), "PUT") != 0 {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response = apihelper.NewErrorResponse(apihelper.NewError(http.StatusMethodNotAllowed, "Method not allowed"))
	} else {
		var entry serviceEntry
		if err = apihelper.ReadJsonRequestBody(r.Body, &entry); err == nil {
			if len(entry.UUID) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				err = apihelper.NewError(http.StatusBadRequest, "No \"uuid\"")
			}

			if err == nil {
				if err = server.remove(entry.UUID); err == nil {
					fmt.Printf("%s remove successfully", entry.UUID)
				} else {
					fmt.Printf("ERROR: %v", err)
				}
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}

	if bytes, err := json.Marshal(response); err == nil {
		if _, e := w.Write(bytes); e != nil {
			fmt.Printf("ERROR: %v", e)
		}
	}
}

func (server *HttpServer) notFound(w http.ResponseWriter, r *http.Request) {
	response := apihelper.NewErrorResponse(apihelper.NewError(http.StatusNotFound, "method not found"))
	w.WriteHeader(http.StatusNotFound)
	if bytes, err := json.Marshal(response); err == nil {
		if _, e := w.Write(bytes); e != nil {
			fmt.Printf("ERROR: %v", e)
		}
	}
}

///////////////////////////////////
// internals

type mdnsServerCache struct {
	mdnsServers map[string]*mdns.Server
	mutex       *sync.RWMutex
}

func (cache *mdnsServerCache) add(uuid string, server *mdns.Server) {
	fmt.Printf("mdns add: %s %+v", uuid, server)

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.mdnsServers[uuid] = server
}

func (cache *mdnsServerCache) remove(uuid string) error {
	fmt.Printf("mdns remove: %s", uuid)

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if server, ok := cache.mdnsServers[uuid]; ok {
		if e := server.Shutdown(); e != nil {
			fmt.Printf("ERROR: %v", e)
		}
		delete(cache.mdnsServers, uuid)
		return nil
	} else {
		return fmt.Errorf("not found: %s", uuid)
	}
}
