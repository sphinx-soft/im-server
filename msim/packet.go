package msim

import (
	"crypto/sha1"
	"encoding/base64"
	"phantom/util"
	"strconv"
	"strings"
)

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

// login
func handleClientAuthentication(client Msim_client) bool {
	util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
		msim_new_data_string("lc", "1"),
		msim_new_data_string("nc", base64.StdEncoding.EncodeToString([]byte(client.Nonce))),
		msim_new_data_string("id", "1"),
	}))

	loginpacket, success := util.ReadTraffic(client.Connection)
	if !success {
		util.Error("Failed to read Login2 Data Packet!")
		return false
	}

	username := findValueFromKey("username", loginpacket)
	version := findValueFromKey("clientver", loginpacket)

	acc := getUserData(username)
	client.Account = acc

	uid := acc.Uid
	sessionkey := GenerateSessionKey()
	screenname := acc.Screenname
	password := acc.Password

	byte_nc2 := make([]byte, 32)
	byte_rc4_key := make([]byte, 16)
	byte_challenge := []byte(client.Nonce)
	for i := 0; i < 32; i++ {
		byte_nc2[i] = byte_challenge[i+32]
	}
	byte_password := util.ConvertToUtf16(password)
	hasher := sha1.New()
	hasher.Write(byte_password)
	byte_hash_phase1 := hasher.Sum(nil)

	byte_hash_phase2 := append(byte_hash_phase1, byte_nc2...)
	hasher.Reset()
	hasher.Write(byte_hash_phase2)
	byte_hash_total := hasher.Sum(nil)
	hasher.Reset()

	for i := 0; i < 16; i++ {
		byte_rc4_key[i] = byte_hash_total[i]
	}
	packetrc4data := findValueFromKey("response", loginpacket)
	byte_rc4_data, err := base64.StdEncoding.DecodeString(packetrc4data)
	if err != nil {
		util.Error("Invalid base64 provided at login packet.")
		return false
	}
	rc4data := util.DecryptRC4(byte_rc4_key, byte_rc4_data)

	if strings.Contains(string(rc4data), username) {
		util.Log("MySpaceIM", "Client Authenticated! | Username: %s | Screenname: %s | Version: 1.0.%s.0", username, screenname, version)
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_string("lc", "2"),
			msim_new_data_int("sesskey", sessionkey),
			msim_new_data_int("proof", sessionkey),
			msim_new_data_int("userid", uid),
			msim_new_data_int("profileid", uid),
			msim_new_data_string("uniquenick", screenname),
			msim_new_data_string("id", "1"),
		}))

		return true
	} else {
		util.WriteTraffic(client.Connection, buildDataPacket([]msim_data_pair{
			msim_new_data_boolean("error"),
			msim_new_data_string("errmsg", "The password provided is incorrect."),
			msim_new_data_string("err", "260"),
			msim_new_data_boolean("fatal"),
		}))
	}
	return false
}

// handleUserLookupPacket
func handleClientPacketUserLookup(client Msim_client, packet []byte) {
	cmd, _ := strconv.Atoi(findValueFromKey("cmd", packet))
	dsn := findValueFromKey("dsn", packet)
	lid := findValueFromKey("lid", packet)

	if cmd == 1 && dsn == "5" && lid == "7" {
		parsedbody := strings.Split(findValueFromKey("body", packet), "=")
		res := buildDataPacket([]msim_data_pair{
			msim_new_data_string("persistr", "1"),
			msim_new_data_string("uid", "1"),
			msim_new_data_int("cmd", cmd^256),
			msim_new_data_string("dsn", dsn),
			msim_new_data_string("lid", lid),
			msim_new_data_string("rid", findValueFromKey("rid", packet)),
			msim_new_data_dictonary("body", buildDataBody([]msim_data_pair{
				msim_new_data_string(parsedbody[0], parsedbody[1]),
				msim_new_data_string("UserID", "12"),
				msim_new_data_string("ImageURL", "http:/1/1tim.is-a-failure.lol/1content/1cdn/11Bft1CL3Yv61/1Code_Z9yVNLY1tA.png"),
				msim_new_data_string("DisplayName", "Tim2"),
				msim_new_data_string("BandName", "soos"),
				msim_new_data_string("SongName", "Hitler"),
				msim_new_data_string("Age", "69"),
				msim_new_data_string("Gender", "Attack Chopper"),
				msim_new_data_string("Location", "Usbekistan"),
			})),
		})
		util.WriteTraffic(client.Connection, res)

	}
}
