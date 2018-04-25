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
func RFLevel(session gosnmp.GoSNMP) (CmtsCM, error) {
	//Session.Target = ip
	var cm CmtsCM
	cm.IPaddr = session.Target
	//var rfdata RFParams
	DSLevel, err := snmpwalk(session, DsOid)
	if err != nil {
		fmt.Println("Error in RFLevel:", err)
		return cm, errors.New(err.Error())
	}

	cm.RF.DSLevel = string2int_a(DSLevel)
	USLevel, err := snmpwalk(session, UsOid)
	if err != nil {
		return cm, fmt.Errorf("Problem with US level retrieval: %s", err)
	}
	cm.RF.USLevel = string2int_a(USLevel)
	//CM.RF = rfparams
	return cm, nil
}

func CmVersion(session gosnmp.GoSNMP) (version string, err error) {
	rs, err := snmpwalk(session, oid_cmVersion)

	if err != nil {
		fmt.Println("Error in CmVersion:", err)
		return version, errors.New(err.Error())
	}
	if len(rs) == 1 {
		version = rs[0]

	} else {
		return version, errors.New("Wrong number of returned varbinds - expected one")
	}
	return version, nil
}

// convert string to slice of integer values
func string2int_a(arstring []string) []int {
	rs := make([]int, len(arstring))
	for i, value := range arstring {
		rs[i], _ = strconv.Atoi(value)
	}
	return rs
}
