package msnp

import (
	"encoding/hex"
	"fmt"
	"phantom/global"
	"phantom/util"
	"strconv"
	"strings"
	"time"
)

/*all of this is DS and NS only not SS/SB*/
func handleClientIncomingPackets(client *global.Client, ctx *msnp_context, data string) {

	switch {
	case strings.HasPrefix(data, "VER"):
		handleClientPacketNegotiateProtocolVersion(client, data)
	case strings.HasPrefix(data, "INF"):
		handleClientPacketAuthenticationMethod(client, ctx, data)
	case strings.HasPrefix(data, "USR"):
		handleClientPacketAuthentication(client, ctx, data)
	case strings.HasPrefix(data, "SYN"):
		handleClientPacketContactListSynchronization(client, data)
	case strings.HasPrefix(data, "CHG"):
		handleClientPacketChangeStatusRequest(client, ctx, data)
	case strings.HasPrefix(data, "CVR"):
		handleClientPacketGetClientServerInformation(client, data)
	case strings.HasPrefix(data, "ADD"):
		handleClientPacketUpdateContactRequest(client, data)
	case strings.HasPrefix(data, "XFR"):
		handleClientPacketSwitchboardSessionRequest(client, ctx, data)
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

	/*
		decode := strings.Replace(packet, "\r\n", "", -1)
			splits := strings.Split(decode, " ")

			for ix := 0; ix < len(splits); ix++ {
				if splits[ix] == data_search {
					//return splits[ix+1+len(offset)]
					//string(bytes.Trim([]byte(splits[1]), "\x00"))
					return string(bytes.Trim([]byte(splits[ix+1+len(offset)]), "\x00"))
				}
			}

			return ""
	*/

	versions := make([]string, 4096)
	strdata := strings.Replace(data, "\r\n", "", -1)
	splits := strings.Split(strdata, " ")

	for ix := 2; ix < len(splits); ix++ {
		if splits[ix] == "CVR0" {
			break
		}

		versions = append(versions, splits[ix])
	}

	protoverstr := versions[len(versions)-1]
	protoverstripped := strings.Replace(protoverstr, "MSNP", "", -1)
	protover, _ := strconv.Atoi(protoverstripped)

	if protover <= 7 {
		client.Protocol = protoverstr
		util.Debug("MSNP -> handleProtocolVersionRequest", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "VER"))))
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", protoverstr))
		return true
	} else {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "VER", "CVR0"))
		return false
	}
}

func handleClientPacketNegotiateProtocolVersion(client *global.Client, data string) {
	handleClientProtocolVersionRequest(client, data)
}

func handleClientPacketAuthenticationMethod(client *global.Client, ctx *msnp_context, data string) {

	protoverstr := strings.Replace(client.Protocol, "MSNP", "", -1)
	protover, _ := strconv.Atoi(protoverstr)
	var authmethod string

	if protover <= 2 {
		authmethod = "CTP"
	} else if protover > 2 && protover <= 7 {
		authmethod = "MD5"
	}

	ctx.authmethod = authmethod

	util.Debug("MSNP -> handleClientPacketAuthenticationMethod", fmt.Sprintf("TrID Dbg: %v", []byte(getTrId(data, "INF"))))
	util.WriteTraffic(client.Connection, msnp_new_command(data, "INF", authmethod))
}

