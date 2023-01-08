package myspace

import (
	"chimera/network"
	"chimera/utility"
	"chimera/utility/database"
	"chimera/utility/encryption"
	"chimera/utility/logging"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// handlers

func MySpaceHandleClientIncomingPackages(cli *network.Client, ctx *MySpaceContext, stream string) {
	switch {
	case strings.HasPrefix(stream, "\\persist"):
		MySpaceHandleClientIncomingCallbacks(cli, ctx, stream)
	case strings.HasPrefix(stream, "\\status"):
		MySpaceHandleClientPacketSetStatusMessage(cli, ctx, stream)
	case strings.HasPrefix(stream, "\\addbuddy"):
		MySpaceHandleClientPacketAddBuddy(cli, ctx, stream)
	case strings.HasPrefix(stream, "\\delbuddy"):
		MySpaceHandleClientPacketDeleteBuddy(cli, stream)
	case strings.HasPrefix(stream, "\\bm\\1"):
		MySpaceHandleClientPacketBuddyInstantMessage(cli, ctx, stream)
	}
}

//	Reading Data: 	\persist\1\sesskey\1\cmd\1\dsn\6\uid\10000\lid\11\rid\7\body\Target=mail∟FriendID=10000\final\
//						\persist\1\sesskey\1\cmd\1\dsn\6\uid\10000\lid\11\rid\23\body\Target=mail∟FriendID=10000\final\
//
// this is hell; Also still missing the return for 1611 which is link callback request
func MySpaceHandleClientIncomingCallbacks(cli *network.Client, ctx *MySpaceContext, stream string) {
	if strings.Contains(stream, "\\persist\\1\\sesskey\\1\\cmd\\1\\dsn\\6\\uid\\10000\\lid\\11") {
		MySpaceHandleClientCallbackNetLinkRequest(cli, stream)
	}
}

func MySpaceHandleClientAuthentication(cli *network.Client, ctx *MySpaceContext) bool {
	cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
		MySpaceNewDataGeneric("lc", "1"),
		MySpaceNewDataGeneric("nc", base64.StdEncoding.EncodeToString([]byte(ctx.Nonce))),
		MySpaceNewDataGeneric("id", "1"),
	}))

	loginPacket, err := cli.Connection.ReadTraffic()
	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to read Login2 data packet! (%s)", err.Error())
		return false
	}

	email := MySpaceRetrieveKeyValue("username", loginPacket) // i have no clue why MySpace called this a "Username" when its the email bruh
	clientver := MySpaceRetrieveKeyValue("clientver", loginPacket)
	response := MySpaceRetrieveKeyValue("response", loginPacket)

	cli.ClientAccount, err = database.GetAccountDataByEmail(email)
	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch Account Data! (%s)", err.Error())
		return false
	}

	logging.Debug("MySpace/Authentication", "NetworkAccount UIN: %d", cli.ClientAccount.UIN)
	logging.Debug("MySpace/Authentication", "NetworkAccount SN: %s", cli.ClientAccount.DisplayName)
	logging.Debug("MySpace/Authentication", "NetworkAccount Mail: %s", cli.ClientAccount.Mail)
	logging.Debug("MySpace/Authentication", "NetworkAccount PW: %v", []byte(cli.ClientAccount.Password))

	/*todo: implement AES, again...*/

	logging.Debug("MySpace/Authentication", "AES Password Stuff: %v", []byte(cli.ClientAccount.Password))

	verifybuf := make([]byte, 32)
	noncebuf := make([]byte, 32)
	rc4key := make([]byte, 16)
	challenge := []byte(ctx.Nonce)
	for i := 0; i < 32; i++ {
		noncebuf[i] = challenge[i+32]
		verifybuf[i] = challenge[i]
	}

	pwdbytes := utility.ConvertToUTF16LE(cli.ClientAccount.Password)
	hasher := sha1.New()
	hasher.Write(pwdbytes)
	stage1 := hasher.Sum(nil)

	logging.Debug("MySpace/Authentication", "SHA1 Stage 1: %v", []byte(stage1))

	stage2 := append(stage1, noncebuf...)
	hasher.Reset()
	hasher.Write(stage2)
	byte_hash_total := hasher.Sum(nil)
	hasher.Reset()

	for i := 0; i < 16; i++ {
		rc4key[i] = byte_hash_total[i]
	}

	logging.Debug("MySpace/Authentication", "SHA1 Stage 2: %v", []byte(stage2))
	logging.Debug("MySpace/Authentication", "RC4 Decryption Key: %v", []byte(rc4key))

	b64data, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		logging.Error("MySpace/Authentication", "Invalid base64 provided at login packet. (%s)", err.Error())
		return false
	}

	blob := encryption.SwapRC4State(rc4key, b64data)

	logging.Debug("MySpace/Authentication", "RC4 Data Blob (str): %s", blob)
	logging.Debug("MySpace/Authentication", "RC4 Data Blob (bytes): %v", []byte(blob))

	if strings.Contains(string(blob), cli.ClientAccount.Mail) && strings.Contains(string(blob), string(verifybuf)) {

		if database.SetLastLoginDate(cli.ClientAccount.UIN) != nil {
			return false
		}

		cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
			MySpaceNewDataGeneric("lc", "2"),
			MySpaceNewDataInt("sesskey", ctx.SessionKey),
			MySpaceNewDataInt("proof", cli.ClientAccount.UIN),
			MySpaceNewDataInt("userid", cli.ClientAccount.UIN),
			MySpaceNewDataInt("profileid", cli.ClientAccount.UIN),
			MySpaceNewDataGeneric("uniquenick", cli.ClientAccount.DisplayName),
			MySpaceNewDataGeneric("id", "1"),
		}))

		cli.ClientInfo.Build = fmt.Sprintf("1.0.%s.0", clientver)
		cli.ClientInfo.Messenger = "MySpaceIM"
		cli.ClientInfo.Protocol = MySpaceIdentifyProtocolRevision(clientver)

		logging.Info("MySpace", "Client successfully authenticated! (UIN: %d, SN: %s, Mail: %s, Build: %s, Proto: %s)", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName, cli.ClientAccount.Mail, cli.ClientInfo.Build, cli.ClientInfo.Protocol)

		return true

	} else {
		cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
			MySpaceNewDataBoolean("error", true),
			MySpaceNewDataGeneric("errmsg", "The password provided is incorrect."),
			MySpaceNewDataGeneric("err", "260"),
			MySpaceNewDataBoolean("fatal", true),
		}))
	}

	return false
}

