package server

import "internal/gamemessages"

func (s *Server) InitializeGameFuncsMap() {
	gamemessages.MESSAGE_TYPE_TO_GAME_FUNC = map[byte]func(interface{}){
		1: func(interface{}) {
			s.Test()
		},
	}
}
func (s *Server) Test() {

}
