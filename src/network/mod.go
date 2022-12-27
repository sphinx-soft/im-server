package network

import (
	"chimera/utility/tcp"
)

type Client struct {
	Connection    tcp.TcpConnection
	ClientInfo    Details
	ClientAccount Account
	ClientUser    User
}

type Details struct {
	Messenger string
	Build     string
	Protocol  string
}

type Account struct {
	UIN         int
	DisplayName string
	Username    string
	Password    string
	SignupDate  int
}

type User struct {
	AvatarBlob string
}

type Contact struct {
	SenderUIN int
	FriendUIN int
	Reason    string
}

type OfflineMessage struct {
	SenderUIN      int
	RecvUIN        int
	MessageDate    int
	MessageContent string
}

var Clients []Client
