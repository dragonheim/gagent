package main

import (
	fmt "fmt"
	log "log"
	http "net/http"
	os "os"
	sync "sync"

	autorestart "github.com/slayer/autorestart"

	fqdn "github.com/Showmax/go-fqdn"

	gstructs "github.com/dragonheim/gagent/internal/gstructs"

	gc "github.com/dragonheim/gagent/internal/client"
	gr "github.com/dragonheim/gagent/internal/router"
	gs "github.com/dragonheim/gagent/internal/setup"
	gw "github.com/dragonheim/gagent/internal/worker"

	docopt "github.com/aviddiviner/docopt-go"

	hclsimple "github.com/hashicorp/hcl/v2/hclsimple"

	logutils "github.com/hashicorp/logutils"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

	uuid "github.com/jakehl/goid"
)

var (
	semVER = "0.0.6"
)

var (
	wg sync.WaitGroup
)

/*
 * Exit Codes
 *  0 Success
 *  1 Configuration file is missing or unreadable
 *  2 Setup failed
 *  3 Invalid mode of operation
 *  4 Agent file is missing or unreadable
 *  5 Agent is missing tags
 *  6 No routers defined
 *  7 No workers defined
 *  8 Agent not defined
 *  9 Agent hints / tags not defined
 * 10 Router not connected
 */

var config gstructs.GagentConfig

func main() {
	log.Printf("[DEBUG] Configuration is %v\n", config)

	switch config.Mode {
	case "client":
		log.Printf("[INFO] Running in client mode\n")

		if len(config.Routers) == 0 {
			log.Printf("[ERROR] No routers defined.\n")
			os.Exit(6)
		}

		wg.Add(1)
		go gc.Main(&wg, config)

	case "router":
		log.Printf("[INFO] Running in router mode\n")

		if len(config.Workers) == 0 {
			log.Printf("[ERROR] No workers defined.\n")
			os.Exit(7)
		}

		// cs.Client.ConsulRegister("gagent-router", config.RouterPort)

		wg.Add(1)
		go gr.Main(&wg, config)

	case "worker":
		log.Printf("[INFO] Running in worker mode\n")

		if len(config.Routers) == 0 {
			log.Printf("[ERROR] No routers defined.\n")
			os.Exit(6)
		}

		// cs.Client.ConsulRegister("gagent-worker", config.WorkerPort)

		wg.Add(1)
		go gw.Main(&wg, config)

	case "setup":
		log.Printf("[INFO] Running in setup mode\n")

		wg.Add(1)
		go gs.Main(&wg, config)

	default:
		log.Printf("[ERROR] Unknown operating mode, exiting.\n")
		os.Exit(3)
	}

	wg.Wait()
	os.Exit(0)
}

func init() {
	// var err error

	autorestart.StartWatcher()

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	http.Handle("/metrics", promhttp.Handler())

	/*
	 * Initialize the configuration
	 */
	config.Version = semVER

	config.File = "/etc/gagent/gagent.hcl"

	config.Mode = "setup"

	config.Name, _ = fqdn.FqdnHostname()

	/*
	 * Set a default UUID for this node.
	 * This is used throughout the G'Agent system to uniquely identify this node.
	 * It can be overridden in the configuration file by setting uuid
	 */
	config.UUID = uuid.NewV4UUID().String()

	/*
	 * By default, we want to listen on all IP addresses.
	 */
	config.ListenAddr = "0.0.0.0"

	/*
	 * G'Agent will use this port for monitoring via prometheus., If set
	 * is set to 0, G'Agent will not listen for prometheus metrics.
	 */
	config.MonitorPort = 9101

	/*
	 * G'Agent client will use this port to communicate with the routers.
	 */
	config.ClientPort = 35572

	/*
	 * G'Agent router will use this port to communicate with other routers.
	 */
	config.RouterPort = 35570

	/*
	 * G'Agent worker will use this port to communicate with the routers.
	 */
	config.WorkerPort = 35571

	config.Clients = make([]*gstructs.ClientDetails, 0)
	config.Routers = make([]*gstructs.RouterDetails, 0)
	config.Workers = make([]*gstructs.WorkerDetails, 0)

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
	usage += "  gagent client (pull|push) [--config=<config>] [--agent=<file>] \n"
	usage += "  gagent router [--config=<config>] \n"
	usage += "  gagent worker [--config=<config>] \n"
	usage += "  gagent setup [--config=<config>] \n"
	usage += "  gagent --version \n"
	usage += "\n"

	usage += "Arguments: \n"
	usage += "  client pull       -- Start as a G'Agent client to pull agent results \n"
	usage += "  client push       -- Start as a G'Agent client to push agent \n"
	usage += "  router            -- Start as a G'Agent router \n"
	usage += "  worker            -- Start as a G'Agent worker \n"
	usage += "  setup             -- Write initial configuration file \n"
	usage += "\n"

	usage += "Options:\n"
	usage += "  -h --help         -- Show this help screen and exit \n"
	usage += "  --version         -- Show version and exit \n"
	usage += "  --config=<config> -- [default: /etc/gagent/gagent.hcl] \n"
	usage += "  --agent=<file>    -- filename of the agent to be uploaded to the G'Agent network. Required in push mode \n"
	usage += "\n"

	/*
	 * Consume the usage variable and the command line arguments to create a
	 * dictionary / map.
	 */
	opts, _ := docopt.ParseArgs(usage, nil, semVER)
	log.Printf("[DEBUG] Arguments are %v\n", opts)

	if opts["--config"] != nil {
		config.File = opts["--config"].(string)
	}

	err := hclsimple.DecodeFile(config.File, nil, &config)
	if err != nil && opts["setup"] == false {
		log.Printf("[ERROR] Failed to load configuration file: %s.\n", config.File)
		os.Exit(1)
	}

	/*
	 * Let the command line mode override the configuration.
	 */
	if opts["setup"] == true {
		config.Mode = "setup"
	} else {
		if opts["client"] == true {
			config.Mode = "client"
			if opts["--agent"] == nil {
				log.Printf("[ERROR] Agent file not specified")
				os.Exit(8)
			} else {
				config.File = opts["--agent"].(string)
			}

			if opts["pull"] == true {
				config.CMode = false
			}

			if opts["push"] == true {
				config.CMode = true
			}
		}

		if opts["router"] == true {
			config.Mode = "router"
		}

		if opts["worker"] == true {
			config.Mode = "worker"
		}
	}

	log.Printf("[DEBUG] Config is %v\n", config)

	/*
	 * Start Prometheus metrics exporter
	 */
	if config.MonitorPort != 0 {
		go func() {
			log.Printf("[INFO] Starting Prometheus metrics exporter on port %d\n", config.MonitorPort)
			log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenAddr, config.MonitorPort), nil))
		}()
	}
}