func MySpaceHandleClientLogoutRequest(data string) bool {
	return strings.HasPrefix(data, "\\logout")
}

func MySpaceHandleClientBroadcastSigninStatus(cli *network.Client, ctx *MySpaceContext) {
	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != cli.ClientAccount.UIN { // make sure we dont fuck the client up by sending it to ourselves
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", cli.ClientAccount.UIN)

			if err != nil {
				logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to get contact list for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN { // send the signon broadcast only to people on our friends list, otherwise the client will add them which is bad.
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, cli.ClientAccount.UIN)

					if err != nil {
						logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
							MySpaceNewDataInt("bm", 100),
							MySpaceNewDataInt("f", cli.ClientAccount.UIN),
							MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", ctx.Status.Code, ctx.Status.Message)),
						}))
						cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
							MySpaceNewDataInt("bm", 100),
							MySpaceNewDataInt("f", network.Clients[ix].ClientAccount.UIN),
							MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", clientContexts[ix].Status.Code, clientContexts[ix].Status.Message)),
						}))

					}
				}
			}
			row.Close()

		}
	}

	logging.System("MySpace", "Broadcasted Sign In Status for UIN: %d / SN: %s", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName)
}

func MySpaceHandleClientKeepalive(cli *network.Client) {
	for {
		time.Sleep(180 * time.Second)
		if cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{MySpaceNewDataBoolean("ka", true)})) != nil {
			break
		}
	}
}

func MySpaceHandleClientBroadcastLogoffStatus(cli *network.Client, ctx *MySpaceContext) {
	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != cli.ClientAccount.UIN {
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", cli.ClientAccount.UIN)

			if err != nil {
				logging.Error("MySpace/MySpaceHandleClientBroadcastLogoffStatus", "Failed to get contact list for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("MySpace/MySpaceHandleClientBroadcastLogoffStatus", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN {
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, cli.ClientAccount.UIN)

					if err != nil {
						logging.Error("MySpace/MySpaceHandleClientBroadcastLogoffStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
							MySpaceNewDataInt("bm", 100),
							MySpaceNewDataInt("f", cli.ClientAccount.UIN),
							MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|0|ss|%s", ctx.Status.Message)),
						}))
					}
				}
			}
			row.Close()
		}
	}

	logging.System("MySpace", "Broadcasted Logoff Status for UIN: %d / SN: %s", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName)
}

