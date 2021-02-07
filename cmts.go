package godocsis

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
)

// GetModemList - provides list of modems seen on CMTS
func CmtsGetModemList(session gosnmp.GoSNMP) (CableModems, error) {
	cms := make(CableModems)
	err := session.Connect()
	defer session.Conn.Close()
	if err != nil {
		//log.Fatalf("Connect() err: %v", err)
		return cms, fmt.Errorf("Connection error: %s", err)
	}
	response, err := session.BulkWalkAll(oid_docsIfCmtsCmStatusTable)
	if err != nil {
		//log.Fatalf("Get() err: %v", err)
		return cms, fmt.Errorf("Walk error: %s", err)
	}
	//var ipAddrType int
	for _, pdu := range response {
		//fmt.Printf("%v -> %v\n ", pdu.Name, pdu.Value)
		oidArray := strings.Split(pdu.Name, ".")
		cmId, err := strconv.Atoi(oidArray[len(oidArray)-1])
		if err != nil {
			return cms, err
		}
		_, ok := cms[cmId]
		//fmt.Println("CMGET:", ok)
		if !ok {
			//log.Println("Creating CM with id", cmId)
			cms[cmId] = &CM{}
			//log.Println(cms)
		}
		//fmt.Println("CMS:", cms)
		cm := cms[cmId]
		//log.Println("CM", cm)
		cm.CmtsIndex = cmId
		if len(oidArray) > 1 {
			snmpColumn, err := strconv.Atoi(oidArray[len(oidArray)-2])
			if err != nil {
				continue
			}
			switch snmpColumn {
			// case CM_STATUS_INET_ADDRESS_TYPE:
			// 	ipAddrType = pdu.Value.(int)
			// case CM_STATUS_INET_ADDRESS:
			// 	if ipAddrType == IP_ADDR_TYPE_v4 {
			// 		fmt.Println("IPv4:", string(pdu.Value.([]byte)))
			// 	}
			case CM_STATUS_INET_ADDRESS:
				ip := pdu.Value.([]byte)
				cmIP := net.IPv4(ip[0], ip[1], ip[2], ip[3])
				// Huawei MA5633 returns IP as a string
				// C4 & E6000 returns byte[4]
				//fmt.Println("IPv4:", string(pdu.Value.([]byte)))
				//fmt.Println("IPv4:", cmIP.String())
				cm.IPaddr = cmIP.String()
			case CM_STATUS_MAC_ADDRESS:
				//fmt.Println("MAC: ", net.HardwareAddr(pdu.Value.([]byte)).String())
				cm.MacAddr = net.HardwareAddr(pdu.Value.([]byte)).String()
			case CM_STATUS_IP_ADDRES:
				// we got IP from CM_STATUS_INET_ADDRESS
				continue
			case CM_STATUS:
				modemState := ModemStatus(pdu.Value.(int))
				//fmt.Println("Modem status: ", modemState)
				cm.State = modemState
			case CM_STATUS_US_SNR:
				cm.RF.USSNR = pdu.Value.(int)
				//fmt.Printf("SNR: %0.1f dB\n", float32(pdu.Value.(int)/10))
			// case CM_STATUS_UPDATE_TS:
			// 	m, _ := time.ParseDuration(fmt.Sprintf("%ds", pdu.Value.(uint)))
			// 	fmt.Printf("Updated ago: %v", m)
			default:
				//fmt.Printf("%d -> %v\n", snmpColumn, pdu.Value)
			}
			//fmt.Printf("%v -> %v\n", snmpColumn, pdu.Value)
			//fmt.Printf("%+v\n", cm)
		}
		//fmt.Printf("PDU TYPE: %v\n", pdu.Value.(string))
		// switch t := pdu.Value.(type) {
		// case string:
		// 	cms = append(cms, t)
		// case []uint8:
		// 	//fmt.Println("UINT8")
		// 	cms = append(cms, string(t))
		// 	// case default:
		// 	// 	fmt.Println("PDU TYPE: %v", pdu.Value)
		// }
	}
	//fmt.Printf("%v\n", cms)
	return cms, nil
}
