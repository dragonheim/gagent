package client

import (
	sha "crypto/sha256"
	fmt "fmt"
	ioutil "io/ioutil"
	log "log"
	os "os"
	sync "sync"
	time "time"

	gstructs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	zmq "github.com/pebbe/zmq4"
)

/*
 Client mode will send an agent file to a router for processing
 Clients do not process the agent files, only send them as
 requests to a router. If started without arguments, the client
 will contact the router and attempt to retrieve the results
 of it's most recent request.
*/
func Main(wg *sync.WaitGroup, config gstructs.GagentConfig) {
	log.Printf("[INFO] Starting client\n")
	defer wg.Done()

	var agent gstructs.AgentDetails
	var err error

	if config.CMode {
		agent.ScriptCode, err = ioutil.ReadFile(config.File)
		if err != nil {
			log.Printf("[ERROR] No such file or directory: %s", config.File)
			os.Exit(6)
		}
		agent.Shasum = fmt.Sprintf("%x", sha.Sum256(agent.ScriptCode))
		agent.Status = "loaded"
		log.Printf("[DEBUG] SHA256 of Agent file: %s", agent.Shasum)
	}

	for key := range config.Routers {
		/*
		 * Generate connect string for this router.
		 */
		rport := config.ClientPort
		if config.Routers[key].ClientPort != 0 {
			rport = config.Routers[key].ClientPort
		}
		connectString := fmt.Sprintf("tcp://%s:%d", config.Routers[key].RouterAddr, rport)

		wg.Add(1)
		go sendAgent(wg, config.UUID, connectString, agent.ScriptCode)
		time.Sleep(10 * time.Millisecond)
	}
}

func sendAgent(wg *sync.WaitGroup, uuid string, connectString string, agent []byte) {
	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	defer wg.Done()

	var mu sync.Mutex
	mu.Lock()

	sock, _ := zmq.NewSocket(zmq.REQ)
	defer sock.Close()

	sock.SetIdentity(uuid)
	sock.Connect(connectString)

	log.Printf("[DEBUG] Start sending agent...\n")
	sock.SendMessage(agent)
	log.Printf("[DEBUG] End sending agent...\n")
	mu.Unlock()
}
