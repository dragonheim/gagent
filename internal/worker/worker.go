package worker

import (
	fmt "fmt"
	log "log"
	sync "sync"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	// picol "git.dragonheim.net/dragonheim/gagent/src/picol"

	prometheus "github.com/prometheus/client_golang/prometheus"
	promauto "github.com/prometheus/client_golang/prometheus/promauto"

	zmq "github.com/pebbe/zmq4"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "agent_requests_collected",
	})
)

/*
 The "worker" processes the agent code. The worker nodes do not know
 anything about the network structure. Instead they know only to which
 router(s) they are connected. The worker will execute the agent code and
 pass the agent and it's results to a router.
*/
func Main(wg *sync.WaitGroup, config gs.GagentConfig) {
	defer wg.Done()
	log.Printf("[INFO] Starting worker\n")

	for key := range config.Routers {
		rport := config.WorkerPort
		if config.Routers[key].WorkerPort != 0 {
			rport = config.Routers[key].WorkerPort
		}

		// Generate connect string for this router.
		connectString := fmt.Sprintf("tcp://%s:%d", config.Routers[key].RouterAddr, rport)

		wg.Add(1)
		go getAgent(wg, config.UUID, connectString)
	}
	// workerListener := fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.WorkerPort)

}

func getAgent(wg *sync.WaitGroup, uuid string, connectString string) {
	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	defer wg.Done()

	subscriber, _ := zmq.NewSocket(zmq.REP)
	defer subscriber.Close()

	subscriber.Connect(connectString)

	msg, err := subscriber.Recv(0)
	if err != nil {
		log.Printf("[DEBUG] Received error: %v", err)
	}
	log.Printf("[DEBUG] Received message: %v", msg[0])
}