func MySpaceHandleClientOfflineMessagesDelivery(cli *network.Client, ctx *MySpaceContext) {
	row, err := database.Query("SELECT * from offlinemsgs WHERE RecvUIN= ?", cli.ClientAccount.UIN)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientOfflineMessagesDelivery", "Failed to get offline messages list for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	for row.Next() {
		var message network.OfflineMessage
		err = row.Scan(&message.SenderUIN, &message.RecvUIN, &message.MessageDate, &message.MessageContent)

		if err != nil {
			logging.Error("MySpace/MySpaceHandleClientOfflineMessagesDelivery", "Failed to scan offline messages list of uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
			row.Close()
			return
		}

		cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
			MySpaceNewDataInt("bm", 1),
			MySpaceNewDataInt("sesskey", ctx.SessionKey),
			MySpaceNewDataInt("f", message.SenderUIN),
			MySpaceNewDataInt("date", message.MessageDate),
			MySpaceNewDataGeneric("msg", message.MessageContent),
		}))
	}
	row.Close()

	row, err = database.Query("DELETE from offlinemsgs WHERE RecvUIN= ?", cli.ClientAccount.UIN)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientOfflineMessagesDelivery", "Failed to delete offline messages for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Close()

	logging.System("MySpace", "Delivered Offline Messages for UIN: %d / SN: %s", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName)
}

// packets

func MySpaceHandleClientPacketSetStatusMessage(cli *network.Client, ctx *MySpaceContext, stream string) {
	ctx.Status.Code, _ = strconv.Atoi(MySpaceRetrieveKeyValue("status", stream))
	ctx.Status.Message = MySpaceRetrieveKeyValue("statstring", stream)

	row, err := database.Query("UPDATE userdetails SET StatusMessage= ? WHERE UIN= ?", ctx.Status.Message, cli.ClientAccount.UIN)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientPacketSetStatusMessage", "Failed to update status message for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Close()

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != cli.ClientAccount.UIN { // make sure we dont fuck the client up by sending it to ourselves
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", cli.ClientAccount.UIN)

			if err != nil {
				logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to get contact list for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN { // send the signon broadcast only to people on our friends list, otherwise the client will add them which is bad.
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, cli.ClientAccount.UIN)

					if err != nil {
						logging.Error("MySpace/MySpaceHandleClientBroadcastSigninStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
							MySpaceNewDataInt("bm", 100),
							MySpaceNewDataInt("f", cli.ClientAccount.UIN),
							MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", ctx.Status.Code, ctx.Status.Message)),
						}))
					}
				}
			}
			row.Close()

		}
	}

	logging.System("MySpace", "Broadcasted New Status for UIN: %d / SN: %s (Code: %d, Message: %s)", cli.ClientAccount.UIN, cli.ClientAccount.DisplayName, ctx.Status.Code, ctx.Status.Message)
}

func MySpaceHandleClientPacketAddBuddy(cli *network.Client, ctx *MySpaceContext, stream string) {
	friend_id := MySpaceRetrieveKeyValue("newprofileid", stream)
	if friend_id == "6221" {
		// System Messages Bot cannot be added as a friend, yet
		return
	}

	var count int
	row, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", cli.ClientAccount.UIN, friend_id)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientPacketAddBuddy", "Failed to count contacts for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Next()
	row.Scan(&count)
	row.Close()

	if count > 0 {
		// already a friend, ignore
		cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
			MySpaceNewDataBoolean("error", true),
			MySpaceNewDataGeneric("errmsg", "The profile requested is already a buddy."),
			MySpaceNewDataInt("err", 1539),
		}))
		return
	}

	row, err = database.Query("INSERT into contacts (`SenderUIN`, `RecvUIN`) VALUES (?, ?)", cli.ClientAccount.UIN, friend_id)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientPacketAddBuddy", "Failed to insert new contact for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Close()

	row, err = database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", friend_id, cli.ClientAccount.UIN)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientPacketAddBuddy", "Failed to count contacts for uin: %d (friend of: %d) (%s)", friend_id, cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Next()
	row.Scan(&count)
	row.Close()

	if count > 0 {
		for ix := 0; ix < len(network.Clients); ix++ {
			if strconv.Itoa(network.Clients[ix].ClientAccount.UIN) == friend_id {
				cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
					MySpaceNewDataInt("bm", 100),
					MySpaceNewDataInt("f", network.Clients[ix].ClientAccount.UIN),
					MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", clientContexts[ix].Status.Code, clientContexts[ix].Status.Message)),
				}))
				network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
					MySpaceNewDataInt("bm", 100),
					MySpaceNewDataInt("f", cli.ClientAccount.UIN),
					MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", ctx.Status.Code, ctx.Status.Message)),
				}))
			}
		}
	}

	logging.System("MySpace", "New Buddy Added! UIN: %d -> UIN: %d", cli.ClientAccount.UIN, friend_id)
}

