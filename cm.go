package godocsis

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/soniah/gosnmp"
)

var logger = log.New(os.Stderr, "cm.go", 0)

// ResetCm function resets cable modem by setting docsDevResetNow.0 to one
// This will make cable modem to reinitialize itself
// The only param is cable modem IP address string
func ResetCm(host string, community string) error {
	Session.Target = host
	Session.Community = community
	err := Session.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect: %s", err)
	}
	defer Session.Conn.Close()
	pdu := []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: ResetOid, Type: gosnmp.Integer, Value: 1}}
	//fmt.Println(pdu)
	_, err = Session.Set(pdu)
	if err != nil {
		return fmt.Errorf("Unable to set reset OID: %s", err)
	}
	return nil
}

// GetRouterIP return built-in e-Router external (WAN) ip address
// used for NAT all user traffic
func GetRouterIP(session gosnmp.GoSNMP) (cm CM, err error) {
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return cm, fmt.Errorf("Connection error: %s", err)
	}
	response, err := session.WalkAll(IPAdEntIfIndex) //
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return cm, fmt.Errorf("Walk error: %s", err)
	}

	for _, pdu := range response {
		// For cablemodems I have ifIndex.1 contains embedded eRouter IP
		if pdu.Value.(int) == 1 {
			cm.RouterIP = strings.Replace(pdu.Name, "."+IPAdEntIfIndex+".", "", 1)
		}
	}
	return cm, nil
}

func GetLogs(session gosnmp.GoSNMP) (logs []string, err error) {
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return logs, fmt.Errorf("Connection error: %s", err)
	}
	response, err := session.WalkAll(oid_cmLogs) //
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return logs, fmt.Errorf("Walk error: %s", err)
	}

	for _, pdu := range response {
		//fmt.Printf("PDU TYPE: %v\n", pdu.Value.(string))
		switch t := pdu.Value.(type) {
		case string:
			logs = append(logs, t)
		case []uint8:
			//fmt.Println("UINT8")
			logs = append(logs, string(t))
			// case default:
			// 	fmt.Println("PDU TYPE: %v", pdu.Value)
		}
	}
	return logs, nil
}

// GetConnetedDevices show devicess connected to CM LAN
// side. Both WiFi and Wired devices are listed
func GetConnetedDevices(session gosnmp.GoSNMP) ([]cgConnectedDevices, error) {
	var macByte []byte
	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return nil, fmt.Errorf("ERR GetConnetedDevices() Connection error: %s", err)
	}
	response, err := session.WalkAll(oid_cgConnectedDevices)
	if err != nil {
		return nil, fmt.Errorf("ERR GetConnetedDevices() WalkAll(): %s\n", err)

	}
	devices := make([]cgConnectedDevices, len(response)/4)
	for _, pdu := range response {
		//fmt.Println(pdu.Name, pdu.Type, pdu.Value)
		oid := strings.Split(pdu.Name, ".")
		id, _ := strconv.Atoi(oid[len(oid)-1])
		switch oid[len(oid)-2] {
		case "2":

			// upstream soniah/gosnmp changed OctetString Value from string to uint8
			macByte = pdu.Value.([]uint8)
			// either my o library fuckup. Some mac are inproperly decoded to string
			// I have to look into Decoding snmp library
			// this is workaround
			switch len(macByte) {
			case 6:
				//devices[id-1].MacAddr = fmt.Sprintf("%X:%X:%X:%X:%X:%X", macByte[0], macByte[1], macByte[2], macByte[2], macByte[4], macByte[5])
				devices[id-1].MacAddr = []byte(macByte)
				//fmt.Println(macByte, pdu.Value)
			case 17:
				mac, err := net.ParseMAC(strings.Replace(pdu.Value.(string), " ", ":", -1))
				if err != nil {
					fmt.Errorf("ERR: MAC parse error", err)
				}
				devices[id-1].MacAddr = mac
			}
		case "3":

			devices[id-1].Name = string(pdu.Value.([]uint8))
		case "4":
			devipByte := pdu.Value.([]byte)
			devices[id-1].IPAddr = devipByte

		}
	}

	// fmt.Println("Device list:")
	// for _, d := range devices {
	// 	fmt.Println("DevName:", d.Name, "MAC:", d.MacAddr, "IP:", d.IPAddr)
	// }
	return devices, nil
}

// CmGetNetiaPlayerList is Netia specific function that return local side (LAN)
// IP address for every connected and detected NetiaPLayer STB.
// Each NetiaPlayer have it's DHCP client-identifier set to NetiaPlayer or NETGEM
func CmGetNetiaPlayerList(session gosnmp.GoSNMP) (npList []cgConnectedDevices, err error) {
	allDevices, err := GetConnetedDevices(session)
	if err != nil {
		return
	}
	for _, device := range allDevices {
		switch true {
		case strings.Contains(device.Name, "NetiaPlayer"):
			npList = append(npList, device)
		case strings.Contains(device.Name, "NETGEM"):
			npList = append(npList, device)
		}
		// if strings.Contains(device.Name, "NetiaPlayer") {
		// 	npList = append(npList, device)
		// }
	}
	return
}

