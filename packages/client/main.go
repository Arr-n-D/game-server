package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/arr-n-d/gns"
)

func StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnecting:
		fmt.Println("Connecting")

	case gns.ConnectionStateConnected:
		fmt.Println("Connected")

		baseMessage := "Hello, world! Sequence: "

		for i := 1; i <= 500; i++ {
			// Create the message with a sequence number
			message := baseMessage + strconv.Itoa(i)

			// Send the message
			_, res := info.Conn().SendMessage([]byte(message), gns.SendReliable)

			if res != gns.ResultOK {
				fmt.Println("Issue fault")
			}
		}

		fmt.Println("Sent data")

	case gns.ConnectionStateProblemDetectedLocally:
		info.Conn().Close(gns.ConnectionEndAppExceptionGeneric, "", false)
		os.Exit(1)
		fmt.Println("Problem detected locally")

	}
}

func main() {
	// var wg sync.WaitGroup
	gns.Init(nil)
	gns.SetDebugOutputFunction(gns.DebugOutputTypeEverything, func(typ gns.DebugOutputType, msg string) {
		log.Print("[DEBUG]", typ, msg)
	})
	defer gns.Kill()

	addrr := &net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: 27015,
	}

	gns.SetGlobalCallbackStatusChanged(StatusCallBackChanged)
	gnsadr := gns.NewIPAddr(addrr)
	c := gns.ConnectByIPAddress(gnsadr, nil)

	if c == gns.InvalidConnection {
		log.Fatal("Invalid connection")
	}

	for {
		gns.RunCallbacks()
	}
	// c, err := gns.Dial(&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 27015}, nil)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// defer c.Close(gns.ConnectionEndAppExceptionMin, "False", false)

}
