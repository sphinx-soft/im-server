package msim

import (
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingPackets(client Msim_client, data []byte) {
	str := string(data)

	switch {
	case strings.Contains(str, "\\persist"):
		handleClientPacketUserLookup(client, data)
	case strings.Contains(str, "\\status"):
		handleClientSetStatusMessages(client, data)
	}

}

func HandleClientKeepalive(client Msim_client) {
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

func HandleClients(client Msim_client) {
	util.Log("MySpaceIM", "Client awaiting authentication from %s", client.Connection.RemoteAddr().String())

	if !handleClientAuthentication(client) {
		client.Connection.Close()
		return
	}

	for {
		data, success := util.ReadTraffic(client.Connection)
		handleClientIncomingPackets(client, data)
		if !success {
			break
		}
	}

	util.Log("MySpaceIM", "Client Disconnected! | Screenname: %s", client.Account.Screenname)
	client.Connection.Close()
}
