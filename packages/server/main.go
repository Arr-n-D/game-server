package main

import (
	"fmt"
	config "internal/configuration"
	"net"

	"log"

	"github.com/arr-n-d/gns"
)

var g_bQuit bool = false
var pollGroup gns.PollGroup

func StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnecting:
		fmt.Println("Connecting")
		conn := info.Conn()
		if conn.Accept() != gns.ResultOK {
			log.Fatalln("Failed to accept client")
		}

		if !conn.SetPollGroup(pollGroup) {
			log.Fatalln("Failed to set poll group")
		}

	}

}

func main() {
	config.InitSentry()

	gns.Init(nil)
	gns.SetDebugOutputFunction(gns.DebugOutputTypeEverything, func(typ gns.DebugOutputType, msg string) {
		log.Print("[DEBUG] ", typ, msg)
	})

	gns.SetGlobalCallbackStatusChanged(StatusCallBackChanged)
	defer gns.Kill()

	l, err := gns.Listen(&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 27015}, nil)
	fmt.Println(l.Addr())

	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		log.Fatal("Invalid poll group")
	}

	pollGroup = poll

	fmt.Println("Server running on port 27015")

	for ok := true; ok; ok = !g_bQuit {
		gns.RunCallbacks()
	}

}
