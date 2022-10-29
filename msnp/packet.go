package msnp

import (
	"fmt"
	"phantom/global"
	"phantom/util"
	"strconv"
	"strings"
)

func handleClientIncomingPackets(client *global.Client, ctx *msnp_context, data string) {

	switch {
	case strings.HasPrefix(data, "VER"):
		handleClientPacketNegotiateProtocolVersion(client, data)
	case strings.HasPrefix(data, "INF"):
		handleClientPacketAuthenticationMethod(client, data)
	case strings.HasPrefix(data, "USR") && strings.Contains(data, " I "):
		handleClientPacketAuthenticationBegin(client, ctx, data)

	case strings.HasPrefix(data, "SYN"):
		handleClientPacketContactListSynchronization(client, data)
	case strings.HasPrefix(data, "CHG"):
		handleClientPacketChangeStatusRequest(client, data)
	case strings.HasPrefix(data, "CVR"):
		handleClientPacketGetClientServerInformation(client, data)
	}

}

func handleClientLogoutRequest(data string) bool {
	if strings.HasPrefix(data, "OUT") {
		return true
	} else {
		return false
	}
}

func handleClientProtocolVersionRequest(client *global.Client, data string) bool {
	if strings.Contains(data, "MSNP2") {
		client.Protocol = "MSNP2"
		util.Debug("MSNP -> handleProtocolVersionRequest", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "VER"))))
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", "MSNP2"))
		return true
	} else {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", "CVR0"))
		return false
	}
}

func handleClientPacketNegotiateProtocolVersion(client *global.Client, data string) {
	handleClientProtocolVersionRequest(client, data)
}

func handleClientPacketAuthenticationMethod(client *global.Client, data string) {

	protoverstr := strings.Replace(client.Protocol, "MSNP", "", -1)
	protover, _ := strconv.Atoi(protoverstr)
	var authmethod string

	if protover <= 2 {
		authmethod = "CTP"
	} else if protover > 2 && protover <= 7 {
		authmethod = "MD5"
	}

	util.Debug("MSNP -> handleClientPacketAuthenticationMethod", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "INF"))))
	util.WriteTraffic(client.Connection, msnp_new_command(data, "INF", authmethod))
}

func handleClientPacketAuthenticationBegin(client *global.Client, ctx *msnp_context, data string) {
	if !ctx.dispatched {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "XFR", "NS localhost:1864"))
		util.Log("MSN Messenger", "Redirecting Client to Notification Server...")
	} else {
		account := strings.Replace(findValueFromData("I", data), "@hotmail.com", "@phantom-im.xyz", -1)
		client.Account, _ = global.GetUserDataFromEmail(account)

		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "aes test: %s", util.GetAESKey())
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "um data test: %s", account)
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data test1: %v", []byte(strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1)))
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data test2: %v", []byte(client.Account.Password))
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data AES: %v", []byte(util.DecryptAES(util.GetAESKey(), client.Account.Password)))

		aes := util.DecryptAES(util.GetAESKey(), client.Account.Password)

		if aes == strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1) {
			util.WriteTraffic(client.Connection, msnp_new_command(data, "USR", fmt.Sprintf("OK %s %s", client.Account.Email, client.Account.Screenname)))
		} else {
			//https://wiki.nina.chat/wiki/Protocols/MSNP/Reference/Error_List#911
			util.WriteTraffic(client.Connection, msnp_new_command_noargs(data, "911"))
		}
	}
}

func handleClientPacketContactListSynchronization(client *global.Client, data string) {

	var clv int

	res, _ := util.GetDatabaseHandle().Query("SELECT clversion from msn WHERE id=?", client.Account.UserId)
	res.Scan(&clv)

	util.WriteTraffic(client.Connection, msnp_new_command(data, "SYN", strconv.Itoa(clv)))

	//todo
	if findValueFromData("SYN", data, 1) == strconv.Itoa(clv) {

		res, _ = util.GetDatabaseHandle().Query("SELECT * from contacts WHERE from_id=?", client.Account.UserId)

		for res.Next() {
			var contact global.Contact
			res.Scan(&contact.FromId, &contact.ToId)

		}

	}

}

func handleClientPacketChangeStatusRequest(client *global.Client, data string) {

	//todo
	util.WriteTraffic(client.Connection, msnp_new_command(data, "CHG", findValueFromData("CHG", data, 1)))

}

// [MySpaceIM] Client Authenticated! | Username: test@phantom-im.xyz | Screenname: TestUser | Version: 1.0.595.0
func handleClientPacketGetClientServerInformation(client *global.Client, data string) {

	//todo
	build := findValueFromData("i386", data, 1)
	client.BuildNumber = build

	util.WriteTraffic(client.Connection, msnp_new_command(data, "CVR", fmt.Sprintf("%s %s %s %s %s", build, build, "1.0.0000", "https://archive.org/download/MsnMessengerClients2/MSN%20Messenger%201.0.0863%20%28English%20-%20United%20States%29.zip", "http://phantom-im.xyz")))

	util.Log("MSN Messenger", "Client Authenticated! -> Email: %s, Screenname: %s, Version: %s, Protocol Version: %s", client.Account.Email, client.Account.Screenname, client.BuildNumber, client.Protocol)
}
