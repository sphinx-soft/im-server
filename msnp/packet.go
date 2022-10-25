package msnp

import (
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingDispatchPackets(client *Msnp_Client, data string) {

	switch {
	case strings.Contains(data, "INF"):
		handleClientPacketAuthenticationMethod(client, data)
	}

}

func handleProtocolVersionRequest(client *Msnp_Client, data string) bool {

	if strings.Contains(data, "MSNP2") {
		util.WriteTraffic(client.Connection, buildCommand(msnp_new_command(data, "VER", "MSNP2 CVR0")))
		return true
	} else {
		util.WriteTraffic(client.Connection, buildCommand(msnp_new_command(data, "VER", "CVR0")))
		return false
	}
}

func handleClientPacketAuthenticationMethod(client *Msnp_Client, data string) {
	//todo
	time.Sleep(time.Millisecond * 150)
	util.WriteTraffic(client.Connection, buildCommand(
		msnp_new_command(data, "INF", "MD5"),
	))
}
