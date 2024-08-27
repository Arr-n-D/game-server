package internal

import "github.com/arr-n-d/gns"

var serverManager *ServerManager

type ServerManager struct {
	bQuit     bool
	pollGroup gns.PollGroup
}

func InitServerManager() {}

func initGameNetworkingSockets() {
	err := gns.Init(nil)

    if err != {
        
    }

}
