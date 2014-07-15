package godocsis

import (
	"github.com/soniah/gosnmp"
	"time"
)

const (
	ResetOid string = "1.3.6.1.2.1.69.1.1.3.0"
	DsOid    string = "1.3.6.1.2.1.10.127.1.1.1.1.6"
	UsOid    string = "1.3.6.1.2.1.10.127.1.2.2.1.3"
)

var session = &gosnmp.GoSNMP{
	Port:      161,
	Community: "public",
	Version:   gosnmp.Version2c,
	Timeout:   time.Duration(1) * time.Second,
	Retries:   2,
}

// basic structure used to hold misc CM RF parameters
type RFParams struct {
	DSLevel []int
	USLevel []int
}

func (rf *RFParams) DsBondingSize() int {
	return len(rf.DSLevel)

}

func (rf *RFParams) UsBondingSize() int {
	return len(rf.USLevel)
}

//type WalkFunc func(dataUnit gosnmp.SnmpPDU) []string, error)
