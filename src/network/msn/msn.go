package msn

import (
	"chimera/bridge"
	"chimera/utility/logging"
)

func LogonMSN() {
	logging.Debug("MSN", "Test")
	bridge.SendMessage(bridge.ServiceMSN, bridge.ServiceMySpace, "Test")

	for {
		bridge.ProcessMessages(bridge.ServiceMSN)
	}
}
