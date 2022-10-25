package msnp

import "phantom/util"

func HandleNotification() {
	tcpServer := util.CreateListener(1864)

	for {
		tcpClient, err := tcpServer.Accept()

		if err != nil {
			util.Error("Failed to accept Client! ", err.Error())
		} else {
			util.Debug("MSNP -> HandleNotification", "Accepted Client")
		}

		util.Log("MSN Messenger", "Client awaiting authentication from %s", tcpClient.RemoteAddr().String())

		client := Msnp_Client{
			Connection: tcpClient,
			Dispatched: true,
		}

		Msnp_Clients = append(Msnp_Clients, &client)

		for {
			data, success := util.ReadTraffic(client.Connection)

			handleClientIncomingPackets(&client, string(data))

			if !success {
				break
			}
		}
	}
}

func HandleDispatch(client *Msnp_Client, firstread string) {
	util.Log("MSN Messenger", "Client awaiting dispatch from %s", client.Connection.RemoteAddr().String())

	// Send first response command to MSN Client, Requesting INF Data
	if !handleProtocolVersionRequest(client, firstread) {
		util.Debug("MSNP -> HandleDispatch", "Unsupported MSNP Version requested, closing...")
		return
	}

	for {
		data, success := util.ReadTraffic(client.Connection)

		handleClientIncomingPackets(client, string(data))

		if !success {
			break
		}
	}

	for i := 0; i < len(Msnp_Clients); i++ {
		if Msnp_Clients[i].Account.Email == client.Account.Email {
			util.Debug("MSNP -> HandleDispatch", "Removing from clients from List...")
			Msnp_Clients = RemoveMsnpClient(Msnp_Clients, i)
		}
	}
	client.Connection.Close()
}
