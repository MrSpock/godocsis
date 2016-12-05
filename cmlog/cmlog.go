package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrspock/godocsis"
)

var (
	community = flag.String("community", "public", "RW community to use when sending restart request")
	csvmode   = flag.Bool("csv", false, "CSV mode for easy import to spreadsheet")
)

func printVerbose(cmd godocsis.CM) {
	fmt.Printf("%s ", cmd.IPaddr)
	fmt.Printf("US(dBmV):%.01f ", float32(cmd.RF.USLevel[0])/10)
	separator := ","
	fmt.Printf("DS(dBmV):")
	for no, ds := range cmd.RF.DSLevel {
		if no == cmd.RF.DsBondingSize()-1 {
			separator = ""
		}
		fmt.Printf("%.01f%v", float32(ds)/10, separator)
	}
	fmt.Println("")

}
func printCSV(cmd godocsis.CM) {
	fmt.Printf("%s:", cmd.IPaddr)
	fmt.Printf("%.01f:", float32(cmd.RF.USLevel[0])/10)
	separator := ","
	for no, ds := range cmd.RF.DSLevel {
		if no == cmd.RF.DsBondingSize()-1 {
			separator = ""
		}
		fmt.Printf("%.01f%v", float32(ds)/10, separator)
	}
	fmt.Println("")

}
func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: cmparams [--csv] [--community <community>] <ip> <ip>")
		return
	}

	s := godocsis.Session
	s.Community = *community
	for _, ip := range flag.Args() {
		s.Target = ip
		rs, err := godocsis.RFLevel(s)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Problem: %v", err)
			fmt.Fprintf(os.Stderr, "%s:OFFLINE\n", ip)
			//panic(err)
		} else {
			if *csvmode {
				printCSV(rs)
			} else {
				printVerbose(rs)
			}
			//			fmt.Printf("%s:", ip)
			//			fmt.Printf("%.01f:", float32(rs.RF.USLevel[0])/10)
			//			separator := ","
			//			for no, ds := range rs.RF.DSLevel {
			//				if no == rs.RF.DsBondingSize()-1 {
			//					separator = ""
			//				}
			//				fmt.Printf("%.01f%v", float32(ds)/10, separator)
			//			}
			//			fmt.Println("")
			//
		}
	}

}
