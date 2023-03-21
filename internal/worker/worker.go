package worker

import (
	log "log"
	strconv "strconv"
	sync "sync"

	gstructs "github.com/dragonheim/gagent/internal/gstructs"

	/*
	 * picol "github.com/dragonheim/gagent/pkg/picol"
	 */

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
 * The "worker" processes the agent code. The worker nodes do not know
 * anything about the network structure. Instead they know only to which
 * router(s) they are connected. The worker will execute the agent code and
 * pass the agent and it's results to a router.
 * Main is the entrypoint for the worker process
 */
func Main(wg *sync.WaitGroup, config gstructs.GagentConfig) {
	log.Printf("[INFO] Starting worker\n")
	defer wg.Done()

	for key := range config.Routers {
		rport := config.WorkerPort
		if config.Routers[key].WorkerPort != 0 {
			rport = config.Routers[key].WorkerPort
		}

		/*
		 * Generate connect string for this router.
		 */
		connectString := "tcp://" + config.Routers[key].RouterAddr + ":" + strconv.Itoa(rport)

		wg.Add(1)
		go getAgent(wg, config.UUID, connectString)
	}
	/*
	 * workerListener := "tcp://" + config.ListenAddr + ":" + strconv.Itoa(config.WorkerPort)
	 */

}

func getAgent(wg *sync.WaitGroup, uuid string, connectString string) {
	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	defer wg.Done()

	subscriber, _ := zmq.NewSocket(zmq.REP)
	defer subscriber.Close()

	_ = subscriber.Connect(connectString)

	msg, err := subscriber.Recv(0)
	if err != nil {
		log.Printf("[DEBUG] Received error: %v", err)
	}
	log.Printf("[DEBUG] Received message: %v", msg[0])
}
