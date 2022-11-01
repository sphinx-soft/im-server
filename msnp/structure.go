package msnp

import "net"

type msnp_context struct {
	dispatched bool
	ctxkey     int
	authmethod string
}

var msn_context_list []*msnp_context

/*
	type Client struct {
		Connection  net.Conn
		Client      string
		BuildNumber string
		Protocol    string
		Account     Account
	}

	type Account struct {
		UserId           int
		Email            string
		Username         string
		Password         string
		Screenname       string
		ICQNumber        int
		RegistrationTime int
	}
*/

type msnp_switchboard_context struct {
	sessionid      int
	username       string
	email          string
	authentication string
	connection     net.Conn
}

var msn_switchboard_list []*msnp_switchboard_context
