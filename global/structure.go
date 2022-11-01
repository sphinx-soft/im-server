package global

import "net"

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

type Contact struct {
	FromId int
	ToId   int
}

type OfflineMsg struct {
	FromId  int
	ToId    int
	Date    int
	Message string
}

type Upload struct {
	UserId int
	Avatar string
}

var Clients []*Client
