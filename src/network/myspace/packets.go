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
	"strings"
)

// packets

func MySpaceHandleClientAuthentication(cli *network.Client, ctx *MySpaceContext) bool {
	cli.Connection.WriteTraffic(MySpaceBuildPackage([]MySpaceDataPair{
		MySpaceNewDataGeneric("lc", "1"),
		MySpaceNewDataGeneric("nc", base64.StdEncoding.EncodeToString([]byte(ctx.Nonce))),
		MySpaceNewDataGeneric("id", "1"),
	}))

	loginPacket, err := cli.Connection.ReadTraffic()
	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to read Login2 data packet!")
		return false
	}

	email := MySpaceRetrieveKeyValue("username", loginPacket) // i have no clue why MySpace called this a "Username" when its the email bruh
	clientver := MySpaceRetrieveKeyValue("clientver", loginPacket)
	response := MySpaceRetrieveKeyValue("response", loginPacket)

	cli.ClientAccount, err = database.GetAccountDataByEmail(email)
	if err != nil {
		logging.Error("MySpace/Authentication", "Failed to fetch Account Data!")
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
		logging.Error("MySpace/Authentication", "Invalid base64 provided at login packet.")
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
