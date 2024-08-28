package internal

import (
	"log"

	"github.com/arr-n-d/gns"
)

var serverManager *ServerManager

type ServerManager struct {
	bQuit     bool
	pollGroup gns.PollGroup
	logger    *log.Logger
}

func InitServerManager() {}

func initGameNetworkingSockets() {
	err := gns.Init(nil)

	// if err != {

	// }

}
