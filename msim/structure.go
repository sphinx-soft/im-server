package msim

import "net"

type msim_data_pair struct {
	Key   string
	Value string
}

type msim_contact struct {
	fromid int
	id     int
	reason string
}

type Msim_Client struct {
	Connection  net.Conn
	Nonce       string
	Sessionkey  int
	Account     msim_account
	StatusCode  string
	StatusText  string
	BuildNumber string
}

type msim_account struct {
	Uid        int
	Username   string
	Password   string
	Screenname string
	Avatar     string
	avatartype string
	BandName   string
	SongName   string
	Age        string
	Gender     string
	Location   string
	headline   string
	lastlogin  int64
}

type msim_offlinemessage struct {
	fromid int
	toid   int
	date   int
	msg    string
}

var Msim_Clients []*Msim_Client
