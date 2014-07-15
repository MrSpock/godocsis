package godocsis

import (
	"errors"
	"fmt"
	"github.com/soniah/gosnmp"
	"log"
	"strconv"
)

func snmpwalk(session *gosnmp.GoSNMP, oid string) ([]string, error) {

	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
		return nil, fmt.Errorf("Connection error", err)
	}
	response, err := session.WalkAll(oid) // Get() accepts up to g.MAX_OIDS
	if err != nil {
		log.Fatalf("Get() err: %v", err)
		return nil, fmt.Errorf("Walk error - no such mib ?", err)
	}
	var result = make([]string, len(response))
	for i, pdu := range response {
		result[i] = strconv.Itoa(pdu.Value.(int))
	}
	return result, nil
}

func RFLevel(ip string) (*RFParams, error) {
	session.Target = ip
	var rfdata RFParams
	DSLevel, err := snmpwalk(session, DsOid)
	if err != nil {
		fmt.Println("Error in RFLevel:", err)
		return &rfdata, errors.New(err.Error())
	}

	rfdata.DSLevel = string2int_a(DSLevel)
	USLevel, err := snmpwalk(session, UsOid)
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
