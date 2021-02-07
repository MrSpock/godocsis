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
	csvmode   = flag.Bool("json", false, "JSON mode for easy import")
)

func printJSON(data []string) {
	var rs []string
	fmt.Println("{ \"log\": [")
	for i, v := range data {
		rs = append(rs, fmt.Sprintf("\n{\n\t\"id\": \"%d\",\n\t\"message\": \"%s\"\n}", i, v))
	}
	//fmt.Println("{")
	fmt.Println(strings.Join(rs, ","))
	fmt.Println("]}")
}
func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: cmlog [--json] [--community <community>] ip")
		return
	}

	s := godocsis.Session
	s.Community = *community
	for _, ip := range flag.Args() {
		s.Target = ip
		rs, err := godocsis.GetLogs(s)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Problem: %v", err)
			fmt.Fprintf(os.Stderr, "%s:OFFLINE\n", ip)
			//panic(err)
		} else {
			if *csvmode {
				printJSON(rs)

			} else {
				//printVerbose(rs)
				for i, l := range rs {
					fmt.Printf("%d: %s\n", i, l)

				}
			}

		}
	}

}
