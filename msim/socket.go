package msim

import (
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingPersistPackets(client *Msim_client, data []byte) {
	str := string(data)

	if strings.Contains(str, "\\persist\\1") {
		if strings.Contains(str, "\\cmd\\1") {
			if strings.Contains(str, "\\dsn\\1") {
				if strings.Contains(str, "lid\\7") || strings.Contains(str, "lid\\17") {
					handleClientPacketUserLookupByUid(client, data)
				}
			}

			if strings.Contains(str, "\\dsn\\5") && strings.Contains(str, "\\lid\\7") {
				handleClientPacketUserLookupByUsernameOrEmail(client, data)
			}
		}
	}
}

func handleClientIncomingPackets(client *Msim_client, data []byte) {
	str := string(data)

	if strings.Contains(str, "\\status") {
		handleClientSetStatusMessages(client, data)
	}
}

func HandleClientKeepalive(client *Msim_client) {
	for {
		time.Sleep(180 * time.Second)
		err := util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("ka", true),
		}))
		if err != nil {
			break
		}
	}
}

func HandleClients(client *Msim_client) {
	util.Log("MySpaceIM", "Client awaiting authentication from %s", client.Connection.RemoteAddr().String())

	if !handleClientAuthentication(client) {
		client.Connection.Close()
		return
	}
	Clients = append(Clients, client)
	for {
		data, success := util.ReadTraffic(client.Connection)
		handleClientIncomingPackets(client, data)
		handleClientIncomingPersistPackets(client, data)
		if !success {
			break
		}
	}

	util.Log("MySpaceIM", "Client Disconnected! | Screenname: %s", client.Account.Screenname)
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Username == client.Account.Username {
			util.Debug("Removing from clients array")
			Clients = ArrayRemove(Clients, i)
		}
	}
	client.Connection.Close()
}
