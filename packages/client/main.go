package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/arr-n-d/gns"
)

const (
	NumPlayers        = 100
	MessagesPerPlayer = 150 // 5 seconds worth of messages at 30 per second
	TotalDuration     = 5 * time.Second
)

func StatusCallBackChanged(info *gns.StatusChangedCallbackInfo) {
	switch state := info.Info().State(); state {
	case gns.ConnectionStateConnecting:
		fmt.Println("Connecting")

	case gns.ConnectionStateConnected:
		fmt.Println("Connected")
		runSimulation(info.Conn())

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

func runSimulation(connection gns.Connection) {
	var wg sync.WaitGroup

	fmt.Printf("Starting simulation with %d players for %v\n", NumPlayers, TotalDuration)

	startTime := time.Now()

	for i := 1; i <= NumPlayers; i++ {
		wg.Add(1)
		go simulatePlayer(i, connection, &wg)
	}

	wg.Wait()

	duration := time.Since(startTime)
	totalMessages := NumPlayers * MessagesPerPlayer
	messagesPerSecond := float64(totalMessages) / duration.Seconds()

	fmt.Printf("Simulation completed in %v\n", duration)
	fmt.Printf("Total messages sent: %d\n", totalMessages)
	fmt.Printf("Average messages per second: %.2f\n", messagesPerSecond)
}

func simulatePlayer(playerID int, connection gns.Connection, wg *sync.WaitGroup) {
	defer wg.Done()

	baseMessage := fmt.Sprintf("Player %d, Sequence: ", playerID)

	for i := 1; i <= MessagesPerPlayer; i++ {
		message := baseMessage + strconv.Itoa(i)

		_, res := connection.SendMessage([]byte(message), gns.SendReliable)
		if res != gns.ResultOK {
			fmt.Printf("Player %d: Issue fault on message %d\n", playerID, i)
		}

		// Sleep to simulate 30 messages per second
		time.Sleep(TotalDuration / MessagesPerPlayer)
	}

	fmt.Printf("Player %d finished sending %d messages\n", playerID, MessagesPerPlayer)
}
