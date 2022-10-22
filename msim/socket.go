package msim

import (
	"phantom/util"
	"time"
)

func handleClientIncomingPackets(client *Msim_client, data []byte) {
	handleClientPacketUserLookup(client, data)
}

func HandleClientKeepalive(client *Msim_client) {
	for {
		time.Sleep(180 * time.Second)
		err := util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("ka"),
		}))
		if err != nil {
			break
		}
	}
}

func HandleClients(client *Msim_client) {
	util.Log("MySpaceIM", "Client awaiting authentication from ", client.Connection.RemoteAddr())

	if !handleClientAuthentication(client) {
		client.Connection.Close()
		return
	}
	Clients = append(Clients, client)
	for {
		data, success := util.ReadTraffic(client.Connection)
		handleClientIncomingPackets(client, data)
		if !success {
			break
		}
	}

	util.Log("MySpaceIM", "Client Disconnected! | Screenname: "+client.Account.Screenname)
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Username == client.Account.Username {
			util.Debug("Removing from clients array")
			Clients = ArrayRemove(Clients, i)
		}
	}
	client.Connection.Close()
}
