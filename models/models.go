package models

type TargetMAC struct {
	ID         int
	MACAddress string
}

type TargetIP struct {
	ID        int
	IPAddress string
}

type PowerShellCredential struct {
	ID           int
	ComputerName string
	Username     string
	Password     string
}
