package msnp

import (
	"fmt"
	"phantom/global"
	"phantom/util"
	"strings"
	"time"
)

func handleClientIncomingSwitchboardPackets(ctx *msnp_switchboard_context, data string) {

	switch {
	case strings.HasPrefix(data, "CAL"):
		handleClientSwitchboardPacketSendSwitchboardInvite(ctx, data)
	}
}

func handleClientSwitchboardPacketAuthentication(ctx *msnp_switchboard_context, data string) bool {
	mail := strings.Replace(findValueFromData("USR", data, 1), "@hotmail.com", util.GetMailDomain(), -1)
	auth := findValueFromData("USR", data, 2)
	acc, _ := global.GetUserDataFromEmail(mail)

	util.Debug("MSNP -> handleClientSwitchboardPacketAuthentication", "auth test1: %s", auth)
	util.Debug("MSNP -> handleClientSwitchboardPacketAuthentication", "auth test2: %s", ctx.authentication)

	if ctx.authentication == auth {
		ctx.email = mail
		ctx.username = acc.Screenname

		util.WriteTraffic(ctx.connection, msnp_new_command(data, "USR", fmt.Sprintf("OK %s %s", mail, ctx.username)))
		return true
	} else {
		util.WriteTraffic(ctx.connection, msnp_new_command_noargs(data, "911"))
		return false
	}
}

// todo
func handleClientSwitchboardPacketSendSwitchboardInvite(ctx *msnp_switchboard_context, data string) {

	mail := findValueFromData("CAL", data, 1)
	var cl *msnp_context
	var cx *global.Client

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.Email == mail {
			cx = global.Clients[i]
		}
	}

	for i := 0; i < len(msn_context_list); i++ {
		if msn_context_list[i].email == cx.Account.Email {
			cl = msn_context_list[i]
		}
	}

	if cl.status == "HDN" {
		util.WriteTraffic(ctx.connection, msnp_new_command_noargs(data, "217"))
		return
	}

	ctx.sessionid = generateContextKey() // generate random int
	util.WriteTraffic(ctx.connection, msnp_new_command(data, "CAL", fmt.Sprintf("RINGING %d", ctx.sessionid)))

	date := time.Now().UTC().UnixMilli()
	util.WriteTraffic(cx.Connection, fmt.Sprintf("RNG %d %s:1865 CKI %d %s %s\r\n", ctx.sessionid, util.GetRootUrl(), date, ctx.email, ctx.username))
}
