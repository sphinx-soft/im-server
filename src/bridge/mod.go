package bridge

import (
	"chimera/utility/logging"
	"sync"

	"golang.org/x/exp/slices"
)

var apptex sync.RWMutex
var deltex sync.RWMutex

func SignOnService(svcname string, svcid int, svcrev string, svccfg bool, svc func()) {

	if !svccfg {
		return
	}

	service := BridgeClient{
		ServiceId:  svcid,
		ServiceRev: svcrev,
	}

	clients = append(clients, service)

	logging.Info("Logon Service", "%s Service SignOn (Id: %d, Rev: %s)", svcname, svcid, svcrev)

	go svc()
}

func SendMessage(sender int, recv int, action string) {
	msg := BridgeDelivery{
		SenderId: sender,
		RecvId:   recv,
		Action:   action,
	}

	apptex.Lock()
	packets = append(packets, msg)
	apptex.Unlock()
}

func ProcessMessages(svcid int) {
	deltex.Lock()
	for ix := 0; ix < len(packets); ix++ {
		if packets[ix].RecvId == svcid {
			logging.Debug("Bridge/ProcessMessages", "SvcId: %d, Action: %s", svcid, packets[ix].Action)
			DeliverMessages(packets[ix])
			packets = slices.Delete(packets, ix, ix+1)
		}
	}
	deltex.Unlock()
}

func DeliverMessages(packet BridgeDelivery) {

}
