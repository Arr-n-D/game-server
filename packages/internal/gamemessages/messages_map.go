package gamemessages

var MESSAGE_TYPE_TO_TYPE_STRUCT = map[byte]interface{}{
	1: Sequence{},
}

var MESSAGE_TYPE_TO_GAME_FUNC map[byte]func(interface{})
