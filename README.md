go-docsis
=========
Implementation of basic CM oprations. Part of larger closed source project.
Code is not using concurrency at the moment.

Example that will fetch DS/US RF levels from CM:
```
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrspock/godocsis"
)

func main() {
	//var ip string
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: cmparams <ip> <ip>")
		return
	}
	s := godocsis.Session
	for _, ip := range flag.Args() {
		s.Target = ip
		rs, err := godocsis.RFLevel(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem: %v", err)
			//panic(err)
		} else {
			fmt.Printf("%s:", ip)
			fmt.Printf("%.01f:", float32(rs.RF.USLevel[0])/10)
			separator := ","
			for no, ds := range rs.RF.DSLevel {
				if no == rs.RF.DsBondingSize()-1 {
					separator = ""
				}
				fmt.Printf("%.01f%v", float32(ds)/10, separator)
			}
			fmt.Println("")
		}
	}

}
```

