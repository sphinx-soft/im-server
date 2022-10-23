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
			if strings.Contains(str, "\\dsn\\0") {
				if strings.Contains(str, "\\lid\\1") {
					handleClientPacketGetContactList(client, data)
				}

				if strings.Contains(str, "\\lid\\2") {
					handleClientPacketGetContactInformation(client, data)
				}

			}
			if strings.Contains(str, "\\dsn\\1") {
				if strings.Contains(str, "lid\\4") {
					handleClientPacketUserLookupIMAboutMyself(client, data)
				}

				if strings.Contains(str, "lid\\7") || strings.Contains(str, "lid\\17") {
					handleClientPacketUserLookupIMByUid(client, data)
				}
			}

			if strings.Contains(str, "\\dsn\\2") && strings.Contains(str, "\\lid\\6") {
				handleClientPacketGetGroups(client, data)
			}

			if strings.Contains(str, "\\dsn\\4") {
				if strings.Contains(str, "lid\\3") || strings.Contains(str, "lid\\5") {
					handleClientPacketUserLookupMySpaceByUid(client, data)
				}
			}

			if strings.Contains(str, "\\dsn\\5") && strings.Contains(str, "\\lid\\7") {
				handleClientPacketUserLookupMySpaceByUsernameOrEmail(client, data)
			}
		}
	}
}

func handleClientIncomingPackets(client *Msim_client, data []byte) {
	str := string(data)

	if strings.Contains(str, "\\status") {
		handleClientPacketSetStatusMessages(client, data)
	}
	if strings.Contains(str, "\\addbuddy") {
		handleClientPacketAddBuddy(client, data)
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
	handleClientOfflineEvents(client)
	for {
		data, success := util.ReadTraffic(client.Connection)

		//split packets that get sent simultaneously
		receivedpackets := strings.Split(string(data), "final\\")
		for i := 0; i < len(receivedpackets); i++ {
			if strings.Contains(receivedpackets[i], "\\") {
				util.Log("TCP", "Reading Data: %s", string(receivedpackets[i]+"final\\"))
				handleClientIncomingPackets(client, []byte(receivedpackets[i]+"final\\"))
				handleClientIncomingPersistPackets(client, []byte(receivedpackets[i]+"final\\"))
			}
		}

		if !success {
			break
		}
	}

	util.Log("MySpaceIM", "Client Disconnected! | Screenname: %s", client.Account.Screenname)

	//notify all users that user logged out
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Uid != client.Account.Uid {
			util.WriteTraffic(Clients[i].Connection, buildDataPacket([]msim_data_pair{
				msim_new_data_int("bm", 100),
				msim_new_data_int("f", client.Account.Uid),
				msim_new_data_string("msg", "|s|0|ss|"+client.StatusText),
			}))
		}
	}
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Username == client.Account.Username {
			util.Debug("Removing from clients array")
			Clients = ArrayRemove(Clients, i)
		}
	}
	client.Connection.Close()
}
