package main

import (
	"fmt"
	"io/ioutil"
	"os"

	//	"math/rand"
	//	"time"

	gs "git.dragonheim.net/dragonheim/gagent/src/gstructs"

	//	client "git.dragonheim.net/dragonheim/gagent/src/client"
	//	router "git.dragonheim.net/dragonheim/gagent/src/router"
	//	worker "git.dragonheim.net/dragonheim/gagent/src/worker"

	docopt "github.com/aviddiviner/docopt-go"
	hclsimple "github.com/hashicorp/hcl/v2/hclsimple"
	uuid "github.com/nu7hatch/gouuid"
)

var exitCodes = struct {
	m map[string]int
}{m: map[string]int{
	"SUCCESS":             0,
	"CONFIG_FILE_MISSING": 1,
	"SETUP_FAILED":        2,
	"INVALID_MODE":        3,
	"AGENT_LOAD_FAILED":   4,
	"AGENT_MISSING_TAGS":  5,
	"NO_ROUTERS_DEFINED":  6,
	"NO_WORKERS_DEFINED":  6,
}}

func main() {
	var config gs.GagentConfig
	var configFile string = "/etc/gagent/gagent.hcl"

	config.Name, _ = os.Hostname()

	/*
	 * Set a default UUID for this node.
	 * This is used throughout the G'Agent system to uniquely identify this node.
	 * It can be overriden in the configuration file by setting uuid
	 */
	// rand.Seed(time.Now().UnixNano())
	identity, _ := uuid.NewV5(uuid.NamespaceURL, []byte("gagent"+config.Name))
	config.UUID = identity.String()

	/*
	 * By default, we want to listen on all IP addresses. It can be overriden
	 * in the configuration file by setting listenaddr
	 */
	config.ListenAddr = "0.0.0.0"

	/*
	 * By default, G'Agent will use port 35570 to communicate with the routers,
	 * but you can override it by setting the listenport in the configuration
	 * file
	 */
	config.ListenPort = 35570

	/*
	 * Create a usage variable and then use that to declare the arguments and
	 * options.  This allows us to use DocOpt for consistency between usage help
	 * and available arguments / options.  Documentation is available at;
	 *   http://docopt.org/
	 */
	usage := "G'Agents \n"
	usage += "\n"
	usage += "  Go based mobile agent system, loosely inspired by the Agent Tcl / D'Agents \n"
	usage += "  system created by Robert S. Gray of Dartmouth college. \n"
	usage += "\n"

	usage += "Usage: \n"
	usage += "  gagent [--config=<config>] [--agent=<file>] \n"
	usage += "  gagent setup [--config=<config>] \n"
	usage += "\n"

	usage += "Arguments: \n"
	usage += "  client -- Start as a G'Agent client \n"
	usage += "  <file> -- filename of the agent to be uploaded to the G'Agent network \n"
	usage += "\n"

	usage += "  setup  -- Write inital configuration file \n"
	usage += "\n"

	usage += "Options:\n"
	usage += "  config=<config> [default: /etc/gagent/gagent.hcl] \n"

	/*
	 * Consume the usage variable and the command line arguments to create a
	 * dictionary of the command line arguments.
	 */
	arguments, _ := docopt.ParseDoc(usage)
	fmt.Printf("Arguments are %v\n", arguments)

	if arguments["--config"] != nil {
		configFile = arguments["--config"].(string)
	}

	/*
	 * Let the command line mode override the configuration.
	 */
	if arguments["setup"] == true {
		config.Mode = "setup"
	} else {
		err := hclsimple.DecodeFile(configFile, nil, &config)
		if err != nil {
			fmt.Printf("Failed to load configuration file: %s.\n", configFile)
			fmt.Println(err)
			os.Exit(exitCodes.m["CONFIG_FILE_MISSING"])
		}
	}

	fmt.Printf("Configuration is %v\n", config)

	switch config.Mode {
	case "client":
		/*
		 * Client mode will send an agent file to a router for processing
		 * Clients do not process the agent files, only send them as
		 * requests to a router. If started without arguments, the client
		 * will contact the router and attempt to retrieve the results
		 * of it's most recent request.
		 */
		fmt.Printf("Running in client mode\n")
		agent, err := ioutil.ReadFile(arguments["--agent"].(string))
		if err == nil {
			fmt.Printf("Agent containts %v\n", string(agent))
			// fmt.Printf("Forking...\n")
			// go client.Main(config.Client, string(agent))
			// fmt.Printf("Forked thread has completed\n")
			// time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Failed to load Agent file: %s.\n", arguments["--agent"].(string))
			os.Exit(exitCodes.m["AGENT_LOAD_FAILED"])
		}

	case "router":
		/*
		 * The 'router' processes routing requests from the agent.  The router does
		 * not handle any of the agent activities beyond processing the agent's
		 * list of tags and passing the agent and it's storage to either a member
		 * or client node. Tags are used by the agent to give hints as to where
		 * it should be routed.
		 */
		fmt.Printf("Running in router mode\n")
		// go router.Main(config.Router)
		// select {}

	case "worker":
		/*
		 * The 'worker' processes the agent code. The worker nodes do not know
		 * anything about the network structure. Instead they know only to which
		 * router(s) they are connected. The worker will execute the agent code and
		 * pass the agent and it's results to a router.
		 */
		fmt.Printf("Running in worker mode\n")
		// go worker.Main(config.Worker)
		// select {}

	case "setup":
		fmt.Printf("Running in setup mode\n")
		os.Exit(exitCodes.m["SETUP_FAILED"])

	default:
		fmt.Printf("Unknown operating mode, exiting.\n")
		os.Exit(exitCodes.m["INVALID_MODE"])
	}

	os.Exit(exitCodes.m["SUCCESS"])
}
