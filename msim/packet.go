package msim

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"phantom/global"
	"phantom/util"
	"strconv"
	"strings"
	"time"
)

var pfproot string = util.GetRootUrl()

/*
	msim_not_a_packet   = -2       	-> garbage
	msim_unknown_packet = -1       	-> unknown packet
	msim_error          = iota - 2 	-> error

	msim_login_initial  			-> lc 1
	msim_login_response 			-> lc 2
	msim_keepalive      			-> keepalive
	msim_callback_reply 			-> persistr

	msim_login_challenge  			-> login2
	msim_logout           			-> logout
	msim_callback_request 			-> persist
*/

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

			if strings.Contains(str, "\\dsn\\6") && strings.Contains(str, "\\lid\\11") {
				handleClientPacketRequestNetLink(client, data)
			}

			if strings.Contains(str, "\\dsn\\7") && strings.Contains(str, "\\lid\\18") {
				handleClientPacketNewNotificationRequest(client, data)
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
		handleClientPacketBuddyInstantMessage(client, ctx, data)
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

// login
func handleClientAuthentication(client *global.Client, ctx *msim_context) bool {
	util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
		msim_new_data_string("lc", "1"),
		msim_new_data_string("nc", base64.StdEncoding.EncodeToString([]byte(ctx.nonce))),
		msim_new_data_string("id", "1"),
	}))

	loginpacket, success := util.ReadTraffic(client.Connection)
	if !success {
		util.Error("MySpace -> handleClientAuthentication", "Failed to read Login2 Data Packet!")
		return false
	}

	username := findValueFromKey("username", loginpacket)
	version := findValueFromKey("clientver", loginpacket)

	acc, _ := global.GetUserDataFromUsername(username)
	client.Account = acc
	client.Protocol = identifyProtocolVersion(version)

	uid := acc.UserId
	screenname := acc.Screenname
	password := strings.Replace(util.DecryptAES(util.GetAESKey(), acc.Password), "\r\n", "", -1)

	util.Debug("MySpace -> handleClientAuthentication", "rc4 pw test: %v", []byte(password))

	byte_nc2 := make([]byte, 32)
	byte_rc4_key := make([]byte, 16)
	byte_challenge := []byte(ctx.nonce)
	for i := 0; i < 32; i++ {
		byte_nc2[i] = byte_challenge[i+32]
	}

	byte_password := util.ConvertToUtf16(password)
	hasher := sha1.New()
	hasher.Write(byte_password)
	byte_hash_phase1 := hasher.Sum(nil)

	util.Debug("MySpace -> handleClientAuthentication", "sha1 pw test1: %v", byte_hash_phase1)

	byte_hash_phase2 := append(byte_hash_phase1, byte_nc2...)
	hasher.Reset()
	hasher.Write(byte_hash_phase2)
	byte_hash_total := hasher.Sum(nil)
	hasher.Reset()

	util.Debug("MySpace -> handleClientAuthentication", "sha1 pw test2: %v", byte_hash_phase2)

	for i := 0; i < 16; i++ {
		byte_rc4_key[i] = byte_hash_total[i]
	}
	packetrc4data := findValueFromKey("response", loginpacket)
	byte_rc4_data, err := base64.StdEncoding.DecodeString(packetrc4data)
	if err != nil {
		util.Error("MySpace -> handleClientAuthentication", "Invalid base64 provided at login packet.")
		return false
	}
	rc4data := util.DecryptRC4(byte_rc4_key, byte_rc4_data)
	util.Debug("MySpace -> handleClientAuthentication", "rc4 data test: %v", rc4data)
	util.Debug("MySpace -> handleClientAuthentication", "rc4 data test: %s", string(rc4data))

	if strings.Contains(string(rc4data), username) {
		res, _ := util.GetDatabaseHandle().Query("UPDATE myspace SET lastlogin = ? WHERE id= ?", time.Now().UnixNano(), acc.UserId)
		res.Close()
		util.Log("MySpaceIM", "Client Authenticated! -> Username: %s, Screenname: %s, Version: 1.0.%s.0, Protocol Version: %s", username, screenname, version, client.Protocol)
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_string("lc", "2"),
			msim_new_data_int("sesskey", ctx.sesskey),
			msim_new_data_int("proof", uid),
			msim_new_data_int("userid", uid),
			msim_new_data_int("profileid", uid),
			msim_new_data_string("uniquenick", screenname),
			msim_new_data_string("id", "1"),
		}))

		client.BuildNumber = fmt.Sprintf("1.0.%s.0", version)
		client.Protocol = "MSIMv?"

		return true
	} else {
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("error", true),
			msim_new_data_string("errmsg", "The password provided is incorrect."),
			msim_new_data_string("err", "260"),
			msim_new_data_boolean("fatal", true),
		}))
	}
	return false
}

