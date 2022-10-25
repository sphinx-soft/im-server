package msim

import (
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingPersistPackets(client *Msim_Client, data []byte) {
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
		if strings.Contains(str, "\\cmd\\514") || strings.Contains(str, "\\cmd\\2") {
			if strings.Contains(str, "\\dsn\\8") && strings.Contains(str, "\\lid\\13") {
				handleClientPacketChangePicture(client, data)
			}
		}
	}
}

func handleClientIncomingPackets(client *Msim_Client, data []byte) {
	str := string(data)

	if strings.Contains(str, "\\status") {
		handleClientPacketSetStatusMessages(client, data)
	}
	if strings.Contains(str, "\\addbuddy") {
		handleClientPacketAddBuddy(client, data)
	}
	if strings.Contains(str, "\\delbuddy") {
		handleClientPacketDelBuddy(client, data)
	}
	if strings.Contains(str, "\\bm\\1") {
		handleClientPacketBuddyInstantMessage(client, data)
	}
}

func HandleClientKeepalive(client *Msim_Client) {
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

func HandleClients(client *Msim_Client) {
	util.Log("MySpaceIM", "Client awaiting authentication from %s", client.Connection.RemoteAddr().String())

	if !handleClientAuthentication(client) {
		client.Connection.Close()
		return
	}

	Msim_Clients = append(Msim_Clients, client)

	//unused for now
	/*global := util.Global_Client{
		Client:   "MySpace",
		Build:    client.BuildNumber,
		Protocol: "MSIMv?",
		Username: client.Account.Screenname,
		Friends:  69,
	}
	util.AddGlobalClient(&global)*/

	handleClientOfflineEvents(client)
	for {
		data, success := util.ReadTraffic(client.Connection)

		//split packets that get sent simultaneously
		receivedpackets := strings.Split(string(data), "final\\")
		for i := 0; i < len(receivedpackets); i++ {
			if strings.Contains(receivedpackets[i], "\\") {
				util.Debug("MySpace -> HandleClients -> TCP", "Reading Split Data: %s", string(receivedpackets[i]+"final\\"))
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
	for i := 0; i < len(Msim_Clients); i++ {
		if Msim_Clients[i].Account.Uid != client.Account.Uid {
			res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE fromid= ?", client.Account.Uid)
			for res.Next() {
				var msg msim_contact
				_ = res.Scan(&msg.fromid, &msg.id, &msg.reason)
				if Msim_Clients[i].Account.Uid == msg.id {
					res2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE fromid= ? AND id= ?", Msim_Clients[i].Account.Uid, client.Account.Uid)
					res2.Next()
					var count int
					res2.Scan(&count)
					res2.Close()
					if count > 0 {
						util.WriteTraffic(Msim_Clients[i].Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", client.Account.Uid),
							msim_new_data_string("msg", "|s|0|ss|"+client.StatusText),
						}))
					}

				}
			}
			res.Close()
		}
	}

	for i := 0; i < len(Msim_Clients); i++ {
		if Msim_Clients[i].Account.Username == client.Account.Username {
			util.Debug("MySpace -> HandleClients", "Removing from clients from List...")
			Msim_Clients = RemoveMsimClient(Msim_Clients, i)
		}
	}
	client.Connection.Close()
}