func handleClientPacketAuthentication(client *global.Client, ctx *msnp_context, data string) {
	if !ctx.dispatched {
		util.WriteTraffic(client.Connection, msnp_new_command(data, "XFR", fmt.Sprintf("NS %s:1864", util.GetRootUrl())))
		util.Log("MSN Messenger", "Redirecting Client to Notification Server...")
	} else {

		account := strings.Replace(findValueFromData("I", data, 0), "@hotmail.com", util.GetMailDomain(), -1)
		ctx.email = account
		client.Account, _ = global.GetUserDataFromEmail(account)
		password := strings.Replace(util.DecryptAES(util.GetAESKey(), client.Account.Password), "\r\n", "", -1)
		var clpw string

		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "aes test: %s", util.GetAESKey())
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "um data test: %s", account)
		util.Debug("MSNP -> handleClientPacketAuthenticationBegin", "pw data AES: %v", []byte(password))

		if ctx.authmethod == "CTP" {
			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> CTP", "pw data CTP test1: %v", []byte(strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1)))
			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> CTP", "pw data CTP test2: %v", []byte(password))

			clpw = strings.Replace(findValueFromData("I", data, 1), "\r\n", "", -1)

		} else if ctx.authmethod == "MD5" {
			saltpw := fmt.Sprintf("%s%s", hex.EncodeToString([]byte(fmt.Sprintf("%d", client.Account.RegistrationTime))), password)
			//unix :=

			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> MD5", "md5 salt test: %v", []byte(hex.EncodeToString([]byte(fmt.Sprintf("%d", client.Account.RegistrationTime)))))
			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> MD5", "pw data MD5 test: %v", []byte(password))
			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> MD5", "pw data MD5 test2: %v", []byte(util.HashMD5(saltpw)))
			util.Debug("MSNP -> handleClientPacketAuthenticationBegin -> MD5", "pw data MD5 test2 plain: %v", util.HashMD5(saltpw))

			util.WriteTraffic(client.Connection, msnp_new_command(data, "USR", fmt.Sprintf("MD5 S %s", hex.EncodeToString([]byte(fmt.Sprintf("%d", client.Account.RegistrationTime))))))

			datanew, _ := util.ReadTraffic(client.Connection)
			clpw = findValueFromData("MD5", string(datanew), 1)
			password = util.HashMD5(saltpw)
		}

		if clpw == password {
			// manually increase trid if not md5
			trid, _ := strconv.Atoi(getTrId(data, "USR"))

			if ctx.authmethod != "CTP" {
				trid++
			}

			// we cant use msnp_new_command here because the data never changes
			resp := fmt.Sprintf("USR %d OK %s %s\r\n", trid, client.Account.Email, client.Account.Screenname)

			util.WriteTraffic(client.Connection, resp)
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

func handleClientPacketChangeStatusRequest(client *global.Client, ctx *msnp_context, data string) {

	//todo
	util.WriteTraffic(client.Connection, msnp_new_command(data, "CHG", findValueFromData("CHG", data, 1)))

	ctx.status = findValueFromData("CHG", data, 1)
}

// [MySpaceIM] Client Authenticated! | Username: test@phantom-im.xyz | Screenname: TestUser | Version: 1.0.595.0
func handleClientPacketGetClientServerInformation(client *global.Client, data string) {

	build := findValueFromData("CVR", data, 6)
	client.BuildNumber = build

	util.WriteTraffic(client.Connection, msnp_new_command(data, "CVR", fmt.Sprintf("%s %s %s %s %s", build, build, "1.0.0000", "https://archive.org/download/MsnMessengerClients2/MSN%20Messenger%201.0.0863%20%28English%20-%20United%20States%29.zip", "http://phantom-im.xyz")))

	util.Log("MSN Messenger", "Client Authenticated! -> Email: %s, Screenname: %s, Version: %s, Protocol Version: %s", client.Account.Email, client.Account.Screenname, client.BuildNumber, client.Protocol)
}

/*todo*/
func handleClientPacketUpdateContactRequest(client *global.Client, data string) {
	list := findValueFromData("ADD", data, 1)
	mail := findValueFromData("ADD", data, 3)

	if !strings.Contains(data, "@hotmail.com") || !strings.Contains(data, util.GetMailDomain()) {
		util.WriteTraffic(client.Connection, msnp_new_command_noargs(data, "201"))
	}

	var count int
	check2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from accounts WHERE email=?", mail)
	check2.Next()
	check2.Scan(&count)
	check2.Close()

	if count < 1 {
		util.WriteTraffic(client.Connection, msnp_new_command_noargs(data, "205"))
	}

	var username string
	if strings.Contains(mail, "@hotmail.com") {
		username = strings.Replace(mail, "@hotmail.com", "", -1)
	} else if strings.Contains(mail, util.GetMailDomain()) {
		username = strings.Replace(mail, util.GetMailDomain(), "", -1)
	}

	toAcc, _ := global.GetUserDataFromUsername(username)

	check2, _ = util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE to_id=? and from_id= ?", toAcc.UserId, client.Account.UserId)
	check2.Next()
	check2.Scan(&count)
	check2.Close()

	if count > 0 {
		util.WriteTraffic(client.Connection, msnp_new_command_noargs(data, "215"))
	}

	if list == "AL" {

	} else if list == "BL" {

	} else if list == "FL" {

		util.WriteTraffic(client.Connection, msnp_new_command(data, "ADD", fmt.Sprintf("%s 1 %s %s", list, mail, mail)))
	}
}

func handleClientPacketSwitchboardSessionRequest(client *global.Client, ctx *msnp_context, data string) {
	date := time.Now().UTC().UnixMilli()
	sbctx := msnp_switchboard_context{
		authentication: strconv.FormatInt(date, 16),
		email:          client.Account.Email,
		nscontext:      ctx,
		nsinterface:    client.Connection,
	}
	addSwitchboardContext(&sbctx)

	util.WriteTraffic(client.Connection, msnp_new_command(data, "XFR", fmt.Sprintf("SB %s:1865 CKI %s", util.GetRootUrl(), sbctx.authentication)))
}
