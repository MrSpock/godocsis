package main

import (
	"flag"
	"fmt"
	"github.com/mrspock/godocsis"
	"net"
	"os"
)

func main() {
	//var ip string
	flag.Parse()
	if len(flag.Args()) == 0 {
		Help(os.Args[0])
		return
	}
	for _, host := range flag.Args() {
		ip, err := net.LookupHost(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Host lookup error:", err)
			//os.Exit(1)
			continue
		}
		for _, address := range ip {
			//fmt.Println(address)
			err := godocsis.ResetCm(address)
			if err != nil {
				fmt.Fprintln(os.Stderr, "NG: Wystąpił błąd komunikacji z modemem", address, ":", err)
				//os.Exit(1)
				continue
			} else {
				fmt.Fprintln(os.Stdout, "OK: Modem", address, "w trakcie restartu..")
				//os.Exit(0)
			}
		}
	}
	os.Exit(1)
}

func Help(name string) {
	fmt.Fprintf(os.Stderr, "======= Cable Modem restarter by Spock (BSD) ========\nUsage: %s cm1_ipaddr cm2_ipaddr\n============================================\n", name)
}
