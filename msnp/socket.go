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

		util.ReadTraffic(tcpClient)
	}
}

func HandleDispatch(client *Msnp_Client, firstread string) {
	util.Log("MSN Messenger", "Client awaiting dispatch from %s", client.Connection.RemoteAddr().String())

	// Send first response command to MSN Client, Requesting INF Data
	if !handleProtocolVersionRequest(client, firstread) {
		util.Debug("MSNP -> HandleDispatch", "Unsupported MSNP Version requested, closing...")
		client.Connection.Close()
		return
	}

	for {
		data, success := util.ReadTraffic(client.Connection)

		handleClientIncomingDispatchPackets(client, string(data))

		if !success {
			break
		}
	}
}
