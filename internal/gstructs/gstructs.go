package gstructs

// GagentConfig is the primary construct used by all modes
type GagentConfig struct {
	Name       string           `hcl:"name,optional"`
	Mode       string           `hcl:"mode,attr"`
	UUID       string           `hcl:"uuid,optional"`
	ListenAddr string           `hcl:"listenaddr,optional"`
	ClientPort int64            `hcl:"clientport,optional"`
	RouterPort int64            `hcl:"routerport,optional"`
	WorkerPort int64            `hcl:"workerport,optional"`
	Clients    []*ClientDetails `hcl:"client,block"`
	Routers    []*RouterDetails `hcl:"router,block"`
	Workers    []*WorkerDetails `hcl:"worker,block"`
	Version    string
	File       string
}

// ClientDetails is details about known clients
type ClientDetails struct {
	/*
	 * Client name for display purposes in logs and
	 * diagnostics.
	 */
	ClientName string `hcl:",label"`

	/*
	 * UUID String for the client node.  This is used by
	 * the router to determine which MQ client to send
	 * the agent's results to. This attempts to keep the
	 * clients unique globally.
	 */
	ClientID string `hcl:"clientid,optional"`
}

// RouterDetails is details about known routers
type RouterDetails struct {
	/*
	 * Router name for display purposes in logs and
	 * diagnostics
	 */
	RouterName string `hcl:",label"`

	/*
	 * UUID String for the router node.  This is used by
	 * the clients, routers, and workers to determine
	 * which MQ router to send the agent's requests to.
	 * This attempts to keep the routers unique globally.
	 */
	RouterID string `hcl:"routerid,attr"`

	/*
	 * This is the IP address or hostname the router
	 * will listen on.  The router will start up a 0MQ
	 * service that clients and workers will connect to.
	 */
	RouterAddr string `hcl:"address,attr"`

	/*
	 * This is the is the port that the router listens
	 * on for clients. If not defined, it will default
	 * to 35571.
	 */
	ClientPort int64 `hcl:"clientport,optional"`

	/*
	 * This is the is the port that the router listens
	 * on for routers. If not defined, it will default
	 * to 35570.
	 */
	RouterPort int64 `hcl:"workerport,optional"`

	/*
	 * This is the is the port that the router listens
	 * on for clients. If not defined, it will default
	 * to 35572.
	 */
	WorkerPort int64 `hcl:"workerport,optional"`

	/*
	 * These tags will be passed to the router upon
	 * connection.  The router will then use these
	 * tags to help determine which worker / client
	 * to send the client's requests and results to.
	 */
	RouterTags []string `hcl:"tags,optional"`
}

// WorkerDetails is details about known workers
type WorkerDetails struct {
	/*
	 * Router name for display purposes in logs and
	 * diagnostics
	 */
	WorkerName string `hcl:",label"`

	/*
	 * UUID String for the worker node.  This is used
	 * by the router to determine which MQ client to
	 * send agents to. This attempts to keep the
	 * workers unique globally.
	 */
	WorkerID string `hcl:"workerid,attr"`

	/*
	 * These tags will be passed to the router upon
	 * connection.  The router will then use these
	 * tags to help determine which worker / client
	 * to send the agent and it's results to.
	 */
	WorkerTags []string `hcl:"tags,optional"`
}
