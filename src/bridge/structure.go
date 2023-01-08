package bridge

type BridgeClient struct {
	ServiceId  int
	ServiceRev string
}

type BridgeDelivery struct {
	SenderId int
	RecvId   int
	Action   string
}

var clients []*BridgeClient
var packets []*BridgeDelivery

const (
	ServiceAPI     int = 0
	ServiceMySpace int = 1
	ServiceMSN     int = 2
	ServiceYahoo   int = 3
	ServiceAIM     int = 4
)
