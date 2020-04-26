package server

type serviceEntry struct {
	UUID        string `json:"uuid"`
	Service     string `json:"service"`
	ServiceName string `json:"serviceName"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Version     string `json:"version"`
	Protocol    string `json:"protocol"`
	ProbeUrl    string `json:"probeUrl"`
}
