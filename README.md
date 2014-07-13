go-docsis
=========
Very rudimental implementation of basic CM queries for larger project.
This is very fresh topic. Look here in future :)

example:
```
package main
import (
    "github.com/MrSpock/godocsis"
    "fmt"
       )



func main() {
        rfdata,err := rf.RFLevel("10.80.0.164")
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println("DS level", rfdata.DSLevel, rfdata.DsBondingSize())
}
```

