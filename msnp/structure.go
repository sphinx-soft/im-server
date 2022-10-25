package msnp

import "net"

type Msnp_Client struct {
	Connection  net.Conn
	Dispatched  bool
	Account     msnp_account
	BuildNumber string
}

type msnp_account struct {
	Uid        int
	Email      string
	Password   string
	Screenname string
}

var Msnp_Clients []*Msnp_Client
