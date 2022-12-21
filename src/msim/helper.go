package msim

import (
	"math/rand"
	"phantom/global"
	"phantom/util"
	"strconv"
	"strings"
)

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
		util.Log(util.INFO, "MySpace -> getMySpaceDataByEmail", err.Error())
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
		util.Log(util.INFO, "MySpace -> getMySpaceDataByUserId", err.Error())
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
	} else if ver > 253 && ver < 366 {
		return "MSIMv2"
	} else if ver > 366 && ver < 404 {
		return "MSIMv3"
	} else if ver > 404 && ver < 594 {
		return "MSIMv4"
	} else if ver > 593 && ver < 673 {
		return "MSIMv5"
	} else if ver > 673 && ver < 697 {
		return "MSIMv6"
	} else if ver > 697 && ver < 812 {
		return "MSIMv7"
	}

	return ""
}
