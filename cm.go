package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
	"strings"
)

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
			cm.IPaddr = strings.Trim(pdu.Name, IpAdEntIfIndex)
		}
	}
	return cm, nil
}
