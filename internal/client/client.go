package client

import (
	sha "crypto/sha256"
	hex "encoding/hex"
	fmt "fmt"
	log "log"
	os "os"
	regexp "regexp"
	strconv "strconv"
	strings "strings"
	sync "sync"
	time "time"

	gs "github.com/dragonheim/gagent/internal/gstructs"

	zmq "github.com/pebbe/zmq4"
)

/*
 * Client mode will send an agent file to a router for processing
 * Clients do not process the agent files, only send them as
 * requests to a router. If started without arguments, the client
 * will contact the router and attempt to retrieve the results
 * of it's most recent request.
 * Main is the entrypoint for the client process
 */
func Main(wg *sync.WaitGroup, config gs.GagentConfig) {
	log.Printf("[INFO] Starting client\n")
	defer wg.Done()

	var agent gs.AgentDetails
	var err error

	if config.CMode {
		agent.Script, err = os.ReadFile(config.Agent)
		if err != nil {
			log.Printf("[ERROR] No such file or directory: %s", config.Agent)
			os.Exit(4)
		}
		log.Printf("[DEBUG] Agent file contents: \n----- -----\n%s\n----- -----\n", agent.Script)
	}
	agent.Client = config.UUID
	tmpsum := sha.Sum256([]byte(agent.Script))
	agent.Shasum = fmt.Sprintf("%v", hex.EncodeToString(tmpsum[:]))
	log.Printf("[INFO] SHA256 of Agent file: %s", agent.Shasum)
	agent.Status = 1
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
		connectString := "tcp://" + config.Routers[key].RouterAddr + ":" + strconv.Itoa(rport)

		wg.Add(1)
		go sendAgent(wg, config.UUID, connectString, agent)
	}
}

/*
 * Parse Agent file for GHINT data to populate the G'Agent hints
 */
func getTagsFromHints(agent gs.AgentDetails) []string {
	var tags []string

	// Use named capture groups to extract the hints
	re := regexp.MustCompile(`^*set\s+GHINT\s*\[\s*split\s*"(?P<Hints>[^"]+)"\s*,\s*\]`)
	res := re.FindStringSubmatch(string(agent.Script))

	// If we don't have at least 2 matches, we have no hints
	if len(res) < 2 {
		log.Printf("[ERROR] Agent is missing GHINT tags")
		os.Exit(4)
	}

	// Use named capturing group index
	hintsIndex := re.SubexpIndex("Hints")
	tags = strings.Split(res[hintsIndex], ",")

	log.Printf("[DEBUG] G'Agent hints: %v\n", tags)

	return tags
}

func sendAgent(wg *sync.WaitGroup, uuid string, connectString string, agent gs.AgentDetails) {
	defer wg.Done()

	var mu sync.Mutex
	mu.Lock()

	sock, _ := zmq.NewSocket(zmq.REQ)
	defer sock.Close()

	_ = sock.SetIdentity(uuid)

	log.Printf("[DEBUG] Attempting to connect to %s\n", connectString)
	err := sock.Connect(connectString)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to %s\n", connectString)
		os.Exit(10)
	}

	log.Printf("[DEBUG] Start sending agent...\n")
	agent.Status = 2
	status, err := sock.SendMessage(agent)
	if err != nil {
		log.Printf("[ERROR] Failed to send agent to router\n")
		return
	}
	log.Printf("[DEBUG] Agent send status: %d\n", status)
	mu.Unlock()
	time.Sleep(10 * time.Second)

}
