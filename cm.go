package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
)

func ResetCm(host string) error {
	session.Target = host
	session.Community = "private"
	err := session.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect:", err)
	}
	defer session.Conn.Close()
	pdu := []gosnmp.SnmpPDU{gosnmp.SnmpPDU{ResetOid, gosnmp.Integer, 1}}
	//fmt.Println(pdu)
	_, err = session.Set(pdu)
	if err != nil {
		return fmt.Errorf("Unable to set reset OID (not cable modem)", err)
	}
	//fmt.Println(result)

	// for i, variable := range result.Variables {
	// 	fmt.Printf("%d: oid: %s ", i, variable.Name)
	// }
	return nil
}
