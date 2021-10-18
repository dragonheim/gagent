package worker

import (
	fmt "fmt"
	log "log"
	http "net/http"
	sync "sync"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	// picol "git.dragonheim.net/dragonheim/gagent/src/picol"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

	zmq "github.com/pebbe/zmq4"
)

/*
 The "worker" processes the agent code. The worker nodes do not know
 anything about the network structure. Instead they know only to which
 router(s) they are connected. The worker will execute the agent code and
 pass the agent and it's results to a router.
*/
func Main(wg *sync.WaitGroup, config gs.GagentConfig, rid int) {
	defer wg.Done()
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("[INFO] Starting worker\n")

	// Generate connect string for this router.
	var rport = config.WorkerPort
	if config.Routers[rid].WorkerPort != 0 {
		rport = config.Routers[rid].WorkerPort
	}
	connectString := fmt.Sprintf("tcp://%s:%d", config.Routers[rid].RouterAddr, rport)
	// workerListener := fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.WorkerPort)
	clientListener := fmt.Sprintf("%s:%d", config.ListenAddr, config.ClientPort)

	subscriber, _ := zmq.NewSocket(zmq.REP)
	defer subscriber.Close()

	go func() {
		http.ListenAndServe(clientListener, nil)
	}()

	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	subscriber.Connect(connectString)

	msg, err := subscriber.Recv(0)
	if err != nil {
		log.Printf("[DEBUG] Received error: %v", err)
	}
	log.Printf("[DEBUG] Received message: %v", msg[0])
}
