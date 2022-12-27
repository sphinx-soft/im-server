package myspace

import (
	"chimera/bridge"
	"chimera/network"
	"chimera/utility/logging"
	"chimera/utility/tcp"
)

func LogonMySpace() {

	tcpServer := tcp.CreateListener(1863)

	for {
		err := tcpServer.AcceptClient()

		go func() {
			if err != nil {
				logging.Error("MySpace/Service", "Failed to accept client! (%s)", err.Error())
				return
			}

			logging.Info("MySpace", "Client awaiting authentication (IP: %s)", tcpServer.GetRemoteAddress())

			client := network.Client{
				Connection: tcpServer,
			}

			network.Clients = append(network.Clients, client)

		}()
	}

	for {
		bridge.ProcessMessages(bridge.ServiceMySpace)
	}

}
