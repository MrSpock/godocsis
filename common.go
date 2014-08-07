// godocsis is package providing misceleaneous functions and few binaries
// that are usefull for users/admins DOCSIS based networks
// Currently all functions are cable modem related.
// TODO is to add support for concurency and some CMTS support
package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"net"
	"strconv"
	"time"
)

const (
	ResetOid               string = "1.3.6.1.2.1.69.1.1.3.0"
	DsOid                  string = "1.3.6.1.2.1.10.127.1.1.1.1.6"
	UsOid                  string = "1.3.6.1.2.1.10.127.1.2.2.1.3"
	IpAdEntIfIndex         string = "1.3.6.1.2.1.4.20.1.2"
	oid_cgConnectedDevices string = "1.3.6.1.4.1.2863.205.10.1.13"
)

var Session = &gosnmp.GoSNMP{
	Port:      161,
	Community: "public",
	Version:   gosnmp.Version2c,
	Timeout:   time.Duration(1) * time.Second,
	Retries:   2,
}

// cable modem structure
type CM struct {
	IPaddr   string
	RouterIP string
	RF       RFParams
	Devices  []cgConnectedDevices
}

// type for holding data related to customer devices connected to cable modem
type cgConnectedDevices struct {
	MacAddr       net.HardwareAddr
	Name          string
	IPAddr        net.IP
	InterfaceType int
}

// Basic structure used to hold misc CM RF parameters
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
//generic SNMPWalk function to walk over SNMP tree
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

// legacy helper function to convert []byte array to human readable form of IP
// currently this is handled by net/ip package functions
func HexIPtoString(octet_a []byte) (string, error) {
	if len(octet_a) == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", octet_a[0], octet_a[1], octet_a[2], octet_a[3]), nil
	} else {
		return "", fmt.Errorf("Unable to make conversion. 4 bytes required")
	}
}

// wrapper for commonly used code for handling errors
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
