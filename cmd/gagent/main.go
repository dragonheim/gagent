package main

import (
	fmt "fmt"
	ioutil "io/ioutil"
	log "log"
	http "net/http"
	os "os"
	sync "sync"

	autorestart "github.com/slayer/autorestart"

	gs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	gc "git.dragonheim.net/dragonheim/gagent/internal/client"
	gr "git.dragonheim.net/dragonheim/gagent/internal/router"
	gw "git.dragonheim.net/dragonheim/gagent/internal/worker"

	cty "github.com/zclconf/go-cty/cty"

	docopt "github.com/aviddiviner/docopt-go"

	hclsimple "github.com/hashicorp/hcl/v2/hclsimple"
	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"

	logutils "github.com/hashicorp/logutils"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

	uuid "github.com/jakehl/goid"
)

var (
	semVER = "0.0.2"
)

var (
	wg sync.WaitGroup
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
	"NO_WORKERS_DEFINED":  7,
	"AGENT_NOT_DEFINED":   8,
}}

var config gs.GagentConfig

func main() {
	autorestart.StartWatcher()
	http.Handle("/metrics", promhttp.Handler())

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	var agent gs.Agent
	// config.File = "/etc/gagent/gagent.hcl"

	// config.Name, _ = os.Hostname()
	// config.Mode = "setup"

	/*
	 * Set a default UUID for this node.
	 * This is used throughout the G'Agent system to uniquely identify this node.
	 * It can be overridden in the configuration file by setting uuid
	 */
	// identity := uuid.NewV4UUID()
	// config.UUID = identity.String()

	/*
	 * By default, we want to listen on all IP addresses. It can be overridden
	 * in the configuration file by setting listenaddr
	 */
	// config.ListenAddr = "0.0.0.0"

	/*
	 * By default, G'Agent client will use port 35571 to communicate with the
	 * routers, but you can override it by setting the clientport in the
	 * configuration file
	 */
	// config.ClientPort = 35571

	/*
	 * By default, G'Agent router will use port 35572 to communicate with
	 * other routers, but you can override it by setting the routerport in
	 * the configuration file
	 */
	// config.RouterPort = 35570

	/*
	 * By default, G'Agent worker will use port 35570 to communicate with the
	 * routers, but you can override it by setting the workerport in the
	 * configuration file
	 */
	// config.WorkerPort = 35572

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

	log.Printf("[DEBUG] Config is %v\n", config)
	if opts["--config"] != nil {
		config.File = opts["--config"].(string)
	}

	/*
	 * Start Prometheus metrics exporter
	 */
	go http.ListenAndServe(fmt.Sprintf("%s:%d", config.ListenAddr, config.ClientPort), nil)

	/*
	 * Let the command line mode override the configuration.
	 */
	if opts["setup"] == true {
		config.Mode = "setup"
	} else {
		err := hclsimple.DecodeFile(config.File, nil, &config)
		if err != nil {
			log.Printf("[ERROR] Failed to load configuration file: %s.\n%s\n", config.File, err)
			os.Exit(exitCodes.m["CONFIG_FILE_MISSING"])
		}
		if opts["client"] == true {
			config.Mode = "client"
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
	config.Version = semVER
	log.Printf("[DEBUG] Configuration is %v\n", config)

	switch config.Mode {
	case "client":
		log.Printf("[INFO] Running in client mode\n")

		if len(config.Routers) == 0 {
			log.Printf("[ERROR] No routers defined.\n")
			os.Exit(exitCodes.m["NO_ROUTERS_DEFINED"])
		}

		if opts["--agent"] == nil {
			log.Printf("[ERROR] Agent file not specified")
			os.Exit(exitCodes.m["AGENT_NOT_DEFINED"])
		}

		// var agent []byte
		var err error
		if config.CMode {
			agent.ScriptCode, err = ioutil.ReadFile(opts["--agent"].(string))
			if err != nil {
				log.Printf("[ERROR] Failed to load Agent file: %s", opts["--agent"])
				os.Exit(exitCodes.m["AGENT_LOAD_FAILED"])
			}
		}

		wg.Add(1)
		go gc.Main(&wg, config, string(agent.ScriptCode))

	case "router":
		log.Printf("[INFO] Running in router mode\n")

		if len(config.Workers) == 0 {
			log.Printf("[ERROR] No workers defined.\n")
			os.Exit(exitCodes.m["NO_WORKERS_DEFINED"])
		}

		wg.Add(1)
		go gr.Main(&wg, config)

	case "worker":
		log.Printf("[INFO] Running in worker mode\n")

		if len(config.Routers) == 0 {
			log.Printf("[ERROR] No routers defined.\n")
			os.Exit(exitCodes.m["NO_ROUTERS_DEFINED"])
		}

		wg.Add(1)
		go gw.Main(&wg, config)

	case "setup":
		log.Printf("[INFO] Running in setup mode\n")
		f := hclwrite.NewEmptyFile()
		rootBody := f.Body()
		rootBody.SetAttributeValue("name", cty.StringVal(config.Name))
		rootBody.SetAttributeValue("mode", cty.StringVal(config.Mode))
		rootBody.SetAttributeValue("uuid", cty.StringVal(config.UUID))
		rootBody.AppendNewline()

		routerBlock1 := rootBody.AppendNewBlock("router", []string{config.Name})
		routerBody1 := routerBlock1.Body()
		routerBody1.SetAttributeValue("routerid", cty.StringVal(config.UUID))
		routerBody1.SetAttributeValue("address", cty.StringVal("127.0.0.1"))
		rootBody.AppendNewline()

		log.Printf("\n%s", f.Bytes())
		os.Exit(exitCodes.m["SUCCESS"])

	default:
		log.Printf("[ERROR] Unknown operating mode, exiting.\n")
		os.Exit(exitCodes.m["INVALID_MODE"])
	}

	wg.Wait()
	os.Exit(exitCodes.m["SUCCESS"])
}

func init() {
	// Initialize the configuration
	config.Mode = "setup"
	config.Name, _ = os.Hostname()
	config.UUID = uuid.NewV4UUID().String()
	config.ListenAddr = "0.0.0.0"
	config.ClientPort = 35571
	config.RouterPort = 35570
	config.WorkerPort = 35572
	config.Clients = make([]*gs.ClientDetails, 0)
	config.Routers = make([]*gs.RouterDetails, 0)
	config.Workers = make([]*gs.WorkerDetails, 0)
	config.File = "/etc/gagent/gagent.hcl"
	config.Version = semVER

	// Initialize the exit codes
	// exitCodes.m = make(map[string]int)
	exitCodes.m["SUCCESS"] = 0
	exitCodes.m["INVALID_MODE"] = 1
	exitCodes.m["CONFIG_FILE_MISSING"] = 2
	exitCodes.m["NO_ROUTERS_DEFINED"] = 3
	exitCodes.m["NO_WORKERS_DEFINED"] = 4
	exitCodes.m["AGENT_NOT_DEFINED"] = 5
	exitCodes.m["AGENT_LOAD_FAILED"] = 6

}
