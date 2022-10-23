package util

type Global_Client struct {
	Client   string
	Build    string
	Protocol string
	Username string
	Friends  int
}

var Global_Clients []*Global_Client
