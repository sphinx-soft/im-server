package msim

import "net"

type msim_data_pair struct {
	Key   string
	Value string
}
type Msim_Contact struct {
	fromid int
	id     int
	reason string
}
type Msim_Client struct {
	Connection  net.Conn
	Nonce       string
	Sessionkey  int
	Account     Msim_Account
	StatusCode  string
	StatusText  string
	BuildNumber string
}

type Msim_Account struct {
	Uid        int
	Username   string
	Password   string
	Screenname string
	Avatar     string
	BandName   string
	SongName   string
	Age        string
	Gender     string
	Location   string
}

type Msim_OfflineMessage struct {
	fromid int
	toid   int
	date   int64
	msg    string
}

var Msim_Clients []*Msim_Client
