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
	"strings"
	"time"
)

type Protocol uint

func (p Protocol) String() string {
	return strconv.Itoa(int(p))
}

func (p Protocol) Value() int {
	return int(p)
}

const (
	ResetOid               string = "1.3.6.1.2.1.69.1.1.3.0"
	DsOid                  string = "1.3.6.1.2.1.10.127.1.1.1.1.6"
	UsOid                  string = "1.3.6.1.2.1.10.127.1.2.2.1.3"
	IpAdEntIfIndex         string = "1.3.6.1.2.1.4.20.1.2"
	oid_cgConnectedDevices string = "1.3.6.1.4.1.2863.205.10.1.13"
)

const (
	oid_docsDevSwServer      string = ".1.3.6.1.2.1.69.1.3.1.0"
	oid_docsDevSwFilename    string = ".1.3.6.1.2.1.69.1.3.2.0"
	oid_docsDevSwAdminStatus string = ".1.3.6.1.2.1.69.1.3.3.0"
)

// list of oids for forwarding table in TC7200
const (
	oid_tc7200_cgUiAdvancedForwardingPortStartValue         string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.2"
	oid_tc7200_cgUiAdvancedForwardingPortEndValue           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.3"
	oid_tc7200_cgUiAdvancedForwardingProtocolType           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.4"
	oid_tc7200_cgUiAdvancedForwardingIpAddrType             string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.5"
	oid_tc7200_cgUiAdvancedForwardingIpAddr                 string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.6"
	oid_tc7200_cgUiAdvancedForwardingEnabled                string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.7"
	oid_tc7200_cgUiAdvancedForwardingRowStatus              string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.8"
	oid_tc7200_cgUiAdvancedForwardingPortInternalStartValue string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.9"
	oid_tc7200_cgUiAdvancedForwardingPortInternalEndValue   string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.10"
	oid_tc7200_cgUiAdvancedForwardingRemoteIpAddr           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.11"
	oid_tc7200_cgUiAdvancedForwardingDescription            string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.12"
)
const (
	Both Protocol = 1
	Tcp  Protocol = 2
	Udp  Protocol = 3
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

type CgForwardingOid struct {
	ExtPortStart        string
	ExtPortEnd          string
	ProtocolType        string
	IpAddrType          string
	LocalIP             string
	ForwardingEnabled   string
	ForwardingRowStatus string
	LocalPortStart      string
	LocalPortEnd        string
	RuleName            string
}

var TC7200ForwardingTree = &CgForwardingOid{
	oid_tc7200_cgUiAdvancedForwardingPortStartValue,
	oid_tc7200_cgUiAdvancedForwardingPortEndValue,
	oid_tc7200_cgUiAdvancedForwardingProtocolType,
	oid_tc7200_cgUiAdvancedForwardingIpAddrType,
	oid_tc7200_cgUiAdvancedForwardingIpAddr,
	oid_tc7200_cgUiAdvancedForwardingEnabled,
	oid_tc7200_cgUiAdvancedForwardingRowStatus,
	oid_tc7200_cgUiAdvancedForwardingPortInternalStartValue,
	oid_tc7200_cgUiAdvancedForwardingPortInternalEndValue,
	oid_tc7200_cgUiAdvancedForwardingDescription,
}

type CgForwardRule struct {
	LocalIP        net.IP
	LocalPortStart uint
	LocalPortEnd   uint
	ExtPortStart   uint
	ExtPortEnd     uint
	RuleName       string
	ProtocolType   Protocol
}

func (p *CgForwardRule) Validate() (err error) {
	if p.LocalPortStart == 0 {
		err = fmt.Errorf("LocalPortStart can't be 0")
	}
	if p.ExtPortStart == 0 {
		err = fmt.Errorf("ExtPortStart can't be 0")
	}
	if p.LocalPortEnd == 0 {
		p.LocalPortEnd = p.LocalPortStart
	}
	if p.ExtPortEnd == 0 {
		p.ExtPortEnd = p.ExtPortStart
	}
	if p.LocalPortStart > p.LocalPortEnd {
		err = fmt.Errorf("LocalPortStart can't be higher than LocalPortEnd")
	}
	if p.ExtPortStart > p.ExtPortEnd {
		err = fmt.Errorf("ExtPortStart can't be higher than ExtPortEnd")
	}
	if len(p.RuleName) == 0 {
		err = fmt.Errorf("RuleName can't be empty")
	}
	if p.ProtocolType == 0 {
		p.ProtocolType = Both
	}
	return
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

func AddOidSuffix(oid string, suffix int) (finalOid string) {
	data := []string{oid, strconv.Itoa(suffix)}
	finalOid = strings.Join(data, ".")
	return
}
