package worker

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	gs "git.dragonheim.net/dragonheim/gagent/src/gstructs"
	picol "git.dragonheim.net/dragonheim/gagent/src/picol"

	zmq "github.com/pebbe/zmq4"
)

func pop(msg []string) (head, tail []string) {
	if msg[1] == "" {
		head = msg[:2]
		tail = msg[2:]
	} else {
		head = msg[:1]
		tail = msg[1:]
	}
	return
}

//  Each worker task works on one request at a time and sends a random number
//  of replies back, with random delays between replies:
func agentHandler(workerNum int) {
	interp := picol.InitInterp()
	interp.RegisterCoreCommands()

	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	worker.Connect("inproc://backend")

	for {
		//  The DEALER socket gives us the reply envelope and message
		msg, _ := worker.RecvMessage(0)
		identity, content := pop(msg)

		//  Send 0..4 replies back
		replies := rand.Intn(5)
		for reply := 0; reply < replies; reply++ {
			//  Sleep for some fraction of a second
			time.Sleep(time.Duration(rand.Intn(1000)+1) * time.Millisecond)

			log.Printf(fmt.Sprintf("Worker %d: %s\n", workerNum, identity))
			worker.SendMessage(identity, content)
		}
	}
}

// Main is the initiation function for a Worker
func Main(config gs.GagentConfig) {
	//  Frontend socket talks to clients over TCP
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	log.Printf("Starting worker\n")

	defer frontend.Close()
	log.Printf("Attempting to connect to: %s(%s)\n", config.Routers[0].RouterName, config.Routers[0].RouterAddr)
	connectString := fmt.Sprintf("tcp://%s", config.Routers[0].RouterAddr)
	frontend.Bind(connectString)

	//  Backend socket talks to workers over inproc
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")

	//  Launch pool of agent handlers
	for i := 0; i < 5; i++ {
		go agentHandler(i)
	}

	//  Connect backend to frontend via a proxy
	// err := zmq.Proxy(frontend, backend, nil)
	// log.Fatalln("Proxy interrupted:", err)
}

