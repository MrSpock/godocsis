// Package godocsis is package providing misceleaneous functions and few binaries
// that are usefull for users/admins DOCSIS based networks
// Currently all functions are cable modem related.
// TODO is to add support for concurency and some CMTS support
package godocsis

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/soniah/gosnmp"
)

// Protocol type used insted uint
type Protocol uint

type ModemStatus int

func (m ModemStatus) String() string {
	switch int(m) {
	case 1:
		return "other(1)"
	case 2:
		return "ranging(2)"
	case 3:
		return "rangindAborted(3)"
	case 4:
		return "rangingComplete(4)"
	case 5:
		return "ipComplete(5)"
	case 6:
		return "registrationComplete(6)"
	case 7:
		return "accessDenied(7)"
	case 8:
		return "operational(8)"
	case 9:
		return "registrationBPIInitializin(9)"
	default:
		return fmt.Sprintf("unknown(%d)", m)

	}

}

// IPAddrType SNMP ipaddr type
type IPAddrType int

func (r IPAddrType) String() string {
	if r == 1 {
		return "ipv4(1)"
	}
	return "unknown()"
}

func (p Protocol) String() string {
	return strconv.Itoa(int(p))
}

// Value return protocol type casted to int for external use
func (p Protocol) Value() int {
	return int(p)
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
	oid_tc7200_cgUiAdvancedForwardingRemoteIpAddr,
	oid_tc7200_cgUiAdvancedForwardingDescription,
}

const (
	// Both SNMP code
	Both Protocol = 1
	// Tcp SNMP code
	Tcp Protocol = 2
	// Udp SNMP code
	Udp Protocol = 3
)

const IPv4 IPAddrType = 1

var Session = gosnmp.GoSNMP{
	Port:      161,
	Community: "public",
	Version:   gosnmp.Version2c,
	Timeout:   time.Duration(2) * time.Second,
	Retries:   1,
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
	ExtIP               string
	RuleName            string
}

type CgForwardRule struct {
	LocalIP        net.IP
	LocalPortStart int
	LocalPortEnd   int
	ExtPortStart   int
	ExtPortEnd     int
	RuleName       string
	ProtocolType   Protocol
	IPAddrType     IPAddrType
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
		//fmt.Fprintf(os.Stderr, "WARN: LocalEndPort not set. Assuming the same value like LocalStartPort\n")
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
	if p.IPAddrType == 0 {
		p.IPAddrType = 1
		fmt.Fprintf(os.Stderr, "WARN: Default IP addr type used (ipv4(1))\n")
	}
	return
}

// Basic structure used to hold misc CM RF parameters
type RFParams struct {
	DSLevel []int
	USLevel []int
	// tens of dB
	USSNR int
}

func (rf *RFParams) DsBondingSize() int {
	return len(rf.DSLevel)

}

func (rf *RFParams) UsBondingSize() int {
	return len(rf.USLevel)
}

//type WalkFunc func(dataUnit gosnmp.SnmpPDU) []string, error)
//generic SNMPWalk function to walk over SNMP tree
func snmpwalk(session gosnmp.GoSNMP, oid string) ([]string, error) {

	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return nil, fmt.Errorf("snmpwalk() Connection error: %s", err)
	}
	response, err := session.WalkAll(oid) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return nil, fmt.Errorf("snmpwalk(): %s", err)
	}
	var result = make([]string, len(response))
	for i, pdu := range response {
		switch pdu.Type {
		case gosnmp.OctetString:
			result[i] = string(pdu.Value.([]uint8))
		case gosnmp.Integer:
			result[i] = strconv.Itoa(pdu.Value.(int))
			// case gosnmp.OctetString:
			// 	result[i] = string(pdu.Value.([]byte))
		}
	}
	return result, nil
}

func snmpset(session gosnmp.GoSNMP, pdus []gosnmp.SnmpPDU) (err error) {
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		return
	}
	_, err = session.Set(pdus)
	return
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

// PanicIf wrapper for commonly used code for handling errors
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// AddOidSuffix will add index suffix and return full OID
func AddOidSuffix(oid string, suffix int) (finalOid string) {
	data := []string{oid, strconv.Itoa(suffix)}
	finalOid = strings.Join(data, ".")
	return
}

//Oid2MAC Replace SNMP integer encoded MAC ("4.222.110.11.22.224")
// MAC address
func Oid2MAC(oid string) (string, error) {
	var rs []string
	for _, v := range strings.Split(oid, ".") {
		i, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		rs = append(rs, fmt.Sprintf("%02X", i))

	}
	return strings.Join(rs, ":"), nil
}
