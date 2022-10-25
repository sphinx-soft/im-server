package msnp

import (
	"fmt"
	"phantom/util"
	"strings"
)

func handleClientIncomingPackets(client *Msnp_Client, data string) {

	switch {
	case strings.Contains(data, "VER"):
		handleClientPacketNegotiateProtocolVersion(client, data)
	case strings.Contains(data, "INF"):
		handleClientPacketAuthenticationMethod(client, data)
	case strings.Contains(data, "USR") && strings.Contains(data, "I"):
		handleClientPacketAuthenticationBegin(client, data)
	}

}

func handleProtocolVersionRequest(client *Msnp_Client, data string) bool {
	if strings.Contains(data, "MSNP2") {
		util.Debug("MSNP -> handleProtocolVersionRequest", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "VER"))))
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", "MSNP2"))
		return true
	} else {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", "CVR0"))
		return false
	}
}

func handleClientPacketNegotiateProtocolVersion(client *Msnp_Client, data string) {
	handleProtocolVersionRequest(client, data)
}

func handleClientPacketAuthenticationMethod(client *Msnp_Client, data string) {
	//todo
	//str := strings.Replace(data, "\r\n", "", -1)
	util.Debug("MSNP -> handleClientPacketAuthenticationMethod", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "INF"))))
	util.WriteTraffic(client.Connection, msnp_new_command(data, "INF", "CTP"))
}

func handleClientPacketAuthenticationBegin(client *Msnp_Client, data string) {
	if !client.Dispatched {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "XFR", "NS localhost:1864"))
		util.Log("MSN Messenger", "Redirecting Client to Notification Server...")
		client.Dispatched = true
	} else {
		account := strings.Replace(findValueFromData("I", data), "@hotmail.com", "@phantom-im.xyz", -1)
		client.Account = getUserData(account)

		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "um data test: %s", account)
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data test1: %v", []byte(strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1)))
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data test2: %v", []byte(client.Account.Password))

		if client.Account.Password == strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1) {
			util.WriteTraffic(client.Connection, msnp_new_command(data, "USR", fmt.Sprintf("OK %s %s", client.Account.Email, client.Account.Screenname)))
		} else {
			//https://wiki.nina.chat/wiki/Protocols/MSNP/Reference/Error_List#911
			util.WriteTraffic(client.Connection, msnp_new_command_noargs(data, "911"))
		}
	}
}
