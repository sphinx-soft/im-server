package msim

import (
	"ghost-msim/utils"
	"strconv"
	"strings"
)

type Tuple struct {
	Item1 string
	Item2 string
}

func NewTuple(item1 string, item2 string) Tuple {
	return Tuple{Item1: item1, Item2: item2}
}

/*
const (

	msim_not_a_packet   = -2       // garbage
	msim_unknown_packet = -1       // unknown packet
	msim_error          = iota - 2 // error

	msim_login_initial  // lc 1
	msim_login_response // lc 2
	msim_keepalive      // keepalive
	msim_callback_reply // persistr

	msim_login_challenge  // login2
	msim_logout           // logout
	msim_callback_request //persist

)

	func getPacketHeader(packetType int) string {
		switch packetType {
		case msim_error:
			return "\\error"
		case msim_not_a_packet:
		case msim_unknown_packet:
			return "**invalid**"

		case msim_login_initial:
			return "\\lc\\1"
		case msim_login_response:
			return "\\lc\\2"
		case msim_keepalive:
			return "\\ka"
		case msim_callback_reply:
			return "\\persistr\\1"

		case msim_login_challenge:
			return "\\login2"
		case msim_logout:
			return "\\logout"
		case msim_callback_request:
			return "\\persist"
		}
		return ""
	}
*/
func decode(key string, packet []byte) string {
	decodedPacket := string(packet)
	splits := strings.Split(decodedPacket, "\\")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == key {
			return splits[ix+1]
		}
	}

	return ""
}

/*
	func decodeString(key string, packet string) string {
		return decode(key, []byte(packet))
	}
*/
func encode(datapairs []Tuple) string {
	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += "\\" + datapairs[i].Item1
		if datapairs[i].Item2 != "" {
			final += "\\" + datapairs[i].Item2
		}
	}
	final += "\\final\\"
	return final
}
func encodeBody(datapairs []Tuple) string {
	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += datapairs[i].Item1 + "="
		final += datapairs[i].Item2 + "\x1c"
	}
	return final
}
func handleUserLookupPacket(c client, packet []byte) {
	cmd, _ := strconv.Atoi(decode("cmd", packet))
	dsn := decode("dsn", packet)
	lid := decode("lid", packet)

	if cmd == 1 && dsn == "5" && lid == "7" {
		utils.Log("DEBUG", "user lookup packet detected")
		parsedbody := strings.Split(decode("body", packet), "=")
		res := encode([]Tuple{
			NewTuple("persistr", "1"),
			NewTuple("uid", "1"),
			NewTuple("cmd", strconv.Itoa(cmd^256)),
			NewTuple("dsn", dsn),
			NewTuple("lid", lid),
			NewTuple("rid", decode("rid", packet)),
			NewTuple("body", encodeBody([]Tuple{
				NewTuple(parsedbody[0], parsedbody[1]),
				NewTuple("UserID", "12"),
				NewTuple("ImageURL", "http:/1/1tim.is-a-failure.lol/1content/1cdn/11Bft1CL3Yv61/1Code_Z9yVNLY1tA.png"),
				NewTuple("DisplayName", "Tim2"),
				NewTuple("BandName", "soos"),
				NewTuple("SongName", "Hitler"),
				NewTuple("Age", "69"),
				NewTuple("Gender", "Attack Chopper"),
				NewTuple("Location", "Usbekistan"),
			})),
		})
		utils.Log("DEBUG", res)
		c.conn.Write([]byte(res))

	}
}
