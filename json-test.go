package main

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Hostname  string `json:"Hostname"`
	IPAddress string `json:"IPAddress"`
}

type Network struct {
	Network string `json:"Network"`
	CIDR    string `json:"CIDR"`
	Hosts   []Host `json:"Hosts"`
}

type Networks struct {
	Networks []Network `json:"Networks"`
}

func main() {
	host11 := Host{"host1-1.something.com", "192.168.10.1"}
	host12 := Host{"host1-2.something.com", "192.168.10.2"}
	host13 := Host{"host1-3.something.com", "192.168.10.3"}
	host14 := Host{"host1-4.something.com", "192.168.10.4"}
	host21 := Host{"host2-1.example.org", "10.0.0.1"}
	host22 := Host{"host2-2.example.org", "10.0.0.2"}
	host23 := Host{"host2-3.example.org", "10.0.0.3"}

	network1 := Network{"192.168.10", "192.168.10.0/24", []Host{host11, host12, host13, host14}}

	network2 := Network{"10.0.0", "10.0.0.0/24", []Host{host21, host22, host23}}

	allNetworks := &Networks{
		Networks: []Network{network1, network2},
	}

	c, err := json.Marshal(allNetworks)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Printf("%s", c)
}
