package router

import (
	fmt "fmt"
	log "log"
	http "net/http"
	sync "sync"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

	zmq "github.com/pebbe/zmq4"
)

// @TODO -- This was documented in the example, and I am unclear what it does
const (
	WORKER_READY = "\001" //  Signals worker is ready
)

// Main is the initiation function for a Router
func Main(wg *sync.WaitGroup, config gs.GagentConfig) {
	defer wg.Done()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", answerClient)

	log.Printf("[INFO] Starting router\n")

	clientSock, _ := zmq.NewSocket(zmq.ROUTER)
	defer clientSock.Close()

	workerSock, _ := zmq.NewSocket(zmq.DEALER)
	defer workerSock.Close()
	workerListener := fmt.Sprintf("tcp://%s:%d", config.ListenAddr, config.WorkerPort)
	clientListener := fmt.Sprintf("%s:%d", config.ListenAddr, config.ClientPort)

	go func() {
		http.ListenAndServe(clientListener, nil)
	}()

	workerSock.Bind(workerListener)

	workers := make([]string, 0)

	poller1 := zmq.NewPoller()
	poller1.Add(workerSock, zmq.POLLIN)

	poller2 := zmq.NewPoller()
	poller2.Add(workerSock, zmq.POLLIN)

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
			case workerSock:
				//  Handle worker activity on backend
				//  Use worker identity for load-balancing
				msg, err := s.RecvMessage(0)
				if err != nil {
					//  Interrupted
					break LOOP
				}
				var identity string
				identity, msg = unwrap(msg)
				log.Printf("[DEBUG] Worker message received: %s", msg)
				workers = append(workers, identity)

				//  Forward message to client if it's not a READY
				// if msg[0] != WORKER_READY {
				// 	clientSock.SendMessage(msg)
				// }

			case clientSock:
				//  Get client request, route to first available worker
				msg, err := s.RecvMessage(0)
				log.Printf("[DEBUG] Client message received: %s", msg)
				if err == nil {
					workerSock.SendMessage(workers[0], "", msg)
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

func answerClient(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprintf(w, "%v\n", r)
		// http.NotFound(w, r)
		return
	}

	// Common code for all requests can go here...

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "%v\n", r)
		// Handle the GET request...

	case http.MethodPost:
		// Handle the POST request...

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
