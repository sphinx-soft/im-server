package msim

import (
	"math/rand"
	"phantom/util"
	"strconv"
	"strings"
)

func msim_new_data_string(key string, value string) msim_data_pair {
	return msim_data_pair{Key: key, Value: value}
}

func msim_new_data_int(key string, value int) msim_data_pair {
	return msim_data_pair{Key: key, Value: strconv.Itoa(value)}
}

func msim_new_data_dictonary(key string, value string) msim_data_pair {
	return msim_new_data_string(key, value)
}

func msim_new_data_boolean(key string, value bool) msim_data_pair {
	if value {
		return msim_data_pair{Key: key, Value: "1"}
	} else {
		return msim_data_pair{Key: key, Value: "1"}
	}
}

func findValueFromKey(key string, packet []byte) string {

	decodedPacket := string(packet)
	splits := strings.Split(decodedPacket, "\\")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == key {
			return splits[ix+1]
		}
	}

	return ""
}

func buildDataPacket(datapairs []msim_data_pair) string {

	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += "\\" + datapairs[i].Key
		final += "\\" + datapairs[i].Value
	}
	final += "\\final\\"
	return final
}

func buildDataBody(datapairs []msim_data_pair) string {

	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += datapairs[i].Key + "="
		final += datapairs[i].Value + "\x1c"
	}
	return final
}

func GenerateNonce() string {

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 0x40)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateSessionKey() int {
	return rand.Intn(100000)
}

func getUserData(username string) Account {

	var acc Account

	row, _ := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE username= ?", username)
	row.Next()
	row.Scan(&acc.Uid, &acc.Username, &acc.Password, &acc.Screenname, &acc.Avatar,
		&acc.BandName, &acc.SongName, &acc.Age, &acc.Gender, &acc.Location)
	row.Close()

	return acc
}

func getUserDataById(userid int) Account {

	var acc Account

	row, _ := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE id= ?", userid)
	row.Next()
	row.Scan(&acc.Uid, &acc.Username, &acc.Password, &acc.Screenname, &acc.Avatar,
		&acc.BandName, &acc.SongName, &acc.Age, &acc.Gender, &acc.Location)
	row.Close()

	return acc
}

func escapeString(data string) string {
	res := strings.Replace(data, "/", "/1", -1)
	res = strings.Replace(res, "\\", "\\2", -1)
	return res
}

func ArrayRemove(s []*Msim_client, i int) []*Msim_client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
