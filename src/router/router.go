package router

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

// Main is the initiation function for a Router
func Main(config gs.GagentConfig) {
	/*
	 * This is our router task.
	 *
	 * It uses the multi-threaded server model to deal requests out to a
	 * pool of workers and route replies back to clients. One worker can
	 * handle one request at a time but one client can talk to multiple
	 * workers at once.
	 *
	 * Frontend socket talks to clients over TCP
	 */
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()

	frontend.Bind(fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.ListenPort))

	//  Backend socket talks to workers over inproc
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")

	//  Launch pool of worker threads, precise number is not critical
	for i := 0; i < 5; i++ {
		go agentRouter(i)
	}

	//  Connect backend to frontend via a proxy
	err := zmq.Proxy(frontend, backend, nil)
	log.Fatalln("Proxy interrupted:", err)
}

//  Each worker task works on one request at a time and sends a random number
//  of replies back, with random delays between replies:

func agentRouter(workerNum int) {
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

			fmt.Println(fmt.Sprintf("Worker %d: %s", workerNum, identity))
			fmt.Println(fmt.Sprintf("Worker %d: %s", workerNum, content))
			worker.SendMessage(identity, content)
		}
	}
}
