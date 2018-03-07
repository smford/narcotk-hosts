package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Host struct {
	PaddedIP  string `json:"PaddedIP"`
	IPAddress string `json:"IPAddress"`
	Hostname  string `json:"Hostname"`
	Short1    string `json:"Short1"`
	Short2    string `json:"Short2"`
	Short3    string `json:"Short3"`
	Short4    string `json:"Short4"`
	MAC       string `json:"MAC"`
}

type Hosts struct {
	Hosts []Host `json:"Hosts"`
}

func PadLeft(str string) string {
	for {
		padding := "00"
		str = padding + str
		startpoint := len(str) - 3
		endpoint := len(str)
		return str[startpoint:endpoint]
	}
}

func makePaddedIp(ipaddress string) string {
	//fmt.Println("starting makePaddedIp")
	f := func(c rune) bool {
		return (c == rune('.'))
	}
	s := strings.FieldsFunc(ipaddress, f)
	paddedIp := PadLeft(s[0]) + PadLeft(s[1]) + PadLeft(s[2]) + PadLeft(s[3])
	fmt.Printf("P=%s\n", paddedIp)
	return paddedIp
}

func main() {
	var myhosts []Host
	myhosts = append(myhosts, Host{makePaddedIp("11.222.3.4"), "11.222.3.4", "server4.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("11.222.3.40"), "11.222.3.40", "server40.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("11.222.3.39"), "11.222.3.39", "server39.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("192.168.1.1"), "192.168.1.1", "server3.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("172.10.10.1"), "172.10.10.1", "server3.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("192.168.2.1"), "11.222.3.3", "server3.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("192.168.2.254"), "11.222.3.3", "server3.domain.com", "", "", "", "", ""})
	myhosts = append(myhosts, Host{makePaddedIp("192.168.2.2"), "11.222.3.3", "server3.domain.com", "", "", "", "", ""})
	fmt.Printf("      hosts=%q\n", myhosts)

	sort.Slice(myhosts, func(i, j int) bool {
		return bytes.Compare([]byte(myhosts[i].IPAddress), []byte(myhosts[j].IPAddress)) < 0
	})
	fmt.Printf("bad sorted=%q\n", myhosts)

	sort.Slice(myhosts, func(i, j int) bool {
		return bytes.Compare([]byte(myhosts[i].PaddedIP), []byte(myhosts[j].PaddedIP)) < 0
	})

	fmt.Printf("      good=%q\n", myhosts)

}
