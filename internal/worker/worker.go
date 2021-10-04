package worker

import (
	"fmt"
	"log"
	"sync"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	// picol "git.dragonheim.net/dragonheim/gagent/src/picol"
	zmq "github.com/pebbe/zmq4"
)

// Main is the initiation function for a Worker
func Main(wg *sync.WaitGroup, config gs.GagentConfig, rid int) {
	defer wg.Done()
	log.Printf("[INFO] Starting worker\n")

	// Generate connect string for this router.
	var rport = int64(config.WorkerPort)
	if config.Routers[rid].WorkerPort != 0 {
		rport = config.Routers[rid].WorkerPort
	}
	connectString := fmt.Sprintf("tcp://%s:%d", config.Routers[rid].RouterAddr, rport)

	subscriber, _ := zmq.NewSocket(zmq.REP)
	defer subscriber.Close()

	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	subscriber.Connect(connectString)

	msg, err := subscriber.Recv(0)
	if err != nil {
		log.Printf("[DEBUG] Received error: %v", err)
	}
	log.Printf("[DEBUG] Received message: %v", msg[0])
}