// broadcast sign on status
func handleClientBroadcastSignOnStatus(client *global.Client, ctx *msim_context) {
	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.UserId != client.Account.UserId {
			res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE from_id= ?", client.Account.UserId)
			for res.Next() {
				var msg global.Contact
				_ = res.Scan(&msg.FromId, &msg.ToId)
				if global.Clients[i].Account.UserId == msg.ToId {
					res2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE from_id= ? AND to_id= ?", global.Clients[i].Account.UserId, client.Account.UserId)
					res2.Next()
					var count int
					res2.Scan(&count)
					res2.Close()
					if count > 0 {
						util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", client.Account.UserId),
							msim_new_data_string("msg", fmt.Sprintf("|s|%d|ss|%s", ctx.statuscode, ctx.statusmessage)),
							//msim_new_data_string("msg", "|s|"+ctx.StatusCode+"|ss|"+client.StatusText),
						}))
						util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", global.Clients[i].Account.UserId),
							msim_new_data_string("msg", fmt.Sprintf("|s|%d|ss|%s", users_context[i].statuscode, users_context[i].statusmessage)),
							//msim_new_data_string("msg", "|s|"+Msim_Clients[i].StatusCode+"|ss|"+Msim_Clients[i].StatusText),
						}))
					}
				}
			}
			res.Close()
		}
	}
}

// broadcast sign off events
func handleClientBroadcastSignOffStatus(client *global.Client, ctx *msim_context) {
	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.UserId != client.Account.UserId {
			res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE from_id= ?", client.Account.UserId)
			for res.Next() {
				var msg global.Contact
				_ = res.Scan(&msg.FromId, &msg.ToId)
				if global.Clients[i].Account.UserId == msg.ToId {
					res2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE from_id= ? AND to_id= ?", global.Clients[i].Account.UserId, client.Account.UserId)
					res2.Next()
					var count int
					res2.Scan(&count)
					res2.Close()
					if count > 0 {
						util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", client.Account.UserId),
							msim_new_data_string("msg", fmt.Sprintf("|s|0|ss|%s", ctx.statusmessage)),
						}))
					}

				}
			}
			res.Close()
		}
	}
}

// handle offline messages
func handleClientHandleOfflineMessages(client *global.Client, ctx *msim_context) {
	res, _ := util.GetDatabaseHandle().Query("SELECT * from offlinemsgs WHERE to_id= ?", client.Account.UserId)
	for res.Next() {
		var msg global.OfflineMsg
		_ = res.Scan(&msg.FromId, &msg.ToId, &msg.Message, &msg.Date)
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_int("bm", 1),
			msim_new_data_int("sesskey", ctx.sesskey),
			msim_new_data_int("f", msg.FromId),
			//	msim_new_data_int("date", msg.date),
			msim_new_data_string("msg", msg.Message),
		}))
		util.Debug("MySpace -> handleClientOfflineEvents", "%d", msg.Date)
	}
	res.Close()
	res2, _ := util.GetDatabaseHandle().Query("DELETE from offlinemsgs WHERE to_id= ?", client.Account.UserId)
	res2.Close()
}

