package msim

import (
	"phantom/global"
	"phantom/util"
	"strings"
)

func HandleClients(client *global.Client) {
	util.Log(util.INFO, "MySpaceIM", "Client awaiting authentication from %s", client.Connection.RemoteAddr().String())

	client.Client = "MySpaceIM"

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

	handleClientBroadcastSignOnStatus(client, &ctx)
	handleClientHandleOfflineMessages(client, &ctx)

	for {
		data, success := util.ReadTraffic(client.Connection)

		//split packets that get sent simultaneously
		receivedpackets := strings.Split(string(data), "final\\")
		for i := 0; i < len(receivedpackets); i++ {
			if strings.Contains(receivedpackets[i], "\\") {
				util.Log(util.TRACE, "MySpace -> HandleClients -> TCP", "Reading Split Data: %s", string(receivedpackets[i]+"final\\"))
				handleClientIncomingPackets(client, &ctx, []byte(receivedpackets[i]+"final\\"))
				handleClientIncomingPersistPackets(client, &ctx, []byte(receivedpackets[i]+"final\\"))
			}
		}

		if !success || handleClientLogoutRequest(string(data)) {
			break
		}
	}

	handleClientBroadcastSignOffStatus(client, &ctx)

	util.Log(util.INFO, "MySpaceIM", "Client Disconnected -> Username: %s", client.Account.Username)

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.Email == client.Account.Email {
			util.Log(util.TRACE, "MySpace -> HandleClients", "Removing from clients from Client List...")
			global.Clients = global.RemoveClient(global.Clients, i)
		}
	}

	for ix := 0; ix < len(users_context); ix++ {
		if users_context[ix].sesskey == ctx.sesskey {
			util.Log(util.TRACE, "MySpace -> HandleClients", "Removing from clients from Context List...")
			users_context = removeUserContext(users_context, ix)
		}
	}

	client.Connection.Close()
}
