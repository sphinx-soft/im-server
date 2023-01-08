package aim

import (
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

			// currently sequence and versionHandshaken are here, but they
			// might move later
			sequence := uint16(rand.Intn(0xFFFF))
			versionHandshaken := false

			versionFlap := FLAPPacket{
				Frame:    FrameSignOn,
				Sequence: sequence,
				Data:     []byte{0x00, 0x00, 0x00, 0x01},
			}

			client.Connection.BinaryWriteTraffic(FLAPSerialize(versionFlap))

			for {
				combined, err := client.Connection.BinaryReadTraffic()

				if err != nil {
					break
				}

				packets, hasErr := FLAPDeserialize(combined)
				if hasErr {
					break
				}

				for _, packet := range packets {
					if packet.Frame != FrameSignOn && !versionHandshaken {
						return
					}

					switch packet.Frame {
					case FrameSignOn:
						versionHandshaken = true

					case FrameData:
						// SNAC data frame

					case FrameError:
						// Error frame

					case FrameSignOff:
						// Sign off frame
					}
				}
			}
		}()
	}

}
