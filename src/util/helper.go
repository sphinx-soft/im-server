package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
	"io"

	externalip "github.com/glendc/go-external-ip"
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
		Log(INFO, "RC4 -> SwapDecryptionState", "Error creating RC4 Ciphertext")
		return nil
	}
	crypted := make([]byte, len(data))
	c.XORKeyStream(crypted, data)
	return crypted
}

func EncryptRC4(pwd []byte, data []byte) []byte {
	return DecryptRC4(pwd, data)
}

func EncryptAES(key string, message string) string {
	plainText := []byte(message)
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		Log(INFO, "AES -> Encryption", "Error Encryption AES #1: %s", err.Error())
		return ""
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		Log(INFO, "AES -> Encryption", "Error Encryption AES #2")
		return ""
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText)
}

func DecryptAES(key string, secure string) string {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		Log(INFO, "AES -> Decryption", "Error Encryption AES #1: %s", err.Error())
		return ""
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(key))

	//IF NewCipher failed, exit:
	if err != nil {
		Log(INFO, "AES -> Decryption", "Error Encryption AES #2")
		return ""
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		Log(INFO, "AES -> Decryption", "Error Encryption AES #3 (short)")
		return ""
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func HashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn := externalip.DefaultConsensus(nil, nil)
	ip, err := conn.ExternalIP()

	if err != nil {
		Log(WARN, "Failed to get Outbound IP: %s", err.Error())
		return "<Unknown>"
	}

	return ip.String()
}
