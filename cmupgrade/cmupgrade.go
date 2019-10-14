package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/mrspock/godocsis"
)

const (
	VERSION string = "1.0.5"
)

var Usage = func() {
	Help(os.Args[0])
	flag.PrintDefaults()
}

// cable modem upgrade protocol

var CmUpgradeProtocol = map[string]int{
	"tftp": 1,
	"http": 2,
}
var upgradeProto godocsis.CmProtocol

func main() {
	//var ip string
	community := flag.String("community", "private", "RW community to use when sending restart request")
	protocol := flag.String("protocol", "tftp", "Upgrade method protocol [tftp],http")
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) < 2 {
		Help(os.Args[0])
		return
	}
	//var upgradeProto int
	if up, exist := CmUpgradeProtocol[*protocol]; !exist {
		//upgradeProto = godocsis.CmProtocol(up)
		fmt.Println("Upgrade method doesn't exist. Use tftp or http")
		os.Exit(-1)
	} else {
		upgradeProto = godocsis.CmProtocol(up)
	}
	fmt.Println("Upgrade proto: ", upgradeProto)
	server := flag.Arg(0)
	path := flag.Arg(1)
	for _, host := range flag.Args()[2:] {

		ip, err := net.LookupHost(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Host lookup error:", err)
			//os.Exit(1)
			continue
		}
		for _, address := range ip {
			//fmt.Println(address)
			snmp := godocsis.Session
			snmp.Community = *community
			snmp.Target = address
			err := godocsis.CmUpgrade(&snmp, server, path, upgradeProto)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				//os.Exit(1)
				continue
			} else {
				fmt.Fprintln(os.Stdout, "OK: Upgrade", address, "in progress.")
				//os.Exit(0)
			}
		}
	}
	os.Exit(0)
}

func Help(name string) {
	fmt.Fprintf(os.Stderr, "======= Cable Modem upgrader by Spock (BSD) ver %s ========\nUsage: %s TFTP_SERVER SW_PATH(relative to tftp root) cm1_ipaddr cm2_ipaddr\n============================================\n", VERSION, name)
}
