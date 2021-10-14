package client

import (
	fmt "fmt"
	log "log"
	http "net/http"
	sync "sync"
	time "time"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	zmq "github.com/pebbe/zmq4"
)

// Main is the initiation function for a Client
func Main(wg *sync.WaitGroup, config gs.GagentConfig, rid int, agent string) {
	defer wg.Done()
	log.Printf("[INFO] Starting client\n")

	// Generate connect string for this router.
	var rport = int64(config.ClientPort)
	if config.Routers[rid].ClientPort != 0 {
		rport = config.Routers[rid].ClientPort
	}
	connectString := fmt.Sprintf("tcp://%s:%d", config.Routers[rid].RouterAddr, rport)
	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)

	var mu sync.Mutex

	sock, _ := zmq.NewSocket(zmq.REQ)
	defer sock.Close()

	sock.SetIdentity(config.UUID)
	sock.Connect(connectString)

	go func() {
		mu.Lock()
		log.Printf("[DEBUG] Start sending agent...\n")
		sock.SendMessage(agent)
		log.Printf("[DEBUG] End sending agent...\n")
		mu.Unlock()
	}()

	//	time.Sleep(10 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	// for {
	// 	time.Sleep(10 * time.Millisecond)
	// 	mu.Lock()
	// 	msg, err := sock.RecvMessage(zmq.DONTWAIT)
	// 	if err == nil {
	// 		log.Println(msg[0], config.UUID)
	// 	}
	// 	mu.Unlock()
	// }

}

func pushAgent(config gs.GagentConfig) {
	http.Get(config.Routers[0].RouterAddr)
}
