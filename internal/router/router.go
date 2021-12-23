package router

import (
	fmt "fmt"
	log "log"
	http "net/http"
	sync "sync"

	gcdb "git.dragonheim.net/dragonheim/gagent/internal/chaindb"
	gstructs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	prometheus "github.com/prometheus/client_golang/prometheus"
	promauto "github.com/prometheus/client_golang/prometheus/promauto"

	zmq "github.com/pebbe/zmq4"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_requests_received",
	})

	db gcdb.GagentDb
)

/*
 The 'router' processes routing requests from the agent.  The router does
 not handle any of the agent activities beyond processing the agent's
 list of tags and passing the agent and it's storage to either a member
 or client node. Tags are used by the agent to give hints as to where
 it should be routed.
 Main is the entrypoint for the router
*/
func Main(wg *sync.WaitGroup, config gstructs.GagentConfig) {
	log.Printf("[INFO] Starting router\n")
	defer wg.Done()

	http.HandleFunc("/hello", answerClient)
	clientSock, _ := zmq.NewSocket(zmq.ROUTER)
	defer clientSock.Close()

	workerSock, _ := zmq.NewSocket(zmq.DEALER)
	defer workerSock.Close()

	db.Init()

	workerListener := fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.WorkerPort)
	_ = workerSock.Bind(workerListener)

	workers := make([]string, 0)

	poller1 := zmq.NewPoller()
	poller1.Add(workerSock, zmq.POLLIN)

	poller2 := zmq.NewPoller()
	poller2.Add(workerSock, zmq.POLLIN)

	wg.Add(1)
	go createClientListener(wg, config)

LOOP:
	for {
		/*
		 *  Poll frontend only if we have available workers
		 */
		var sockets []zmq.Polled
		var err error
		if len(workers) > 0 {
			sockets, err = poller2.Poll(-1)
		} else {
			sockets, err = poller1.Poll(-1)
		}
		if err != nil {
			/*
			 *  Interrupt
			 */
			break
		}
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case workerSock:
				/*
				 *  Handle worker activity on backend
				 *  Use worker identity for load-balancing
				 */
				msg, err := s.RecvMessage(0)
				if err != nil {
					/*
					 *  Interrupt
					 */
					break LOOP
				}
				var identity string
				identity, msg = unwrap(msg)
				log.Printf("[DEBUG] Worker message received: %s", msg)
				workers = append(workers, identity)

			case clientSock:
				wg.Add(1)
				go createClientListener(wg, config)
			}
		}
	}
}

/*
 * Create listener for client requests
 */
func createClientListener(wg *sync.WaitGroup, config gstructs.GagentConfig) {
	defer wg.Done()

	clientSock, _ := zmq.NewSocket(zmq.ROUTER)
	defer clientSock.Close()

	clientListener := fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.ClientPort)
	log.Printf("[DEBUG] Binding to: %s", clientListener)
	_ = clientSock.Bind(clientListener)

	for {
		msg, err := clientSock.RecvMessage(0)
		if err != nil {
			break
		}
		log.Printf("[DEBUG] Client message received: %s", msg)
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

func answerClient(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		opsProcessed.Inc()
		// fmt.Fprintf(w, "%v\n", r)
		http.NotFound(w, r)
		return
	}

	/*
	 * Common code for all requests can go here...
	 */
	switch r.Method {
	/*
	 * Handle GET requests
	 */
	case http.MethodGet:
		fmt.Fprintf(w, "%v\n", r)

	/*
	 * Handle POST requests
	 */
	case http.MethodPost:
		fmt.Fprintf(w, "%v\n", r)

	/*
	 * Handle PUT requests
	 */
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	/*
	 * Handle everything else
	 */
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
