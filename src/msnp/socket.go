package msnp

import (
	"bytes"
	"phantom/global"
	"phantom/util"
	"strings"
)

func HandleNotification() {
	tcpServer := util.CreateListener(1864)

	for {
		tcpClient, err := tcpServer.Accept()

		go func() {
			if err != nil {
				util.Error("MSNP -> HandleNotification", "Failed to accept Client! ", err.Error())
			} else {
				util.Debug("MSNP -> HandleNotification", "Accepted Client")
			}

			util.Log(util.INFO, "MSN Messenger", "Client awaiting authentication from %s", tcpClient.RemoteAddr().String())

			client := global.Client{
				Connection: tcpClient,
			}

			global.AddClient(&client)

			ctx := msnp_context{
				dispatched: true,
				ctxkey:     generateContextKey(),
			}

			addUserContext(&ctx)

			for {
				data, success := util.ReadTraffic(client.Connection)

				recv := strings.Split(string(data), "\r\n")

				for ix := 0; ix < len(recv); ix++ {
					recv[ix] = string(bytes.Trim([]byte(recv[ix]), "\x00"))
					if recv[ix] != "" {
						util.Debug("MSNP -> HandleNotification -> TCP", "Reading Split Data: %s", string(recv[ix]))
						handleClientIncomingPackets(&client, &ctx, recv[ix])
						//util.Debug("MSNP -> HandleNotification", "TCP dbg: %v", []byte(string(recv[ix])))
					}
				}

				if !success || handleClientLogoutRequest(string(data)) {
					break
				}
			}

			if client.Account.Email != "" {
				util.Log(util.INFO, "MSN Messenger", "Client Disconnected -> Email: %s", client.Account.Email)
			} else {
				util.Log(util.INFO, "MSN Messenger", "Client Disconnected -> Email: Unknown")
			}

			for i := 0; i < len(global.Clients); i++ {
				if global.Clients[i].Account.Email == client.Account.Email {
					util.Debug("MSNP -> HandleNotification", "Removing from clients from Clients List...")
					global.Clients = global.RemoveClient(global.Clients, i)
				}
			}

			for ix := 0; ix < len(msn_context_list); ix++ {
				if msn_context_list[ix].ctxkey == ctx.ctxkey {
					util.Debug("MSNP -> HandleNotification", "Removing from clients from Context List...")
					msn_context_list = removeUserContext(msn_context_list, ix)
				}
			}

			client.Connection.Close()
		}()
	}
}

func HandleDispatch(client *global.Client, firstread string) {
	util.Log(util.INFO, "MSN Messenger", "Client awaiting dispatch from %s", client.Connection.RemoteAddr().String())

	client.Client = "MSN Messenger"

	ctx := msnp_context{
		dispatched: false,
		ctxkey:     generateContextKey(),
	}

	addUserContext(&ctx)

	// Send first response command to MSN Client, Requesting INF Data
	if !handleClientProtocolVersionRequest(client, firstread) {
		util.Debug("MSNP -> HandleDispatch", "Unsupported MSNP Version requested, closing...")
		client.Connection.Close()
		return
	}

	for {
		data, success := util.ReadTraffic(client.Connection)

		recv := strings.Split(string(data), "\r\n")

		for ix := 0; ix < len(recv); ix++ {
			recv[ix] = string(bytes.Trim([]byte(recv[ix]), "\x00"))
			if recv[ix] != "" {
				util.Debug("MSNP -> HandleDispatch -> TCP", "Reading Split Data: %s", string(recv[ix]))
				handleClientIncomingPackets(client, &ctx, recv[ix])
				//util.Debug("MSNP -> HandleDispatch", "TCP dbg: %v", []byte(string(recv[ix])))
			}
		}

		if !success {
			break
		}
	}

	if client.Account.Email != "" {
		util.Log(util.INFO, "MSN Messenger", "Client Redirected -> Email: %s", client.Account.Email)
	} else {
		util.Log(util.INFO, "MSN Messenger", "Client Disconnected (DS) -> Email: Unknown")
	}

	for i := 0; i < len(global.Clients); i++ {
		if global.Clients[i].Account.Email == client.Account.Email {
			util.Debug("MSNP -> HandleDispatch", "Removing from clients from Clients List...")
			global.Clients = global.RemoveClient(global.Clients, i)
		}
	}

	for ix := 0; ix < len(msn_context_list); ix++ {
		if msn_context_list[ix].ctxkey == ctx.ctxkey {
			util.Debug("MSNP -> HandleDispatch", "Removing from clients from Context List...")
			msn_context_list = removeUserContext(msn_context_list, ix)
		}
	}

	client.Connection.Close()
}

func HandleSwitchboard() {
	tcpServer := util.CreateListener(1865)

	for {
		tcpClient, err := tcpServer.Accept()

		go func() {
			if err != nil {
				util.Error("MSNP -> HandleSwitchboard", "Failed to accept Client! ", err.Error())
			} else {
				util.Debug("MSNP -> HandleSwitchboard", "Accepted Client")
			}

			util.Log(util.INFO, "MSN Messenger", "Client joining switchboard from %s", tcpClient.RemoteAddr().String())

			data, _ := util.ReadTraffic(tcpClient)

			var ctx msnp_switchboard_context
			mail := strings.Replace(findValueFromData("USR", string(data), 1), "@hotmail.com", util.GetMailDomain(), -1)
			for i := 0; i < len(msn_switchboard_list); i++ {
				if msn_switchboard_list[i].email == mail {
					util.Debug("MSNP -> HandleSwitchboard", "Found Switchboard Context by Mail!")
					ctx = *msn_switchboard_list[i]
					ctx.connection = tcpClient
				}
			}

			if !handleClientSwitchboardPacketAuthentication(&ctx, string(data)) {
				util.Debug("MSNP -> HandleSwitchboard", "Failed to authenticate Switchboard session, closing...")
				tcpClient.Close()
				return
			}

			for {
				data, success := util.ReadTraffic(ctx.connection)

				recv := strings.Split(string(data), "\r\n")

				for ix := 0; ix < len(recv); ix++ {
					recv[ix] = string(bytes.Trim([]byte(recv[ix]), "\x00"))
					if recv[ix] != "" {
						util.Debug("MSNP -> HandleSwitchboard -> TCP", "Reading Split Data: %s", string(recv[ix]))
						handleClientIncomingSwitchboardPackets(&ctx, recv[ix])
						//util.Debug("MSNP -> HandleNotification", "TCP dbg: %v", []byte(string(recv[ix])))
					}
				}

				if !success || handleClientLogoutRequest(string(data)) {
					break
				}
			}

			if ctx.email != "" {
				util.Log(util.INFO, "MSN Messenger", "Client Left Switchboard -> Email: %s", ctx.email)
			} else {
				util.Log(util.INFO, "MSN Messenger", "Client Disconnected (SB) -> Email: Unknown")
			}

			for i := 0; i < len(msn_switchboard_list); i++ {
				if msn_switchboard_list[i].email == ctx.email {
					util.Debug("MSNP -> HandleSwitchboard", "Removing from clients from Context List...")
					msn_switchboard_list = removeSwitchboardContext(msn_switchboard_list, i)
				}
			}

			ctx.connection.Close()
		}()
	}
}
