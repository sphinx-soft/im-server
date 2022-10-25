package msnp

import "net"

type Msnp_Client struct {
	Connection  net.Conn
	Dispatched  bool
	Account     Msnp_Account
	BuildNumber string
}

type Msnp_Account struct {
	Uid        int
	Email      string
	Password   string
	Screenname string
}

var Msnp_Clients []*Msnp_Client
