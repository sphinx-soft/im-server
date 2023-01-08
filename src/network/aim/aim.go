package aim

import (
	"bytes"
	"chimera/network"
	"chimera/utility/logging"
	"chimera/utility/tcp"
	"math/rand"
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

			flapVersion := []byte{0x00, 0x00, 0x00, 0x01}

			// currently sequence are here, but they
			// might move later
			sequence := uint16(rand.Intn(0xFFFF))

			versionFlap := FLAPPacket{
				Frame:    FrameSignOn,
				Sequence: sequence,
				Data:     flapVersion,
			}

			client.Connection.BinaryWriteTraffic(FLAPDeserialize(versionFlap))

			for {
				combined, err := client.Connection.BinaryReadTraffic()
				if err != nil {
					break
				}

				packets, err := FLAPSerialize(combined)
				if err != nil {
					break
				}

				for _, packet := range packets {
					switch packet.Frame {
					case FrameSignOn:
						if bytes.Equal(packet.Data, flapVersion) {
							continue
						}
						// this is old FLAP authentication, TODO: implement this

					case FrameData:
						// this is snac data
						
					}
				}
			}

			client.Connection.CloseConnection()
		}()
	}

}
