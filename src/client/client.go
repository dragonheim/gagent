package client

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	gs "git.dragonheim.net/dragonheim/gagent/src/gstructs"

	zmq "github.com/pebbe/zmq4"
)

// Main is the initiation function for a Client
func Main(config gs.GagentConfig, routerID int, agent string) {
	var mu sync.Mutex
	var rport = int(config.ListenPort)
	if config.Routers[routerID].RouterPort != "" {
		rport, _ = strconv.Atoi(config.Routers[routerID].RouterPort)
	}

	log.Printf("--|%#v|--\n", agent)

	connectString := fmt.Sprintf("tcp://%s:%d",
		config.Routers[routerID].RouterAddr,
		rport)
	log.Printf("Attempting to connect to %s\n", connectString)

	sock, _ := zmq.NewSocket(zmq.DEALER)
	defer sock.Close()

	sock.SetIdentity(config.UUID)
	sock.Connect(connectString)

	go func() {
		mu.Lock()
		sock.SendMessage(agent)
		mu.Unlock()
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		msg, err := sock.RecvMessage(zmq.DONTWAIT)
		if err == nil {
			log.Println(msg[0], config.UUID)
		}
		mu.Unlock()
	}
}
