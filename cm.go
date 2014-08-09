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
		return fmt.Errorf("Unable to set reset OID (not cable modem)", err)
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
			mac_byte := []byte(pdu.Value.(string))
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
			devices[id-1].Name = pdu.Value.(string)
		case "4":
			devip_byte := []byte(pdu.Value.(string))
			devices[id-1].IPAddr = devip_byte

		}
	}

	// fmt.Println("Device list:")
	// for _, d := range devices {
	// 	fmt.Println("DevName:", d.Name, "MAC:", d.MacAddr, "IP:", d.IPAddr)
	// }
	return devices, nil
}