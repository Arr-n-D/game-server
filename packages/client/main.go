package main

import (
	"log"
	"net"

	"github.com/arr-n-d/gns"
)

func main() {
	gns.Init(nil)
	gns.SetDebugOutputFunction(gns.DebugOutputTypeEverything, func(typ gns.DebugOutputType, msg string) {
		log.Print("[DEBUG]", typ, msg)
	})
	defer gns.Kill()
	str := "Hello, world!"

	c, err := gns.Dial(&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 27015}, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close()

	if _, err := c.Write([]byte(str)); err != nil {
		log.Print(err)
		return
	}

}
