package main

// Configuration File Opjects
type configuration struct {
	ProjectID      string
	AppVersion   string
	ServerName   string
	Broker       string
	BrokerUser   string
	BrokerPwd    string
	BrokerVhost  string
	BrokerQueue  string
	ChannelCount int
	DstURL       string
	LocalEcho    bool
	SysLogSrv    string
	SysLogPort   string
	LogLevel     int
}
