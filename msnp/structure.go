package msnp

import "net"

type msnp_command struct {
	command       string
	transactionID string
	arguments     string
}

type Msnp_Client struct {
	Connection net.Conn
}
