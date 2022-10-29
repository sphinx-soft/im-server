package msim

import (
	"math/rand"
	"phantom/global"
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

func msim_new_data_int64(key string, value int64) msim_data_pair {
	return msim_data_pair{Key: key, Value: strconv.FormatInt(value, 10)}
}

func msim_new_data_dictonary(key string, value string) msim_data_pair {
	return msim_new_data_string(key, value)
}

func msim_new_data_boolean(key string, value bool) msim_data_pair {
	if value {
		return msim_data_pair{Key: key, Value: "1"}
	} else {
		return msim_data_pair{Key: key, Value: ""}
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
		if datapairs[i].Value != "" {
			final += "\\" + datapairs[i].Value
		}
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

func generateNonce() string {

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 0x40)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateSessionKey() int {
	return rand.Intn(100000)
}

func getMySpaceDataByEmail(email string) (msim_user_details, bool) {
	var user msim_user_details

	acc, _ := global.GetUserDataFromEmail(email)

	row, err := util.GetDatabaseHandle().Query("SELECT * from myspace WHERE id= ?", acc.UserId)
	if err != nil {
		util.Error("MySpace -> getMySpaceDataByEmail", err.Error())
		return user, true
	}

	row.Next()
	row.Scan(&user.userid, &user.avatartype, &user.bandname, &user.songname, &user.age, &user.gender, &user.location, &user.headline, &user.lastlogin)
	row.Close()

	return user, true
}

func getMySpaceDataByUserId(uid int) (msim_user_details, bool) {
	var user msim_user_details

	acc, _ := global.GetUserDataFromUserId(uid)

	row, err := util.GetDatabaseHandle().Query("SELECT * from myspace WHERE id= ?", acc.UserId)
	if err != nil {
		util.Error("MySpace -> getMySpaceDataByUserId", err.Error())
		return user, true
	}
	row.Next()
	row.Scan(&user.userid, &user.avatartype, &user.bandname, &user.songname, &user.age, &user.gender, &user.location, &user.headline, &user.lastlogin)
	row.Close()

	return user, true
}

func escapeString(data string) string {
	res := strings.Replace(data, "/", "/1", -1)
	res = strings.Replace(res, "\\", "\\2", -1)
	return res
}
func unescapeString(data string) string {
	res := strings.Replace(data, "/1", "/", -1)
	res = strings.Replace(res, "\\2", "\\", -1)
	return res
}

func addUserContext(ctx *msim_context) {
	users_context = append(users_context, ctx)
}

func removeUserContext(s []*msim_context, i int) []*msim_context {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func identifyProtocolVersion(clientver string) string {
	ver, _ := strconv.Atoi(clientver)

	if ver <= 253 {
		return "MSIMv1"
	} else if ver > 253 && ver <= 366 {
		return "MSIMv2"
	} else if ver > 366 && ver <= 595 {
		return "MSIMv3"
	} else if ver > 595 && ver <= 673 {
		return "MSIMv4"
	} else if ver > 673 && ver <= 697 {
		return "MSIMv5"
	} else if ver > 697 && ver <= 812 {
		return "MSIMv6"
	}

	return ""
}
