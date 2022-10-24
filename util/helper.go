package util

import (
	"crypto/rc4"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func ConvertToUtf16(input string) []byte {
	bytes, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte(input))
	if err != nil {
		return nil
	}
	return bytes
}

func DecryptRC4(pwd []byte, data []byte) []byte {
	c, err := rc4.NewCipher(pwd)
	if err != nil {
		Error("Error creating RC4 Ciphertext")
		return nil
	}
	crypted := make([]byte, len(data))
	c.XORKeyStream(crypted, data)
	return crypted
}

func EncryptRC4(pwd []byte, data []byte) []byte {
	return DecryptRC4(pwd, data)
}

func AddGlobalClient(client *Global_Client) {
	Global_Clients = append(Global_Clients, client)
}

func GetClientFromGlobalList(username string) *Global_Client {
	for i := 0; i < len(Global_Clients); i++ {
		if Global_Clients[i].Username == username {
			return Global_Clients[i]
		}
	}

	return nil
}