func MySpaceHandleClientPacketDeleteBuddy(cli *network.Client, stream string) {
	delete_id := MySpaceRetrieveKeyValue("delprofileid", stream)

	var count int
	row, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", cli.ClientAccount.UIN, delete_id)

	if err != nil {
		logging.Error("MySpace/MySpaceHandleClientPacketDeleteBuddy", "Failed to count contacts for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
		return
	}

	row.Next()
	row.Scan(&count)
	row.Close()

	if count > 0 {
		row, err := database.Query("DELETE from contacts WHERE SenderUIN= ? AND RecvUIN= ?", cli.ClientAccount.UIN, delete_id)

		if err != nil {
			logging.Error("MySpace/MySpaceHandleClientPacketDeleteBuddy", "Failed to delete contact for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
			return
		}

		row.Close()

		for ix := 0; ix < len(network.Clients); ix++ {

			row, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", cli.ClientAccount.UIN, delete_id)

			if err != nil {
				logging.Error("MySpace/MySpaceHandleClientPacketDeleteBuddy", "Failed to count contacts for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
				return
			}

			row.Next()
			row.Scan(&count)
			row.Close()

			if count > 0 {
				network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
					MySpaceNewDataInt("bm", 100),
					MySpaceNewDataInt("f", cli.ClientAccount.UIN),
					MySpaceNewDataGeneric("msg", "|s|0|ss|Offline"),
				}))
			}

		}
	}
}

// bm type 1 -> IM
func MySpaceHandleClientPacketBuddyInstantMessage(cli *network.Client, ctx *MySpaceContext, stream string) {
	isOnline := false
	msg := MySpaceRetrieveKeyValue("msg", stream)
	recv_id := MySpaceRetrieveKeyValue("t", stream)

	for ix := 0; ix < len(network.Clients); ix++ {
		if strconv.Itoa(network.Clients[ix].ClientAccount.UIN) == recv_id {
			isOnline = true
			network.Clients[ix].Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
				MySpaceNewDataInt("bm", 1),
				MySpaceNewDataInt("sesskey", clientContexts[ix].SessionKey),
				MySpaceNewDataInt("f", cli.ClientAccount.UIN),
				MySpaceNewDataGeneric("msg", msg),
			}))
		}
	}

	if !isOnline {
		if !strings.Contains(msg, "%typing%") && !strings.Contains(msg, "%stoptyping%") {
			row, err := database.Query("INSERT into offlinemsgs (`SenderUIN`, `RecvUIN`, `MessageDate`, `MessageContent`)", cli.ClientAccount.UIN, recv_id, time.Now().UTC().Unix(), msg)

			if err != nil {
				logging.Error("MySpace/MySpaceHandleClientPacketBuddyInstantMessage", "Failed to insert offline for uin: %d (%s)", recv_id, err.Error())
				return
			}

			row.Close()
		}
	}
}

func MySpaceHandleClientCallbackNetLinkRequest(cli *network.Client, stream string) {
	cmd, _ := strconv.Atoi(MySpaceRetrieveKeyValue("cmd", stream))

	//test
	res := MySpaceBuildPackage([]MySpaceDataPair{
		MySpaceNewDataBoolean("persistr", true),
		MySpaceNewDataInt("uid", cli.ClientAccount.UIN),
		MySpaceNewDataInt("cmd", cmd^256),
		MySpaceNewDataGeneric("dsn", MySpaceRetrieveKeyValue("dsn", stream)),
		MySpaceNewDataGeneric("lid", MySpaceRetrieveKeyValue("lid", stream)),
		MySpaceNewDataGeneric("rid", MySpaceRetrieveKeyValue("rid", stream)),
		MySpaceNewDataGeneric("body", MySpaceBuildInnerBody([]MySpaceDataPair{
			MySpaceNewDataGeneric("SourceURL", MySpaceEscapeString("http://google.de")),
		})),
	})
	cli.Connection.WriteTraffic(res)
}
