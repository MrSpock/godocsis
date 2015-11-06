package main

import (
	"flag"
	"fmt"
	"github.com/mrspock/godocsis"
	"net"
	"os"
)

const (
	VERSION string = "1.0.4"
)

var Usage = func() {
	Help(os.Args[0])
	flag.PrintDefaults()
}

func main() {
	//var ip string
	community := flag.String("community", "private", "RW community to use when sending restart request")
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) < 2 {
		Help(os.Args[0])
		return
	}
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
			err := godocsis.CmUpgrade(&snmp, server, path)
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
