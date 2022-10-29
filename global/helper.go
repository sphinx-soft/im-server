package global

import (
	"fmt"
	"phantom/util"
	"strings"
)

func AddClient(client *Client) {
	Clients = append(Clients, client)
}

func RemoveClient(s []*Client, i int) []*Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func GetClient(username string) *Client {
	for i := 0; i < len(Clients); i++ {
		if Clients[i].Account.Email == username {
			return Clients[i]
		}
	}

	return nil
}

func GetUserDataFromEmail(email string) (Account, bool) {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE email= ?", email)

	if err != nil {
		util.Error("Failed to get email userdata: %s", err.Error())
		return acc, false
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()
	acc.Username = strings.Replace(email, "@phantom-im.xyz", "", -1)

	return acc, true
}

func GetUserDataFromUsername(username string) (Account, bool) {

	var acc Account

	user := fmt.Sprintf("%s@phantom-im.xyz", username)

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE email= ?", user)

	if err != nil {
		util.Error("Failed to get username userdata: %s", err.Error())
		return acc, false
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()
	acc.Username = username

	return acc, true
}

func GetUserDataFromIcqNumber(uin int) (Account, bool) {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE uin= ?", uin)

	if err != nil {
		util.Error("Failed to get icq number userdata: %s", err.Error())
		return acc, false
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()
	acc.Username = strings.Replace(acc.Email, "@phantom-im.xyz", "", -1)

	return acc, true
}

func GetUserDataFromUserId(uid int) (Account, bool) {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE id= ?", uid)

	if err != nil {
		util.Error("Failed to get icq number userdata: %s", err.Error())
		return acc, false
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()
	acc.Username = strings.Replace(acc.Email, "@phantom-im.xyz", "", -1)

	return acc, true
}

func GetUploadDataFromUserId(uid int) (Upload, bool) {

	var upl Upload

	row, err := util.GetDatabaseHandle().Query("SELECT * from upload WHERE id= ?", uid)

	if err != nil {
		util.Error("Failed to get data userdata: %s", err.Error())
		return upl, false
	}

	row.Next()
	row.Scan(&upl.UserId, &upl.Avatar)
	row.Close()

	return upl, true
}
