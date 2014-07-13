package rf


import (
    "github.com/alouca/gosnmp"
    "strconv"
    "fmt"
)


// errorString is a trivial implementation of error.
type errorString struct {
        s string
}

func (e *errorString) Error() string {
        return e.s
}
// New returns an error that formats as the given text.
func (e *errorString) New(text string) error {
        return &errorString{text}
}

// Struct and methods for each retured object by this module
type RFParams struct {
    DSLevel [] string
    USLevel [] int
}

func (rf *RFParams) DsBondingSize() int {
    return len(rf.DSLevel)
}

func (rf *RFParams) UsBondingSize() int {
        return len(rf.USLevel)
}


// DOCS-IF-MIB::docsIfDownChannelPower
const DsOid string = ".1.3.6.1.2.1.10.127.1.1.1.1.6"

func snmpwalk(ip string, oid string) ([]string,error) {
    var error errorString
    s, err := gosnmp.NewGoSNMP(ip, "public", gosnmp.Version2c, 5)
    if err != nil  {
        return nil, error.New("Error makeing SNMP connection")
    }
    resp, err := s.Walk(oid)
    if err != nil {
        return nil, error.New("Error getting Oid")
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
   var error errorString
   var rfdata RFParams
   DSLevel, err := snmpwalk(ip, DsOid)
   if err != nil {
       fmt.Println("Error in RFLevel:", err)
       return &rfdata, error.New(err.Error())
   }
   rfdata.DSLevel = DSLevel
   return &rfdata, nil
}