func handleClientLogoutRequest(data string) bool {
	if strings.HasPrefix(data, "\\logout") {
		return true
	} else {
		return false
	}
}

// Status Messages
func handleClientPacketSetStatusMessages(client *global.Client, ctx *msim_context, packet []byte) {
	status := findValueFromKey("status", packet)
	statstring := findValueFromKey("statstring", packet)

	ctx.statuscode, _ = strconv.Atoi(status)
	ctx.statusmessage = statstring
	res, _ := util.GetDatabaseHandle().Query("UPDATE myspace SET headline= ? WHERE id= ?", statstring, client.Account.UserId)
	res.Close()
	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.UserId != client.Account.UserId {
			res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE from_id= ?", client.Account.UserId)
			for res.Next() {
				var msg global.Contact
				_ = res.Scan(&msg.FromId, &msg.ToId)
				if global.Clients[i].Account.UserId == msg.ToId {
					res2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE from_id= ? AND to_id= ?", global.Clients[i].Account.UserId, client.Account.UserId)
					res2.Next()
					var count int
					res2.Scan(&count)
					res2.Close()
					if count > 0 {
						util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
							msim_new_data_int("bm", 100),
							msim_new_data_int("f", client.Account.UserId),
							msim_new_data_string("msg", fmt.Sprintf("|s|%s|ss|%s", status, statstring)),
							//msim_new_data_string("msg", "|s|"+status+"|ss|"+statstring+""),
						}))
					}

				}
			}
			res.Close()
		}
	}
}

// addbuddy message
func handleClientPacketAddBuddy(client *global.Client, ctx *msim_context, packet []byte) {
	if findValueFromKey("newprofileid", packet) == "6221" {
		util.Debug("MySpace -> handleClientPacketAddBuddy", "MySpace Chatbot Friend Request Detected! Skipping...")
		return
	}
	newprofileid := findValueFromKey("newprofileid", packet)

	var count int
	check, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE to_id=? and from_id= ?", newprofileid, client.Account.UserId)
	check.Next()
	check.Scan(&count)
	check.Close()
	if count > 0 {
		util.Debug("MySpace -> handleClientPacketAddBuddy", "Buddy is already added to Contact List! Returning Error...")
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("error", true),
			msim_new_data_string("errmsg", "The profile requested is already a buddy."),
			msim_new_data_int("err", 1539),
		}))
		return
	}
	util.Debug("addbuddy", "%d:%d", client.Account.UserId, newprofileid)
	dbres, _ := util.GetDatabaseHandle().Query("INSERT into contacts (`from_id`, `to_id`) VALUES (?, ?)", client.Account.UserId, newprofileid)
	dbres.Close()
	var count2 int
	check2, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE from_id=? and to_id= ?", newprofileid, client.Account.UserId)
	check2.Next()
	check2.Scan(&count2)
	check2.Close()
	if count2 > 0 {
		for i := 0; i < len(global.Clients); i++ {
			if strconv.Itoa(global.Clients[i].Account.UserId) == newprofileid {
				util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
					msim_new_data_int("bm", 100),
					msim_new_data_int("f", global.Clients[i].Account.UserId),
					msim_new_data_string("msg", fmt.Sprintf("|s|%d|ss|%s", users_context[i].statuscode, users_context[i].statusmessage)),
					//msim_new_data_string("msg", "|s|"+[i].StatusCode+"|ss|"+Msim_Clients[i].StatusText),
				}))
				util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
					msim_new_data_int("bm", 100),
					msim_new_data_int("f", client.Account.UserId),
					msim_new_data_string("msg", fmt.Sprintf("|s|%d|ss|%s", ctx.statuscode, ctx.statusmessage)),
				}))
			}
		}
	}
}

