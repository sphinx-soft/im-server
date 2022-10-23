package msim

import "net"

func ArrayRemove(s []*Msim_client, i int) []*Msim_client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type msim_data_pair struct {
	Key   string
	Value string
}
type Contact struct {
	fromid int
	id     int
	reason string
}
type Msim_client struct {
	Connection net.Conn
	Nonce      string
	Sessionkey int
	Account    Account
	StatusCode string
	StatusText string
}

type Account struct {
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

var Clients []*Msim_client
