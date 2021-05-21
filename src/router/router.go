package router

import (
	"fmt"
	"log"

	gs "git.dragonheim.net/dragonheim/gagent/src/gstructs"

	zmq "github.com/pebbe/zmq4"
)

const (
	WORKER_READY = "\001" //  Signals worker is ready
)

// Main is the initiation function for a Router
func Main(config gs.GagentConfig) {
	log.Printf("[INFO] Starting router\n")

	client_sock, _ := zmq.NewSocket(zmq.ROUTER)
	defer client_sock.Close()

	worker_sock, _ := zmq.NewSocket(zmq.DEALER)
	defer worker_sock.Close()

	client_sock.Bind(fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.ClientPort))
	worker_sock.Bind(fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.WorkerPort))

	workers := make([]string, 0)

	poller1 := zmq.NewPoller()
	poller1.Add(worker_sock, zmq.POLLIN)

	poller2 := zmq.NewPoller()
	poller2.Add(worker_sock, zmq.POLLIN)
	poller2.Add(client_sock, zmq.POLLIN)

LOOP:
	for {
		//  Poll frontend only if we have available workers
		var sockets []zmq.Polled
		var err error
		if len(workers) > 0 {
			sockets, err = poller2.Poll(-1)
		} else {
			sockets, err = poller1.Poll(-1)
		}
		if err != nil {
			break //  Interrupted
		}
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case worker_sock: //  Handle worker activity on backend
				//  Use worker identity for load-balancing
				msg, err := s.RecvMessage(0)
				if err != nil {
					break LOOP //  Interrupted
				}
				var identity string
				identity, msg = unwrap(msg)
				log.Printf("[DEBUG] Worker message received: %s", msg)
				workers = append(workers, identity)

				//  Forward message to client if it's not a READY
				if msg[0] != WORKER_READY {
					client_sock.SendMessage(msg)
				}

			case client_sock:
				//  Get client request, route to first available worker
				msg, err := s.RecvMessage(0)
				log.Printf("[DEBUG] Client message received: %s", msg)
				if err == nil {
					worker_sock.SendMessage(workers[0], "", msg)
					workers = workers[1:]
				}
			}
		}
	}
}

func unwrap(msg []string) (head string, tail []string) {
	head = msg[0]
	if len(msg) > 1 && msg[1] == "" {
		tail = msg[2:]
	} else {
		tail = msg[1:]
	}
	return
}
