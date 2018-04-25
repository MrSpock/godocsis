package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mrspock/godocsis"
)

var (
	community = flag.String("community", "public", "RW community to use when sending restart request")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: cmts [-community <community>] <ip>")
		return
	}
	log.Println("CMTS IP: ", flag.Arg(0))
	s := godocsis.Session
	s.Community = *community
	s.Target = flag.Arg(0)
	cableModems, err := godocsis.CmtsGetModemList(s)
	if err != nil {
		log.Println("CMTS query error:", err)
	}
	fmt.Printf("Total (%d) modem list for CMTS %s\n", len(cableModems), s.Target)
	fmt.Println("ID\tMAC\t\t\tIPADDR\t\tState\t\tUS_SNR")
	for _, cm := range cableModems {
		fmt.Printf("%d\t%s\t%s\t%s\t\t%.01fdB\n", cm.CmtsIndex, cm.MacAddr, cm.IPaddr, cm.State, float32(cm.RF.USSNR/10))
	}
}
