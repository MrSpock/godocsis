package godocsis


// list of oids for forwarding table in TC7200
const (
	oid_tc7200_cgUiAdvancedForwardingPortStartValue         string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.2"
	oid_tc7200_cgUiAdvancedForwardingPortEndValue           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.3"
	oid_tc7200_cgUiAdvancedForwardingProtocolType           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.4"
	oid_tc7200_cgUiAdvancedForwardingIpAddrType             string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.5"
	oid_tc7200_cgUiAdvancedForwardingIpAddr                 string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.6"
	oid_tc7200_cgUiAdvancedForwardingEnabled                string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.7"
	oid_tc7200_cgUiAdvancedForwardingRowStatus              string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.8"
	oid_tc7200_cgUiAdvancedForwardingPortInternalStartValue string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.9"
	oid_tc7200_cgUiAdvancedForwardingPortInternalEndValue   string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.10"
	oid_tc7200_cgUiAdvancedForwardingRemoteIpAddr           string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.11"
	oid_tc7200_cgUiAdvancedForwardingDescription            string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.12"
	oid_tc7200_cgUiAdvancedForwardingRemove                 string = ".1.3.6.1.4.1.2863.205.10.1.33.2.5.1.13"
)


// cm upgrade
const (
	DocsDevSwServerOid       string = ".1.3.6.1.2.1.69.1.3.1.0"
	oid_docsDevSwFilename    string = ".1.3.6.1.2.1.69.1.3.2.0"
	oid_docsDevSwAdminStatus string = ".1.3.6.1.2.1.69.1.3.3.0"
	oid_docsDevSwCurrentVers string = ".1.3.6.1.2.1.69.1.3.5.0"
	oid_cmVersion            string = ".1.3.6.1.2.1.1.1.0"
	oid_cmLogs               string = ".1.3.6.1.2.1.69.1.5.8.1.7"
)

const (
	// ResetOid - generic DOCSIS cable modem reset oid
	ResetOid string = "1.3.6.1.2.1.69.1.1.3.0"
	// DsOid contans table of active downstreams
	DsOid string = "1.3.6.1.2.1.10.127.1.1.1.1.6"
	// UsOid contains table of used upstream channels
	UsOid string = "1.3.6.1.2.1.10.127.1.2.2.1.3"
	// IPAdEntIfIndex will provide tree with list of IP addressess
	IPAdEntIfIndex string = "1.3.6.1.2.1.4.20.1.2"
	// oid_cgConnectedDevices is Technicolor TC7200 specific mib
	// with list of connected devices
	oid_tc7200_cgConnectedDevices string = "1.3.6.1.4.1.2863.205.10.1.13"
)



//CMTS
const (
	// deprecated
	oid_docsIfCmtsCmStatusIpAddress = ".1.3.6.1.2.1.10.127.1.3.3.1.3"
	oid_docsIfCmtsCmStatus = ".1.3.6.1.2.1.10.127.1.3.3.1.9"
	oid_docsIfCmtsCmInetAddress = ".1.3.6.1.2.1.10.127.1.3.3.1.21"
	oid_docsIfCmtsCmMacAddress = ".1.3.6.1.2.1.10.127.1.3.3.1.2"
	oid_docsIfCmtsCmStatusSignalNoise = ".1.3.6.1.2.1.10.127.1.3.3.1.13"
)