// delbuddy message
func handleClientPacketDelBuddy(client *global.Client, packet []byte) {
	delprofileid := findValueFromKey("delprofileid", packet)
	dbres, _ := util.GetDatabaseHandle().Query("DELETE from contacts WHERE to_id=? and from_id= ?", delprofileid, client.Account.UserId)
	dbres.Close()
	for i := 0; i < len(global.Clients); i++ {
		if strconv.Itoa(global.Clients[i].Account.UserId) == delprofileid {
			var count int
			dbres, _ := util.GetDatabaseHandle().Query("SELECT COUNT(*) from contacts WHERE to_id=? and from_id= ?", client.Account.UserId, delprofileid)
			dbres.Next()
			dbres.Scan(&count)
			dbres.Close()
			if count > 0 {
				util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
					msim_new_data_int("bm", 100),
					msim_new_data_int("f", client.Account.UserId),
					msim_new_data_string("msg", "|s|0|ss|Offline"),
				}))
			}
		}
	}
}

// bm type 1
func handleClientPacketBuddyInstantMessage(client *global.Client, ctx *msim_context, packet []byte) {
	t, _ := strconv.Atoi(findValueFromKey("t", packet))
	msg := findValueFromKey("msg", packet)
	date := time.Now().UTC().UnixMilli()
	found := false
	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.UserId == t {
			found = true
			util.WriteTraffic(global.Clients[i].Connection, buildDataPacket([]msim_data_pair{
				msim_new_data_int("bm", 1),
				msim_new_data_int("sesskey", users_context[i].sesskey),
				msim_new_data_int("f", client.Account.UserId),
				msim_new_data_string("msg", msg),
			}))
		}
	}
	if !found {
		if !strings.Contains(msg, "%typing%") && !strings.Contains(msg, "%stoptyping%") {
			res, _ := util.GetDatabaseHandle().Query("INSERT INTO offlinemsgs (`from_id`, `to_id`, `message`, `date`) VALUES (?, ?, ?, ?)", client.Account.UserId, t, msg, date)
			res.Close()
		}
	}
}

// persist 1;0;1 get_contact_information
func handleClientPacketGetContactList(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)
	util.Debug("MySpace -> handleClientPacketGetContactList", "Requested Contact List...")
	res, _ := util.GetDatabaseHandle().Query("SELECT * from contacts WHERE from_id=?", client.Account.UserId)
	body := ""
	for res.Next() {
		var contact global.Contact
		_ = res.Scan(&contact.FromId, &contact.ToId)

		accountRow, _ := global.GetUserDataFromUserId(contact.ToId)
		accountData, _ := getMySpaceDataByUserId(contact.ToId)

		body += buildDataBody([]msim_data_pair{
			msim_new_data_int("ContactID", accountRow.UserId),
			msim_new_data_string("Headline", accountData.headline),
			msim_new_data_int("Position", 1),                //TODO
			msim_new_data_string("GroupName", "IM Friends"), //TODO
			msim_new_data_int("Visibility", 1),
			msim_new_data_string("ShowAvatar", "true"),
			msim_new_data_string("AvatarUrl", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_int64("LastLogin", accountData.lastlogin),
			msim_new_data_string("IMName", accountRow.Email),
			msim_new_data_string("NickName", accountRow.Screenname),
			msim_new_data_int("NameSelect", 0),
			msim_new_data_string("OfflineMsg", "im offline"),
			msim_new_data_int("SkyStatus", 0),
		})
	}
	res.Close()
	resp := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", body),
	})
	util.WriteTraffic(client.Connection, resp)
}

