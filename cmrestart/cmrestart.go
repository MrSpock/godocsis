package main

import (
	"flag"
	"fmt"
	"github.com/mrspock/godocsis"
	"os"
)

func main() {
	var ip string
	flag.Parse()
	if len(flag.Args()) == 0 {
		Help(os.Args[0])
		return
	} else {
		ip = flag.Args()[0]
	}
	err := godocsis.ResetCm(ip)
	if err != nil {
		fmt.Println("NG: Wystąpił błąd komunikacji z modemem:", err)
		os.Exit(1)
	} else {
		fmt.Println("OK: Modem w trakcie restartu..")
		os.Exit(0)
	}

}

func Help(name string) {
	fmt.Fprintf(os.Stderr, "======= Cable Modem restarter by Spock (BSD) ========\nUsage: %s cm_ipaddr\n============================================\n", name)
}
