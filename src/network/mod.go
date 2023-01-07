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
	Mail        string
	Password    string
}

type User struct {
	UIN             int
	AvatarBlob      string
	AvatarImageType string
	StatusMessage   string
	LastLogin       int64
	SignupDate      int
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

var Clients []*Client
