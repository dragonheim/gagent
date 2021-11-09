package client

import (
	sha "crypto/sha256"
	fmt "fmt"
	ioutil "io/ioutil"
	log "log"
	os "os"
	regexp "regexp"
	strings "strings"
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
			os.Exit(4)
		}
		log.Printf("[DEBUG] Agent file contents: \n----- -----\n%s\n----- -----\n", agent.ScriptCode)
	}
	agent.Client = config.UUID
	agent.Shasum = fmt.Sprintf("%x", sha.Sum256(agent.ScriptCode))
	log.Printf("[INFO] SHA256 of Agent file: %s", agent.Shasum)
	agent.Status = "loaded"
	agent.Hints = getTagsFromHints(agent)
	agent.Answer = nil

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
		go sendAgent(wg, config.UUID, connectString, agent)
	}
}

/*
 * Parse Agent file for GHINT data to populate the G'Agent hints
 */
func getTagsFromHints(agent gstructs.AgentDetails) []string {
	var tags []string
	re := regexp.MustCompile(`\s*set\s+GHINT\s*\[\s*split\s*"(?P<Hints>.+)"\s*\,\s*\]`)
	res := re.FindStringSubmatch(string(agent.ScriptCode))
	if len(res) < 1 {
		log.Printf("[ERROR] Agent is missing GHINT tags")
		os.Exit(4)
	}
	tags = strings.Split(res[1], ",")
	log.Printf("[DEBUG] G'Agent hints: %v\n", tags)

	return tags
}

func sendAgent(wg *sync.WaitGroup, uuid string, connectString string, agent gstructs.AgentDetails) {
	defer wg.Done()

	var mu sync.Mutex
	mu.Lock()

	sock, _ := zmq.NewSocket(zmq.REQ)
	defer sock.Close()

	sock.SetIdentity(uuid)

	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	err := sock.Connect(connectString)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to %s\n", connectString)
		os.Exit(10)
	}

	log.Printf("[DEBUG] Start sending agent...\n")
	status, err := sock.SendMessage(agent)
	if err != nil {
		log.Printf("[ERROR] Failed to send agent to router\n")
		// os.Exit(11)
		return
	}
	log.Printf("[DEBUG] Agent send status: %d\n", status)
	mu.Unlock()
	time.Sleep(10 * time.Second)

}
