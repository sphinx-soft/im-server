package msim

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"ghost-msim/utils"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type client struct {
	conn        net.Conn
	nonce       string
	displayname string
}

func Initialise(port int) {
	l, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(port))
	if err != nil {
		utils.Error("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	utils.Log("TcpServer", "Listening on 0.0.0.0:%s", strconv.Itoa(port))
	for {
		conn, err := l.Accept()
		if err != nil {
			utils.Error("Error accepting: ", err.Error())
		}
		c := client{
			conn:  conn,
			nonce: generateNonce(),
		}
		go handleClient(c)
		go handleKeepAlive(c)
	}
}
func handleKeepAlive(c client) {
	for {
		time.Sleep(180 * time.Second)
		utils.Log("DEBUG", "Sending keep alive packet")
		_, err := c.conn.Write([]byte("\\ka\\\\\\final\\"))
		if err != nil {
			break
		}
	}
}
func login(c client) bool {
	c.conn.Write([]byte(encode([]Tuple{
		NewTuple("lc", "1"),
		NewTuple("nc", b64.StdEncoding.EncodeToString([]byte(c.nonce))),
		NewTuple("id", "1"),
	})))
	loginpacket := make([]byte, 4096)
	_, err := c.conn.Read(loginpacket)
	if err != nil {
		utils.Error("Error reading login packet from client: ", err.Error())
		return false
	}
	username := decode("username", loginpacket)
	version := decode("clientver", loginpacket)

	uid := 2
	sessionkey := rand.Intn(100000)
	screenname := "EinTim"
	password := "11ee..22"

	byte_nc2 := make([]byte, 32)
	byte_rc4_key := make([]byte, 16)
	byte_challenge := []byte(c.nonce)
	for i := 0; i < 32; i++ {
		byte_nc2[i] = byte_challenge[i+32]
	}
	byte_password := utils.ConvertToUtf16(password)
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
	packetrc4data := decode("response", loginpacket)
	byte_rc4_data, err := b64.StdEncoding.DecodeString(packetrc4data)
	if err != nil {
		utils.Error("Invalid base64 provided at login packet.")
		return false
	}
	rc4data := utils.DecryptRC4(byte_rc4_key, byte_rc4_data)

	if strings.Contains(string(rc4data), username) {
		utils.Log("MSIM", "User authed successfully with username: %s and client version: %s", username, version)
		c.conn.Write([]byte(encode([]Tuple{
			NewTuple("lc", "2"),
			NewTuple("sesskey", strconv.Itoa(sessionkey)),
			NewTuple("proof", strconv.Itoa(sessionkey)),
			NewTuple("userid", strconv.Itoa(uid)),
			NewTuple("profileid", strconv.Itoa(uid)),
			NewTuple("uniquenick", screenname),
			NewTuple("id", "1"),
		})))
		return true
	} else {
		c.conn.Write([]byte(encode([]Tuple{
			NewTuple("error", ""),
			NewTuple("\\errmsg", "The password provided is incorrect."),
			NewTuple("err", "260"),
			NewTuple("fatal", "\\"),
		})))
	}
	return false
}

func handleClient(c client) {
	utils.Log("TcpServer", "Accepted client. Awaiting MSIM authentication")
	if !login(c) {
		c.conn.Close()
		return
	}
	for {
		buf := make([]byte, 4096)
		_, err := c.conn.Read(buf)
		utils.Log("RECV", string(buf))
		handleUserLookupPacket(c, buf)
		if err != nil {
			utils.Error("Error reading message from client: ", err.Error())
			break
		}
	}
	utils.Log("TcpServer", "Client with displayname: %s disconnected", c.displayname)
	c.conn.Close()
}
