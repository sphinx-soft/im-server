package msnp

import (
	"fmt"
	"phantom/global"
	"phantom/util"
	"strconv"
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
	if strings.HasPrefix(data, "USR") {
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
	} else { //[Debug] [TCP -> ReadTraffic] Reading Data: ANS 1 test2@hotmail.com 1843e8b2e6b 31847
		mail := strings.Replace(findValueFromData("ANS", data, 1), "@hotmail.com", util.GetMailDomain(), -1)
		authenticate := findValueFromData("ANS", data, 2)
		sessionid, _ := strconv.Atoi(findValueFromData("ANS", data, 3))

		if authenticate != ctx.authentication {
			util.WriteTraffic(ctx.connection, msnp_new_command_noargs(data, "911"))
			return false
		}

		for i := 0; i < len(msn_switchboard_sessions); i++ {
			if msn_switchboard_sessions[i] != nil {
				if msn_switchboard_sessions[i].sessionid == sessionid {
					msn_switchboard_sessions[i].clients = append(msn_switchboard_sessions[i].clients, ctx)
				}

				if msn_switchboard_sessions[i].clients[i].email != mail {
					util.WriteTraffic(ctx.connection, msnp_new_command("IRO", data, fmt.Sprintf("")))
				}
			}
		}
	}

	return false
}

// todo
func handleClientSwitchboardPacketSendSwitchboardInvite(ctx *msnp_switchboard_context, data string) {

	mail := strings.Replace(findValueFromData("CAL", data, 1), "@hotmail.com", util.GetMailDomain(), -1)
	var cl *msnp_context
	var cx *global.Client

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i] != nil {
			if global.Clients[i].Account.Email == mail {
				cx = global.Clients[i]
			}
		}
	}

	for i := 0; i < len(msn_context_list); i++ {
		if msn_context_list[i] != nil {
			if msn_context_list[i].email == cx.Account.Email {
				cl = msn_context_list[i]
			}
		}
	}

	if cl.status == "HDN" {
		util.WriteTraffic(ctx.connection, msnp_new_command_noargs(data, "217"))
		return
	}

	ctx.sessionid = generateContextKey() // generate random int
	util.WriteTraffic(ctx.connection, msnp_new_command(data, "CAL", fmt.Sprintf("RINGING %d", ctx.sessionid)))

	date := time.Now().UTC().UnixMilli()
	sbctx := msnp_switchboard_context{
		sessionid:      ctx.sessionid,
		username:       cx.Account.Screenname,
		email:          cx.Account.Email,
		authentication: strconv.FormatInt(date, 16),
	}
	addSwitchboardContext(&sbctx)

	for i := 0; i < len(msn_switchboard_sessions); i++ {
		if msn_switchboard_sessions[i] != nil {
			if msn_switchboard_sessions[i].sessionid != ctx.sessionid {
				sbsess := msnp_switchboard_session{
					sessionid: ctx.sessionid,
				}

				sbsess.clients = append(sbsess.clients, ctx)
				addSwitchboardSession(&sbsess)
				break
			}
		}
	}

	err := util.WriteTraffic(cx.Connection, fmt.Sprintf("RNG %d %s:1865 CKI %s %s %s\r\n", ctx.sessionid, util.GetRootUrl(), strconv.FormatInt(date, 16), ctx.email, ctx.username))
	if err != nil {
		util.WriteTraffic(ctx.connection, msnp_new_command_noargs(data, "217"))
		return
	}
}
