package main

import (
	"flag"
	"fmt"
	"github.com/mrspock/godocsis"
)

func main() {
	var ip string
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Missing argument - cm ip address.")
		return
	} else {
		ip = flag.Args()[0]
	}

	s := godocsis.Session
	s.Target = ip
	fmt.Println(s.Target, "device list:")
	devices, err := godocsis.GetConnetedDevices(s)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("IP Address\tMac Address\t\tHostname")
	for _, device := range devices {
		fmt.Println(device.IPAddr.String() + "\t" + device.MacAddr.String() + "\t" + device.Name)
	}
}
