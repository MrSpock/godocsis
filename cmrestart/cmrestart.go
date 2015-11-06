package main

import (
	"flag"
	"fmt"
	//	"net"
	"os"

	"github.com/mrspock/godocsis"
)

const (
	VERSION string = "1.0.3"
)

var Usage = func() {
	Help(os.Args[0])
	flag.PrintDefaults()
}

func main() {
	community := flag.String("community", "private", "RW community to use when sending restart request")
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) == 0 {
		Help(os.Args[0])
		return
	}
	for _, address := range flag.Args() {
		//		ip, err := net.LookupHost(host)
		//		if err != nil {
		//			fmt.Fprintln(os.Stderr, "Host lookup error:", err)
		//os.Exit(1)
		//			continue
		//		}
		//		for _, address := range flag.Args() {
		//fmt.Println(address)
		err := godocsis.ResetCm(address, *community)
		if err != nil {
			fmt.Fprintln(os.Stderr, "NG: Wystąpił błąd komunikacji z modemem", address, ":", err)
			//os.Exit(1)
			continue
		} else {
			fmt.Fprintln(os.Stdout, "OK: Modem", address, "w trakcie restartu..")

			//os.Exit(0)
		}
		//		}
	}
	os.Exit(0)
}

func Help(name string) {
	fmt.Fprintf(os.Stderr, "======= Cable Modem restarter by Spock (BSD) ver %s ========\nUsage: %s cm1_ipaddr [cm2_ipaddr cmN_ipaddr]\n============================================\nOptions:\n", VERSION, name)
}
