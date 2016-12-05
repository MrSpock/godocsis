package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrspock/godocsis"
)

var (
	community = flag.String("community", "public", "RW community to use when sending restart request")
	jsonout   = flag.Bool("json", false, "CSV mode for easy import to spreadsheet")
)

type CmData struct {
	Cm_ip           string `json:"cm_ip"`
	Sys_description string `json:"sysdesc"`
}

func main() {
	var out_a []string
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: cmversion [-json] [-community <community>] <ip> <ip>")
		return
	}

	s := godocsis.Session
	s.Community = *community
	for _, ip := range flag.Args() {
		s.Target = ip
		rs, err := godocsis.CmVersion(s)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Problem: %v", err)
			fmt.Fprintf(os.Stderr, "%s:OFFLINE\n", ip)
			//panic(err)
		} else {
			if *jsonout {
				out_a = append(out_a, fmt.Sprintf("{\"ip\": \"%s\", \"sysdesc\": \"%s\"}", ip, rs))

			} else {
				fmt.Println(rs)
			}

		}
	}
	if *jsonout {
		fmt.Println("[")
		fmt.Println(strings.Join(out_a, ","))
		fmt.Println("]")
	}
}
