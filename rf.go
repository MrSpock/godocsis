package godocsis

import (
	"errors"
	"fmt"
	"github.com/alouca/gosnmp"
	"strconv"
)

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

//const ResetOid string = ".1.3.6.1.2.1.69.1.1.3.0"

func snmpwalk(ip string, oid string) ([]string, error) {
	s, err := gosnmp.NewGoSNMP(ip, "public", gosnmp.Version2c, 5)
	if err != nil {
		return nil, errors.New("Error makeing SNMP connection")
	}
	resp, err := s.Walk(oid)
	if err != nil {
		return nil, errors.New("Error getting Oid")
	}
	var result = make([]string, len(resp))
	for i, pdu := range resp {
		//switch pdu.Value {
		//case gosnmp.OctetString:
		//case gosnmp.Integer:
		//result[i] = strconv.Itoa(pdu.Value)
		responseValue := pdu.Value.(int)
		//fmt.Println("Index:", i, ",Value:", responseValue)
		result[i] = strconv.Itoa(responseValue)

		//}
	}
	return result, nil
}

func RFLevel(ip string) (*RFParams, error) {

	var rfdata RFParams
	DSLevel, err := snmpwalk(ip, DsOid)
	if err != nil {
		fmt.Println("Error in RFLevel:", err)
		return &rfdata, errors.New(err.Error())
	}

	rfdata.DSLevel = string2int_a(DSLevel)
	USLevel, err := snmpwalk(ip, UsOid)
	if err != nil {
		return &rfdata, fmt.Errorf("Problem with US level retrieval: %s", err)
	}
	rfdata.USLevel = string2int_a(USLevel)
	return &rfdata, nil
}

// helpers
func string2int_a(arstring []string) []int {
	rs := make([]int, len(arstring))
	for i, value := range arstring {
		rs[i], _ = strconv.Atoi(value)
	}
	return rs
}
