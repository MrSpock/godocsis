package godocsis

import (
	"fmt"
	"github.com/soniah/gosnmp"
)



func ResetCm(host string) error {
	gosnmp.Default.Target = host
	gosnmp.Default.Community = "private"
	err := gosnmp.Default.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect:", err)
	}
	defer gosnmp.Default.Conn.Close()
	pdu := []gosnmp.SnmpPDU{gosnmp.SnmpPDU{ResetOid, gosnmp.Integer, 1}}
	//fmt.Println(pdu)
	_, err = gosnmp.Default.Set(pdu)
	if err != nil {
		return fmt.Errorf("Unable to set reset OID (not cable modem)", err)
	}
	//fmt.Println(result)

	// for i, variable := range result.Variables {
	// 	fmt.Printf("%d: oid: %s ", i, variable.Name)
	// }
	return nil
}
