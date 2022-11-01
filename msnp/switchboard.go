package msnp

import (
	"fmt"
	"phantom/global"
	"phantom/util"
	"strings"
)

func handleClientIncomingSwitchboardPackets(ctx *msnp_switchboard_context, data string) {

	switch {

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
