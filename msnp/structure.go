package msnp

import (
	"net"
)

type msnp_context struct {
	dispatched bool
	ctxkey     int
	authmethod string
	status     string
}

var msn_context_list []*msnp_context

type msnp_switchboard_context struct {
	sessionid      int
	username       string
	email          string
	authentication string
	connection     net.Conn
	nsinterface    net.Conn
	nscontext      *msnp_context
}

var msn_switchboard_list []*msnp_switchboard_context

type msnp_switchboard_session struct {
	sessionid int
	clients   []*msnp_switchboard_context
}

var msn_switchboard_sessions []*msnp_switchboard_session
