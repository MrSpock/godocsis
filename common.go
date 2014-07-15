package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"strconv"
	"time"
)

const (
	ResetOid       string = "1.3.6.1.2.1.69.1.1.3.0"
	DsOid          string = "1.3.6.1.2.1.10.127.1.1.1.1.6"
	UsOid          string = "1.3.6.1.2.1.10.127.1.2.2.1.3"
	IpAdEntIfIndex string = "1.3.6.1.2.1.4.20.1.2"
)

var Session = &gosnmp.GoSNMP{
	Port:      161,
	Community: "public",
	Version:   gosnmp.Version2c,
	Timeout:   time.Duration(1) * time.Second,
	Retries:   2,
}

type CM struct {
	IPaddr   string
	RouterIP string
	RF       RFParams
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
func snmpwalk(session *gosnmp.GoSNMP, oid string) ([]string, error) {

	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return nil, fmt.Errorf("Connection error", err)
	}
	response, err := session.WalkAll(oid) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return nil, fmt.Errorf("Walk error - no such mib ?", err)
	}
	var result = make([]string, len(response))
	for i, pdu := range response {
		result[i] = strconv.Itoa(pdu.Value.(int))
	}
	return result, nil
}
