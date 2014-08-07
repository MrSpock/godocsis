package godocsis

import (
	"errors"
	"fmt"
	"github.com/soniah/gosnmp"
	//"log"
	"strconv"
)

// Return Cable modem struct with filled fields related to CM RF parameters
// This will work on any cable modem since those are generic DOCSIS MIBS
func RFLevel(session *gosnmp.GoSNMP) (CM, error) {
	//Session.Target = ip
	var cm CM
	//var rfdata RFParams
	DSLevel, err := snmpwalk(Session, DsOid)
	if err != nil {
		fmt.Println("Error in RFLevel:", err)
		return cm, errors.New(err.Error())
	}

	cm.RF.DSLevel = string2int_a(DSLevel)
	USLevel, err := snmpwalk(Session, UsOid)
	if err != nil {
		return cm, fmt.Errorf("Problem with US level retrieval: %s", err)
	}
	cm.RF.USLevel = string2int_a(USLevel)
	//CM.RF = rfparams
	return cm, nil
}

// convert string to slice of integer values
func string2int_a(arstring []string) []int {
	rs := make([]int, len(arstring))
	for i, value := range arstring {
		rs[i], _ = strconv.Atoi(value)
	}
	return rs
}
