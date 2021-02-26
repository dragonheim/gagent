package client

import (
	"fmt"
	"sync"
	"time"

	gs "git.dragonheim.net/dragonheim/gagent/src/gstructs"

	zmq "github.com/pebbe/zmq4"
)

// Main is the initiation function for a Client
func Main(config gs.GagentConfig, agent string) {
	var mu sync.Mutex

	fmt.Printf("Did we make it this far?\n")
	fmt.Printf("--|%#v|--\n", agent)

	connectString := fmt.Sprintf("tcp://%s", config.Routers[0].RouterAddr)

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
			fmt.Println(msg[0], config.UUID)
		}
		mu.Unlock()
	}
}