// persist 1;0;2 get_contact_information
func handleClientPacketGetContactInformation(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	parsedbody := strings.Split(findValueFromKey("body", packet), "=")

	util.Debug("MySpace -> handleClientPacketGetContactInformation", "Requesting Contact Information...")
	parse, _ := strconv.Atoi(parsedbody[1])

	accountRow, _ := global.GetUserDataFromUserId(parse)
	accountData, _ := getMySpaceDataByUserId(parse)
	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_int("ContactID", accountRow.UserId),
			msim_new_data_string("Headline", accountData.headline),
			msim_new_data_int("Position", 1),                 //TODO
			msim_new_data_string("!GroupName", "IM Friends"), //TODO
			msim_new_data_int("Visibility", 1),
			msim_new_data_string("!ShowAvatar", "true"),
			msim_new_data_string("!AvatarUrl", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_int("!NameSelect", 0),
			msim_new_data_string("IMName", accountRow.Email),
			msim_new_data_string("!NickName", accountRow.Screenname),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;1;4
func handleClientPacketUserLookupIMAboutMyself(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	parse := client.Account.UserId

	accountRow, _ := global.GetUserDataFromUserId(parse)
	accountData, _ := getMySpaceDataByUserId(parse)

	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_int("UserID", accountRow.UserId),
			msim_new_data_string("Sound", "true"),
			msim_new_data_int("!PrivacyMode", 0),
			msim_new_data_string("!ShowOnlyToList", "False"),
			msim_new_data_int("!OfflineMessageMode", 2),
			msim_new_data_string("Headline", accountData.headline),
			msim_new_data_string("Avatarurl", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_int("Alert", 1),
			msim_new_data_string("!ShowAvatar", "true"),
			msim_new_data_string("IMName", accountRow.Screenname),
			msim_new_data_int("!ClientVersion", 999),
			msim_new_data_string("!AllowBrowse", "true"),
			msim_new_data_string("IMLang", "English"),
			msim_new_data_int("LangID", 8192),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;1;17
func handleClientPacketUserLookupIMByUid(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	parsedbody := strings.Split(findValueFromKey("body", packet), "=")
	parse, _ := strconv.Atoi(parsedbody[1])

	accountRow, _ := global.GetUserDataFromUserId(parse)
	accountData, _ := getMySpaceDataByUserId(parse)

	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_int("UserID", accountRow.UserId),
			msim_new_data_string("Sound", "true"),
			msim_new_data_int("!PrivacyMode", 0),             // TODO
			msim_new_data_string("!ShowOnlyToList", "False"), // TODO
			msim_new_data_int("!OfflineMessageMode", 2),      // TODO
			msim_new_data_string("Headline", accountData.headline),
			msim_new_data_string("Avatarurl", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_int("Alert", 1),               //TODO
			msim_new_data_string("!ShowAvatar", "true"), // TODO
			msim_new_data_string("IMName", accountRow.Screenname),
			msim_new_data_int("!ClientVersion", 999),
			msim_new_data_string("!AllowBrowse", "true"), // TODO
			msim_new_data_string("IMLang", "English"),
			msim_new_data_int("LangID", 8192),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// persist 1;2;6
// \persist\1\sesskey\7920\cmd\1\dsn\2\uid\1\lid\6\rid\8\body\\final\
func handleClientPacketGetGroups(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	util.Debug("MySpace -> handleClientPacketGetGroups", "Requesting Contact Groups")
	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_int("GroupID", 21672248),
			msim_new_data_string("GroupName", "IM Friends"),
			msim_new_data_int("Position", 1),
			msim_new_data_int("GroupFlag", 131073),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;4;3, 1;4;5
func handleClientPacketUserLookupMySpaceByUid(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)
	parsedbody := strings.Split(findValueFromKey("body", packet), "=")

	parse, _ := strconv.Atoi(parsedbody[1])
	accountRow, _ := global.GetUserDataFromUserId(parse)
	accountData, _ := getMySpaceDataByUserId(parse)

	util.Debug("MySpace -> handleClientPacketUserLookupMySpaceByUid", "http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype)

	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_string("UserName", accountRow.Email),
			msim_new_data_int("UserID", accountRow.UserId),
			msim_new_data_string("ImageURL", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_string("DisplayName", accountRow.Screenname),
			msim_new_data_string("BandName", accountData.bandname),
			msim_new_data_string("SongName", accountData.songname),
			msim_new_data_int("Age", accountData.age),
			msim_new_data_string("Gender", accountData.gender),
			msim_new_data_string("Location", accountData.location),
			msim_new_data_int("!TotalFriends", 1), //TODO
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;5;7
func handleClientPacketUserLookupMySpaceByUsernameOrEmail(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	parsedbody := strings.Split(findValueFromKey("body", packet), "=")
	accountRow, _ := global.GetUserDataFromEmail(parsedbody[1])
	accountData, _ := getMySpaceDataByEmail(parsedbody[1])

	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_string(parsedbody[0], parsedbody[1]),
			msim_new_data_int("UserID", accountRow.UserId),
			msim_new_data_string("ImageURL", escapeString(fmt.Sprintf("http://%s/pfp/id=%d.%s", pfproot, accountRow.UserId, accountData.avatartype))),
			msim_new_data_string("DisplayName", accountRow.Screenname),
			msim_new_data_string("BandName", accountData.bandname),
			msim_new_data_string("SongName", accountData.songname),
			msim_new_data_int("Age", accountData.age),
			msim_new_data_string("Gender", accountData.gender),
			msim_new_data_string("Location", accountData.location),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;6;11
func handleClientPacketRequestNetLink(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))

	//test
	res := buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", findValueFromKey("dsn", packet)),
		msim_new_data_string("lid", findValueFromKey("lid", packet)),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
			msim_new_data_string("!URL", escapeString("http://google.de")),
		})),
	})
	util.WriteTraffic(client.Connection, res)
}

// Persist 1;7;18
func handleClientPacketNewNotificationRequest(client *global.Client, packet []byte) {
	/*	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))

		//test
		res := buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("persistr", true),
			msim_new_data_int("uid", client.Account.UserId),
			msim_new_data_int("cmd", cmd^256),
			msim_new_data_string("dsn", findValueFromKey("dsn", packet)),
			msim_new_data_string("lid", findValueFromKey("lid", packet)),
			msim_new_data_string("rid", findValueFromKey("rid", packet)),
			msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
				msim_new_data_string("FriendRequest", "On"),
			})),
		})
		util.WriteTraffic(client.Connection, res)
	*/
}

var buf []byte

// persist 514;8;13 2;8;13 change_profile_picture
func handleClientPacketChangePicture(client *global.Client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	body := strings.Split(strings.Replace(findValueFromKey("body", packet), "\x1c", "", -1), "=")
	part, err := base64.StdEncoding.DecodeString(unescapeString(body[len(body)-1]))
	if err != nil {
		return
	}
	buf = append(buf, part...)

	if strings.Contains(findValueFromKey("body", packet), "True") {
		var pfpType string
		if strings.HasPrefix(string(buf), "GIF") {
			pfpType = "gif"
		} else if strings.Contains(string(buf), "PNG") {
			pfpType = "png"
		} else {
			pfpType = "jpg"
		}
		res, _ := util.GetDatabaseHandle().Query("UPDATE upload SET avatar= ? WHERE id= ?", base64.StdEncoding.EncodeToString(buf), client.Account.UserId)
		res.Close()
		res, _ = util.GetDatabaseHandle().Query("UPDATE myspace SET avatartype= ? WHERE id= ?", pfpType, client.Account.UserId)
		res.Close()
		buf = nil
	}

	util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
		msim_new_data_boolean("persistr", true),
		msim_new_data_int("uid", client.Account.UserId),
		msim_new_data_int("cmd", cmd^256),
		msim_new_data_string("dsn", dsn),
		msim_new_data_string("lid", lid),
		msim_new_data_string("rid", findValueFromKey("rid", packet)),
		msim_new_data_dictonary("body", ""),
	}))
}
