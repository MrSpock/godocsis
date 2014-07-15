go-docsis
=========
Very rudimental implementation of basic CM queries for larger project.
This is very fresh topic. Look here in future :)
Code is not using concurrency at the moment.

example:
```
package main

import (
    "flag"
    "fmt"
    "github.com/mrspock/godocsis"
)

func main() {
    var ip string
    flag.Parse()
    if len(flag.Args()) < 1 {
        // default IP - my default test modem  (3dc1)
        ip = "10.80.0.164"
    }
    ip = flag.Args()[0]
    rs, err := rf.RFLevel(ip)
    if err != nil {
        fmt.Println("Error")
        panic(err)
    }

    fmt.Println("DS", rs.DSLevel, "\nDS Bonding size:", rs.DsBondingSize())
    fmt.Println("US", rs.USLevel)
    // new method (will be aplied to RFLevel())
    s := godocsis.Session
    s.Target = ip
    cm, err := godocsis.GetRouterIP(s)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(cm.RouterIP)
}

```