//CmUpgrade request cable modem upgrade
//params are gosnmp.GoSNMP object, server IP address (string) and path relative to tftp root (string)
func CmUpgrade(session *gosnmp.GoSNMP, server string, filename string) (err error) {
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "private"
	}

	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		return
	}
	// gosnmp supports only one set at request
	serverIP := net.ParseIP(server)
	pdu := make([]gosnmp.SnmpPDU, 1)
	pdu[0] = gosnmp.SnmpPDU{Name: DocsDevSwServerOid, Type: gosnmp.IPAddress, Value: []byte(serverIP)[12:]}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}

	pdu[0] = gosnmp.SnmpPDU{oid_docsDevSwFilename, gosnmp.OctetString, filename, logger}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	// docsDevSwAdminStatus.0 = 1 -> start upgrade (upgradeFromMgt(1))
	pdu[0] = gosnmp.SnmpPDU{oid_docsDevSwAdminStatus, gosnmp.Integer, 1, logger}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}

	return
}

// CmGetFwdRuleCount return number of alredy present forwarding rules providing
// SNMP index where to pun new one without overrwriting exisiting one
func CmGetFwdRuleCount(session gosnmp.GoSNMP, oids *CgForwardingOid) (count int, err error) {
	count = 0
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "public"
	}
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return 0, fmt.Errorf("CmGetFwdRuleCount() Connection error: %s", err)
	}
	response, err := session.WalkAll(oids.ExtPortStart)
	for _, pdu := range response {
		//fmt.Println("Name:", pdu.Name, "Value:", pdu.Value)
		if pdu.Value != 0 {
			count++
		} else {
			//fmt.Println("OID:", AddOidSuffix(oids.ExtPortStart, count))
			return
		}
	}
	return 1, err
}

// CmSetForwardRule sets new firewall forwarding rule for TC7200 CM
func CmSetForwardRule(session gosnmp.GoSNMP, rule *CgForwardRule, oids *CgForwardingOid) (err error) {
	err = rule.Validate()
	if err != nil {
		return
	}
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "private"
	}
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		return
	}
	ruleNo, err := CmGetFwdRuleCount(session, oids)
	if ruleNo == 0 {
		ruleNo = 1
	}
	pdu := make([]gosnmp.SnmpPDU, 1)
	//fmt.Printf("SET ExtPortStart: %d, ", rule.ExtPortStart)
	fmt.Printf("SET: . ")
	// session.Connect()
	// defer session.Conn.Close()
	//Ext Port Start
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ExtPortStart, ruleNo), gosnmp.Integer, rule.ExtPortStart, logger}
	//err = snmpset(*session, pdu)
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	//fmt.Printf("ExtPortEnd: %d, ", rule.ExtPortEnd)
	fmt.Printf(". ")
	// Ext Port End
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ExtPortEnd, ruleNo), gosnmp.Integer, rule.ExtPortEnd, logger}
	//err = snmpset(*session, pdu)
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	//fmt.Printf("ProtocolType: %s, ", rule.ProtocolType)
	fmt.Printf(". ")
	// Ext Port End
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ProtocolType, ruleNo), gosnmp.Integer, int(rule.ProtocolType), logger}
	//err = snmpset(*session, pdu)
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	//fmt.Printf("IpAddrType: %s, ", rule.IPAddrType.String())
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.IpAddrType, ruleNo), gosnmp.Integer, int(rule.IPAddrType), logger}
	//_, err = session.Set(pdu)
	err = snmpset(session, pdu)
	if err != nil {
		return
	}
	//fmt.Printf("LocalIP: %s, ", rule.LocalIP)
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.LocalIP, ruleNo), gosnmp.OctetString, []byte(rule.LocalIP)[12:], logger}
	//_, err = session.Set(pdu)
	err = snmpset(session, pdu)
	if err != nil {
		return
	}
	//fmt.Printf("EnableRule: (1)")
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ForwardingEnabled, ruleNo), gosnmp.Integer, 1, logger}
	err = snmpset(session, pdu)
	if err != nil {
		return
	}

	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ForwardingRowStatus, ruleNo), gosnmp.Integer, 3, logger}
	err = snmpset(session, pdu)
	if err != nil {
		return
	}
	//Local Port Start
	//fmt.Printf("LocalPortStart: %d, ", rule.LocalPortStart)
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.LocalPortStart, ruleNo), gosnmp.Integer, rule.LocalPortStart, logger}
	err = snmpset(session, pdu)
	if err != nil {
		return
	}
	// Local Port End
	//fmt.Printf("LocalPortEnd: %d, ", rule.LocalPortEnd)
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.LocalPortEnd, ruleNo), gosnmp.Integer, rule.LocalPortEnd, logger}
	err = snmpset(session, pdu)
	if err != nil {
		return
	}

	/* To wygląda na nie wymagane */
	//fmt.Printf("ExtIP: %s, ", "0.0.0.0"
	fmt.Printf(". ")
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ExtIP, ruleNo), gosnmp.OctetString, []byte{0, 0, 0, 0}, logger}
	//_, err = session.Set(pdu)
	err = snmpset(session, pdu)
	if err != nil {
		return
	}

	//fmt.Printf("Description: %s, ", rule.RuleName)
	fmt.Printf(". ")
	//fmt.Println(oids.RuleName)
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.RuleName, ruleNo), gosnmp.OctetString, rule.RuleName, logger}
	err = snmpset(session, pdu)
	if err != nil {
		return
	}

	fmt.Println("")
	return
}
