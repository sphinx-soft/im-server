package msim

import (
	"phantom/global"
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingPersistPackets(client *global.Client, ctx *msim_context, data []byte) {
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

func handleClientIncomingPackets(client *global.Client, ctx *msim_context, data []byte) {
	str := string(data)

	if strings.Contains(str, "\\status") {
		handleClientPacketSetStatusMessages(client, ctx, data)
	}
	if strings.Contains(str, "\\addbuddy") {
		handleClientPacketAddBuddy(client, ctx, data)
	}
	if strings.Contains(str, "\\delbuddy") {
		handleClientPacketDelBuddy(client, data)
	}
	if strings.Contains(str, "\\bm\\1") {
		handleClientPacketBuddyInstantMessage(client, data)
	}
}

func HandleClientKeepalive(client *global.Client) {
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

func HandleClients(client *global.Client) {
	util.Log("MySpaceIM", "Client awaiting authentication from %s", client.Connection.RemoteAddr().String())

	ctx := msim_context{
		nonce:   generateNonce(),
		sesskey: generateSessionKey(),
	}

	addUserContext(&ctx)

	if !handleClientAuthentication(client, &ctx) {
		client.Connection.Close()
		return
	}

	global.AddClient(client)

	handleClientOfflineEvents(client, &ctx)
	for {
		data, success := util.ReadTraffic(client.Connection)

		//split packets that get sent simultaneously
		receivedpackets := strings.Split(string(data), "final\\")
		for i := 0; i < len(receivedpackets); i++ {
			if strings.Contains(receivedpackets[i], "\\") {
				util.Debug("MySpace -> HandleClients -> TCP", "Reading Split Data: %s", string(receivedpackets[i]+"final\\"))
				handleClientIncomingPackets(client, &ctx, []byte(receivedpackets[i]+"final\\"))
				handleClientIncomingPersistPackets(client, &ctx, []byte(receivedpackets[i]+"final\\"))
			}
		}

		if !success {
			break
		}
	}

	util.Log("MySpaceIM", "Client Disconnected! | Screenname: %s", client.Account.Screenname)

	//notify all users that user logged out
	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.UserId != client.Account.UserId {
			res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE fromid= ?", client.Account.UserId)
			for res.Next() {
				var msg global.Contact
				_ = res.Scan(&msg.FromId, &msg.ToId)
				if global.Clients[i].Account.UserId == msg.ToId {
					res2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE fromid= ? AND id= ?", global.Clients[i].Account.UserId, client.Account.UserId)
					res2.Next()
					var count int
					res2.Scan(&count)
					res2.Close()
					if count > 0 {
						util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", client.Account.UserId),
							msim_new_data_string("msg", "|s|0|ss|"+ctx.statusmessage),
						}))
					}

				}
			}
			res.Close()
		}
	}

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.Email == client.Account.Email {
			util.Debug("MySpace -> HandleClients", "Removing from clients from List...")
			global.Clients = global.RemoveClient(global.Clients, i)
		}
	}
	client.Connection.Close()
}
