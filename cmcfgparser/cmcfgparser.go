package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// TLVType is byte describing type of TLV
type TLVType byte

type TLVT struct {
	Value byte
	Name  string
}

// global type values

var EndOfDataMarker = TLVT{0xff, "EndOfDataMarker"}

// EndOfDataMarker This is a special marker for end of data.
const (

	//This has no length or value fields and is only used following the end of
	//data marker to pad the file to an integral number of 32-bit words.
	PadConfigSetting TLVType = 0x0

	//Downstream Frequency Configuration Setting
	DowstreamFrequency TLVType = 0x1

	//UpstreamChannelID
	UpstreamChannelID TLVType = 0x2

	//NetworkAccess
	//If the value field is a 1 this CM is allowed access to the network;
	//if a 0 it is not.
	NetworkAccess TLVType = 0x3

	// Class of Service Configuration Setting
	ClassOfService TLVType = 0x4

	//Modem capabilities settings
	ModemCap TLVType = 0x5

	// CM Message Integrity Check (MIC) Configuration Setting
	CMMIC TLVType = 0x6

	// CMTS Message Integrity Check (CMTS MIC)
	CMTSMIC TLVType = 0x7

	// Vendor ID settings
	VIDSettings TLVType = 0x8

	//SWUpgradeFilename
	SwUpgradeFilename TLVType = 0x9

	//SNMPWrite-AccessControl
	SnmpWriteAccess TLVType = 0xa

	//SNMPMIBObject
	SnmpMibObj TLVType = 0xb

	//Vendor Specific Information
	VSI TLVType = 0x2b

	// Modem IP
	CmIPAddr TLVType = 0xc

	// CPE MAC Addr
	CpeMacAddr TLVType = 0xe

	TriCfg01 TLVType = 0xf

	BpiSettings TLVType = 0x11

	MaxCPE TLVType = 0x12

	//C.7.21 TFTPServerTimestamp
	TftpServTs TLVType = 0x13

	//C.7.22 TFTP Server Provisioned Modem Addres
	TftpIPAddr TLVType = 0x14
)

//C.7.6.1 Internal Class of Service Encodings
const (
	ClassID    TLVType = 0x1
	DsMaxRate  TLVType = 0x2
	UsMaxRate  TLVType = 0x3
	UsChanPrio TLVType = 0x4
	MinCIR     TLVType = 0x5
	MaxUsBurst TLVType = 0x6
	CoSPrv     TLVType = 0x7
)

// TLV is basic pdu content
type TLV struct {
	Type TLVType
	//Length byte
	Value  interface{}
	Parent *TLV
	Childs []*TLV
}

func (t *TLV) Length() byte {
	switch v := t.Value.(type) {
	case uint8:
		return 1
	case uint32:
		return 4
	case uint64:
		return 8
	case string:
		return byte(len(v))
	}
	return 0
}
func MarshalTLV(tlv TLV) (result *bytes.Buffer, err error) {
	result = new(bytes.Buffer)
	//	result = bytes.NewBuffer([]byte{})
	switch v := tlv.Value.(type) {
	case uint8, uint32, uint64:
		err = binary.Write(result, binary.BigEndian, tlv.Type)
		if err != nil {
			panicIf("MarshalTLV() uint8 tlv.Type", err)
			return nil, err
		}
		err = binary.Write(result, binary.BigEndian, tlv.Length())
		if err != nil {
			panicIf("MarshalTLV() uint8 tlv.Length", err)
			return nil, err
		}
		err = binary.Write(result, binary.BigEndian, v)
		if err != nil {
			panicIf("MarshalTLV() uint8 v", err)
			return nil, err
		}
		return
	case string:
		err = binary.Write(result, binary.BigEndian, tlv.Type)
		if err != nil {
			return
		}
		err = binary.Write(result, binary.BigEndian, tlv.Length())
		if err != nil {
			return
		}
		err = binary.Write(result, binary.BigEndian, []byte(v))
		if err != nil {
			return
		}
		return
	}

	return
}

func main() {
	var fpdu bytes.Buffer
	t := TLV{DowstreamFrequency, uint32(402000000), nil, nil}
	filename := "upgrade"
	t2 := TLV{SwUpgradeFilename, filename, nil, nil}
	t3 := TLV{NetworkAccess, byte(1), nil, nil}
	pdu, err := MarshalTLV(t)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fpdu.Write(pdu.Bytes())

	pdu, err = MarshalTLV(t2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fpdu.Write(pdu.Bytes())
	pdu, _ = MarshalTLV(t3)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fpdu.Write(pdu.Bytes())
	outFile, err := os.Create("config.cm")
	defer outFile.Close()
	if err != nil {
		panicIf(err)
		os.Exit(1)
	}
	outFile.Write(fpdu.Bytes())
	//fmt.Printf("%v", fpdu)
}

func panicIf(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg)
}
