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
				if strings.Contains(str, "lid\\7") || strings.Contains(str, "lid\\4") || strings.Contains(str, "lid\\17") {
					handleClientPacketUserLookupIMByUidOrMyself(client, data)
				}
			}
			if strings.Contains(str, "\\dsn\\4") {
				if strings.Contains(str, "lid\\5") || strings.Contains(str, "lid\\3") {
					handleClientPacketUserLookupMySpaceByUidOrMyself(client, data)
				}
			}
			if strings.Contains(str, "\\dsn\\0") && strings.Contains(str, "\\lid\\2") {
				handleClientPacketGetContactInformation(client, data)
			}
			if strings.Contains(str, "\\dsn\\2") && strings.Contains(str, "\\lid\\6") {
				handleClientPacketGetGroups(client, data)
			}
			if strings.Contains(str, "\\dsn\\0") && strings.Contains(str, "\\lid\\1") {
				handleClientPacketGetContactList(client, data)
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
	if strings.Contains(str, "\\addbuddy") {
		handleClientAddBuddy(client, data)
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
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Username == client.Account.Username {
			util.Debug("Removing from clients array")
			Clients = ArrayRemove(Clients, i)
		}
	}
	client.Connection.Close()
}
