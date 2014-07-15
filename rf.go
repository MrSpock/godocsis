package godocsis

import (
	"errors"
	"fmt"
	//"github.com/soniah/gosnmp"
	//"log"
	"strconv"
)

func RFLevel(ip string) (*RFParams, error) {
	Session.Target = ip
	//var cm CM
	var rfdata RFParams
	DSLevel, err := snmpwalk(Session, DsOid)
	if err != nil {
		fmt.Println("Error in RFLevel:", err)
		return &rfdata, errors.New(err.Error())
	}

	rfdata.DSLevel = string2int_a(DSLevel)
	USLevel, err := snmpwalk(Session, UsOid)
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
