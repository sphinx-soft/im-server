package msnp

import (
	"phantom/global"
	"phantom/util"
)

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

		client := global.Client{
			Connection: tcpClient,
		}

		global.AddClient(&client)

		ctx := msnp_context{
			dispatched: true,
		}

		for {
			data, success := util.ReadTraffic(client.Connection)

			handleClientIncomingPackets(&client, &ctx, string(data))

			if !success {
				break
			}
		}

		for i := 0; i < len(global.Clients); i++ {
			if global.Clients[i].Account.Email == client.Account.Email {
				util.Debug("MSNP -> HandleNotification", "Removing from clients from List...")
				global.Clients = global.RemoveClient(global.Clients, i)
			}
		}
		client.Connection.Close()
	}
}

func HandleDispatch(client *global.Client, firstread string) {
	util.Log("MSN Messenger", "Client awaiting dispatch from %s", client.Connection.RemoteAddr().String())

	ctx := msnp_context{
		dispatched: true,
	}

	// Send first response command to MSN Client, Requesting INF Data
	if !handleProtocolVersionRequest(client, firstread) {
		util.Debug("MSNP -> HandleDispatch", "Unsupported MSNP Version requested, closing...")
		return
	}

	for {
		data, success := util.ReadTraffic(client.Connection)

		handleClientIncomingPackets(client, &ctx, string(data))

		if !success {
			break
		}
	}

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.Email == client.Account.Email {
			util.Debug("MSNP -> HandleDispatch", "Removing from clients from List...")
			global.Clients = global.RemoveClient(global.Clients, i)
		}
	}
	client.Connection.Close()
}
