package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"net"
	"strconv"
	"strings"
)

// Reset cable modem by setting docsDevResetNow.0 to one
// This will make cable modem to reinitialize itself
// The only param is cable modem IP address string
func ResetCm(host string) error {
	Session.Target = host
	Session.Community = "private"
	err := Session.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect:", err)
	}
	defer Session.Conn.Close()
	pdu := []gosnmp.SnmpPDU{gosnmp.SnmpPDU{ResetOid, gosnmp.Integer, 1}}
	//fmt.Println(pdu)
	_, err = Session.Set(pdu)
	if err != nil {
		return fmt.Errorf("Unable to set reset OID: ", err)
	}
	return nil
}

// For cable modems with Router built-in this will return e-Router CPE
// external IP address used for NAT all user traffic
func GetRouterIP(session *gosnmp.GoSNMP) (CM, error) {
	var cm CM
	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return cm, fmt.Errorf("Connection error", err)
	}
	response, err := session.WalkAll(IpAdEntIfIndex) //
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return cm, fmt.Errorf("Walk error - no such mib ?", err)
	}

	for _, pdu := range response {
		// For cablemodems I have ifIndex.1 contains embedded eRouter IP
		if pdu.Value.(int) == 1 {
			cm.RouterIP = strings.Trim(pdu.Name, IpAdEntIfIndex)
		}
	}
	return cm, nil
}

// For TC7200 cable modems this will return list of devicess connected to its LAN
// side. Both WiFi and Wired devices are listed
func GetConnetedDevices(session *gosnmp.GoSNMP) ([]cgConnectedDevices, error) {
	var mac_byte []byte
	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return nil, fmt.Errorf("Connection error", err)
	}
	response, err := session.WalkAll(oid_cgConnectedDevices)
	if err != nil {
		return nil, fmt.Errorf("ERR GetConnetedDevices()", err)

	}
	devices := make([]cgConnectedDevices, len(response)/4)
	for _, pdu := range response {
		//fmt.Println(pdu.Name, pdu.Type, pdu.Value)
		oid := strings.Split(pdu.Name, ".")
		id, _ := strconv.Atoi(oid[len(oid)-1])
		switch oid[len(oid)-2] {
		case "2":

			// upstream soniah/gosnmp changed OctetString Value from string to uint8
			mac_byte = pdu.Value.([]uint8)
			// either my o library fuckup. Some mac are inproperly decoded to string
			// I have to look into Decoding snmp library
			// this is workaround
			switch len(mac_byte) {
			case 6:
				//devices[id-1].MacAddr = fmt.Sprintf("%X:%X:%X:%X:%X:%X", mac_byte[0], mac_byte[1], mac_byte[2], mac_byte[2], mac_byte[4], mac_byte[5])
				devices[id-1].MacAddr = []byte(mac_byte)
				//fmt.Println(mac_byte, pdu.Value)
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
			devip_byte := pdu.Value.([]byte)
			devices[id-1].IPAddr = devip_byte

		}
	}

	// fmt.Println("Device list:")
	// for _, d := range devices {
	// 	fmt.Println("DevName:", d.Name, "MAC:", d.MacAddr, "IP:", d.IPAddr)
	// }
	return devices, nil
}

func CmGetNetiaPlayerList(session *gosnmp.GoSNMP) (npList []cgConnectedDevices, err error) {
	allDevices, err := GetConnetedDevices(session)
	if err != nil {
		return
	}
	for _, device := range allDevices {
		if strings.Contains(device.Name, "NetiaPlayer") {
			npList = append(npList, device)
		}
	}
	return
}

//request cable modem upgrade
//params are gosnmp.GoSNMP object, server IP address (string) and path relative to tftp root (string)
func CmUpgrade(session *gosnmp.GoSNMP, server string, filename string) (err error) {
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "private"
	}

	err = Session.Connect()
	defer Session.Conn.Close()
	if err != nil {
		return
	}
	// gosnmp supports only one set at request
	serverIP := net.ParseIP(server)
	pdu := make([]gosnmp.SnmpPDU, 1)
	pdu[0] = gosnmp.SnmpPDU{oid_docsDevSwServer, gosnmp.IPAddress, []byte(serverIP)[12:]}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}

	pdu[0] = gosnmp.SnmpPDU{oid_docsDevSwFilename, gosnmp.OctetString, filename}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	// docsDevSwAdminStatus.0 = 1 -> start upgrade (upgradeFromMgt(1))
	pdu[0] = gosnmp.SnmpPDU{oid_docsDevSwAdminStatus, gosnmp.Integer, 1}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	//fmt.Println(pdu)
	//_, err = session.Set(pdu)
	return
}

// type CgForwardingOid struct {
// 	ExtPortStart        string
// 	ExtPortEnd          string
// 	ProtocolType        string
// 	IpAddrType          string
// 	LocalIP             string
// 	ForwardingEnabled   string
// 	ForwardingRowStatus string
// 	LocalPortStart      string
// 	LocalPortEnd        string
// 	RuleName            string
// }
//
func CmGetFwdRuleCount(session *gosnmp.GoSNMP, oids *CgForwardingOid) (count int, err error) {
	count = 0
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "public"
	}
	err = session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return 0, fmt.Errorf("Connection error", err)
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

func CmSetForwardRule(session *gosnmp.GoSNMP, rule CgForwardRule, oids *CgForwardingOid) (err error) {
	err = rule.Validate()
	if err != nil {
		return
	}
	if len(session.Community) == 0 {
		// set default community if none is set
		session.Community = "private"
	}
	err = Session.Connect()
	defer Session.Conn.Close()
	if err != nil {
		return
	}
	ruleNo, err := CmGetFwdRuleCount(session, oids)
	pdu := make([]gosnmp.SnmpPDU, 1)
	//Ext Port Start
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ExtPortStart, ruleNo), gosnmp.Integer, rule.ExtPortStart}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	// Ext Port End
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.ExtPortEnd, ruleNo), gosnmp.Integer, rule.ExtPortEnd}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	//Local Port Start
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.LocalPortStart, ruleNo), gosnmp.Integer, rule.LocalPortStart}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}
	// Local Port End
	pdu[0] = gosnmp.SnmpPDU{AddOidSuffix(oids.LocalPortEnd, ruleNo), gosnmp.Integer, rule.LocalPortEnd}
	_, err = session.Set(pdu)
	if err != nil {
		return
	}

	return
}
