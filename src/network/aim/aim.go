package aim

import (
	"chimera/network"
	"chimera/utility/logging"
	"chimera/utility/tcp"
	"encoding/hex"
)

func LogonAIM() {

	// The BUCP server listens on 5190 and handles authentication, then the client is transported to the BOS server which listens on 5191
	// Right now the BUCP server is the only one which is listening
	tcpServer := tcp.CreateListener(5190)

	for {
		err := tcpServer.AcceptClient()

		go func() {
			if err != nil {
				logging.Error("AIM/BUCP", "Failed to accept client! (%s)", err.Error())
				return
			}

			logging.Info("AIM/BUCP", "Client connected! (IP: %s)", tcpServer.GetRemoteAddress())

			client := network.Client{
				Connection: tcpServer,
			}

			// a FLAP header consists of the following:

			// 2A - the marker
			// 01 - the frame, in this case this frame means to initialize the connection and the data contains the server FLAP version
			// which is a DWORD and always 1
			// 0001 - the sequence, which is initialized to a random value between 0x0000 and 0xFFFF and it wraps to 0x0000 if it's 0xFFFF,
			// otherwise increments for each packet sent to the client.
			// 0004 - the length. i feel that this is self-explanatory.

			// currently this is simply a example packet in hexadecimal which is decoded and sent to the client, but eventually I'll add a
			// FLAP packet builder soon

			data, err := hex.DecodeString("2A010001000400000001")
			if err != nil {
				return
			}

			client.Connection.BinaryWriteTraffic(data)

			for {
				_, err := client.Connection.BinaryReadTraffic()

				if err != nil {
					break
				}

				// TODO: parse the packets, add structures, etc.
			}
		}()
	}

}
