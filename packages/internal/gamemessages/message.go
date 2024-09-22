package gamemessages

import "github.com/arr-n-d/gns"

type GameMessage struct {
	MessageType    byte
	MessageContent []byte

	MsgNumber int64
	UserData  int64

	Connection gns.Connection
	ReceivedAt gns.Timestamp
	Flags      gns.SendFlags
}
