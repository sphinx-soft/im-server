package global

import (
	"fmt"
	"phantom/util"
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

func GetUserDataFromEmail(email string) Account {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE email= ?", email)

	if err != nil {
		util.Error("Failed to get email userdata: %s", err.Error())
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()

	return acc
}

func GetUserDataFromUsername(username string) Account {

	var acc Account

	user := fmt.Sprintf("%s@phantom-im.xyz", username)

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE email= ?", user)

	if err != nil {
		util.Error("Failed to get username userdata: %s", err.Error())
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()

	return acc
}

func GetUserDataFromIcqNumber(uin int) Account {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE uin= ?", uin)

	if err != nil {
		util.Error("Failed to get icq number userdata: %s", err.Error())
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()

	return acc
}

func GetUserDataFromUserId(uid int) Account {

	var acc Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from accounts WHERE id= ?", uid)

	if err != nil {
		util.Error("Failed to get icq number userdata: %s", err.Error())
	}

	row.Next()
	row.Scan(&acc.UserId, &acc.Email, &acc.Password, &acc.Screenname, &acc.ICQNumber)
	row.Close()

	return acc
}
