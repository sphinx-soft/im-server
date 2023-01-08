package aim

import (
	"encoding/binary"
)

type SNACMessage struct {
	Foodgroup uint16
	Subgroup  uint16
	Flags     uint16
	RequestID uint32
	Data      []byte
}

func SNACSerialize(snac []byte) SNACMessage {
	message := SNACMessage{
		Foodgroup: binary.BigEndian.Uint16(snac[0:2]),
		Subgroup:  binary.BigEndian.Uint16(snac[2:4]),
		Flags:     binary.BigEndian.Uint16(snac[4:6]),
		RequestID: binary.BigEndian.Uint32(snac[6:10]),
	}
	message.Data = make([]byte, len(snac)-10)
	copy(message.Data, snac[10:])

	return message
}

func SNACDeserialize(message SNACMessage) []byte {
	snac := make([]byte, len(message.Data)+10)

	binary.BigEndian.PutUint16(snac[0:2], message.Foodgroup)
	binary.BigEndian.PutUint16(snac[2:4], message.Subgroup)

	binary.BigEndian.PutUint16(snac[4:6], message.Flags)
	binary.BigEndian.PutUint32(snac[6:10], message.RequestID)

	copy(snac[10:], message.Data)
	return snac
}
