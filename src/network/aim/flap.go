package aim

import (
	"encoding/binary"
	"errors"
)

const (
	FrameSignOn  = 0x01
	FrameData    = 0x02
	FrameError   = 0x03
	FrameSignOff = 0x04
)

type FLAPPacket struct {
	Frame    byte
	Sequence uint16
	Data     []byte
}

func FLAPSerialize(flap []byte) (outPackets []FLAPPacket, err error) {
	packets := []FLAPPacket{}

	i := 0

	for i < len(flap) {
		if flap[0] != 0x2A {
			return packets, errors.New("invalid marker")
		} else if len(flap)-i < 6 {
			return packets, errors.New("incorrect length")
		}

		length := binary.BigEndian.Uint16(flap[i+4 : i+6])
		packet := FLAPPacket{
			Frame:    flap[i+1],
			Sequence: binary.BigEndian.Uint16(flap[i+2 : i+4]),
		}

		if int(length) > len(flap)-i {
			return packets, errors.New("incorrect length")
		} else if packet.Frame != FrameSignOn && packet.Frame != FrameData && packet.Frame != FrameError && packet.Frame != FrameSignOff {
			return packets, errors.New("incorrect frame")
		}

		packet.Data = make([]byte, length)
		copy(packet.Data, flap[i+6:i+6+int(length)])

		i += 6 + int(length)
		packets = append(packets, packet)
	}

	return packets, nil
}

func FLAPDeserialize(packet FLAPPacket) []byte {
	flap := make([]byte, 6+len(packet.Data))

	flap[0] = 0x2A
	flap[1] = packet.Frame

	// I'm so glad I can use slicing here
	binary.BigEndian.PutUint16(flap[2:4], packet.Sequence)
	binary.BigEndian.PutUint16(flap[4:6], uint16(len(packet.Data)))

	copy(flap[6:], packet.Data)
	return flap
}
