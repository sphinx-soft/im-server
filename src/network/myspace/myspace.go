package myspace

import (
	"chimera/network"
	"chimera/utility"
	"chimera/utility/logging"
	"chimera/utility/tcp"
	"fmt"
	"strings"

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

			MySpaceHandleClientBroadcastSigninStatus(&client, &context)
			MySpaceHandleClientOfflineMessagesDelivery(&client, &context)

			for {
				stream, err := client.Connection.ReadTraffic()

				recv := strings.Split(stream, "final\\")
				for ix := 0; ix < len(recv); ix++ {
					if strings.Contains(recv[ix], "\\") {
						fix := fmt.Sprintf("%sfinal\\", recv[ix]) // this is a side effect of the split, we need to reattach the final as to not break everything
						logging.Trace("MySpace/Service", "Split TCP Readout: %s", fix)
						MySpaceHandleClientIncomingPackages(&client, &context, fix)
					}
				}

				if err != nil || MySpaceHandleClientLogoutRequest(stream) {
					break
				}
			}

			MySpaceHandleClientBroadcastLogoffStatus(&client, &context)

			logging.Info("MySpace", "Client signed out! (UIN: %d, SN: %s)", client.ClientAccount.UIN, client.ClientAccount.DisplayName)

			for i := 0; i < len(network.Clients); i++ {
				if network.Clients[i].ClientAccount.UIN == client.ClientAccount.UIN {
					logging.Debug("MySpace/Service", "Removing from client from Client List...")
					network.Clients = slices.Delete(network.Clients, i, i+1)
				}
			}

			for ix := 0; ix < len(clientContexts); ix++ {
				if clientContexts[ix].SessionKey == context.SessionKey {
					logging.Debug("MySpace/Service", "Removing from client from Context List...")
					clientContexts = slices.Delete(clientContexts, ix, ix+1)
				}
			}
		}()
	}

	//	for {
	//		bridge.ProcessMessages(bridge.ServiceMySpace)
	//	}

}
