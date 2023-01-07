package myspace

import (
	"chimera/network"
	"chimera/utility"
	"chimera/utility/logging"
	"chimera/utility/tcp"

	"golang.org/x/exp/slices"
)

var lastSessKey int = 0

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

			context := MySpaceContext{
				Nonce:      utility.RandomString(0x40),
				SessionKey: lastSessKey + 1,
			}

			// here's where the client actually starts doing shit

			if !MySpaceHandleClientAuthentication(&client, &context) {
				logging.Warn("MySpace", "Client Failed Authentication! Signing Off...")
				client.Connection.CloseConnection()
				return
			}

			network.Clients = append(network.Clients, &client)
			clientContexts = append(clientContexts, &context)

			LogoutMySpace(&client, &context)
		}()
	}

	//	for {
	//		bridge.ProcessMessages(bridge.ServiceMySpace)
	//	}

}

func LogoutMySpace(cli *network.Client, ctx *MySpaceContext) {

	logging.Info("MySpace", "Client signed out! (UIN: %d, SN: %s)", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName)

	for i := 0; i < len(network.Clients); i++ {
		if network.Clients[i].ClientAccount.UIN == cli.ClientAccount.UIN {
			logging.Debug("MySpace/Service", "Removing from client from Client List...")
			network.Clients = slices.Delete(network.Clients, i, i+1)
		}
	}

	for ix := 0; ix < len(clientContexts); ix++ {
		if clientContexts[ix].SessionKey == ctx.SessionKey {
			logging.Debug("MySpace/Service", "Removing from client from Context List...")
			clientContexts = slices.Delete(clientContexts, ix, ix+1)
		}
	}
}
