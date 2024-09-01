package server

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"internal/configuration"

	"github.com/arr-n-d/gns"
)

var ServerInstance *Server

const (
	tickRate     = 24
	tickDuration = time.Second / tickRate
)

func InitServer() error {
	if err := configuration.InitSentry(); err != nil {
		return fmt.Errorf("failed to initialize sentry. %w", err)
	}
	if err := initGameNetworkingSockets(); err != nil {
		return fmt.Errorf("failed to initialize game networking sockets. %w", err)
	}
	if err := initServer(); err != nil {
		return fmt.Errorf("failed to initServer. %w", err)
	}

	gns.SetGlobalCallbackStatusChanged(ServerInstance.StatusCallBackChanged)

	ServerInstance.init()
	return nil
}

func initGameNetworkingSockets() error {
	err := gns.Init(nil)
	if err != nil {
		return fmt.Errorf("gns.Init failed: %w", err)
	}
	return nil
}

func setDebugOutputFunction(detailLevel gns.DebugOutputType) {
	gns.SetDebugOutputFunction(detailLevel, func(typ gns.DebugOutputType, msg string) {
		slog.With("type", typ, "msg", msg).Debug("[DEBUG]")
		log.Print("[DEBUG] ", typ, msg)
	})
}

func initServer() error {
	conf := configuration.GetConfiguration()

	if conf.Env == "" {
		return errors.New("no configuration found")
	}

	if conf.Env == configuration.DevEnv {
		setDebugOutputFunction(gns.DebugOutputTypeEverything)
	}

	l, err := gns.Listen(&net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: int(conf.GameServerPort),
	},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d. %w", conf.GameServerPort, err)
	}

	poll := gns.NewPollGroup()
	if poll == gns.InvalidPollGroup {
		return errors.New("failed to create poll group")
	}

	ServerInstance = &Server{
		PollGroup:              poll,
		listener:               l,
		Quit:                   false,
		ReceiveMessagesChannel: make(chan []byte, 200),
		SendMessagesChannel:    make(chan []byte, 200),
	}
	return nil
}
