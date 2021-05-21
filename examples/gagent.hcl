/*
 * This is the name of this node and is only used
 * for logging purposes.
 *
 * Optional.
 */
// name = "gagent-zulu.example.org"

/*
 * This is the mode that this node operates in. There
 * are three modes;
 *   client == Clients read the local agent file and
 *             forwards the contents on to a router
 *
 *   router == Routers accept agents from clients and
 *             other routers and accepts responses to
 *             agents from workers and other routers.
 *
 *   worker == Workers collect and process agents and
 *             send responses to routers for return
 *             the requesting client.
 *
 * If it is not defined, G'Agent will start in setup
 * mode and attempt to write a new configuration file
 * to the local directory.  The file will be called
 * gagent.hcl
 *
 * Required.
 */
mode = "router"

/*
 * @TODO: Add authentication based on UUID
 * This is the UUID used throughout the G'Agent system
 * to uniquely identify this node. It is generated
 * during setup if it doesn't exist.
 *
 * Required.
 */
// uuid = "7e9d13fe-5151-5876-66c0-20ca03e8fca4"

/*
 * This is the IP Address to bind to, it defaults to
 * 0.0.0.0
 *
 * Optional.
 */
// listenaddr =  0.0.0.0

/*
 * This is the port to the router will listen for on
 * for clients. It defaults to 35570.
 *
 * Optional.
 */
// clientport = 35571

/*
 * This is the port to the router will listen for on
 * for other routers. It defaults to 35570.
 *
 * Optional.
 */
// routerport = 35570

/*
 * This is the port to the router will listen for on
 * for workers. It defaults to 35571.
 *
 * Optional.
 */
// workerport = 35572

/*
 * @TODO
 * This is the list of known G'Agent clients. Clients
 * are not registered dynamically, instead the only
 * clients that may connect are those listed here,
 * but client's of other routers may route, via tags,
 * their agent here.
 *
 * Optional.
 */
// client "alpha" {
//   clientid = "04f97538-270d-4ce3-b782-e09ef35830e9"
// }

// client "beta" {
//   clientid = "04f97538-270d-4cf3-b782-e09ef35830e9"
// }

/*
 * This is a list of known G'Agent routers. At least
 * one router is required for workers and clients. If
 * there is more than one router, clients and workers
 * will connect to them in sequential order.
 */
//  router "alpha" {
//    routerid = "04f97538-270d-4cb3-b782-e09ef35830e9"
//    address = "gagent-alpha.example.org"
//  }

//  router "beta" {
//    routerid = "04f97538-270d-4cc3-b782-e09ef35830e9"
//    address = "gagent-beta.example.org"
//  }

/*
 * This is a list of known G'Agent workers. This is only
 * used by routers to determine which workers are
 * allowed to accept and respond to agents.
 *
 * At least one worker is reuqired for routers.
 */
// worker "alpha" {
//   workerid = "04f97538-270d-4ce3-b782-e09ef35830e9"
// }

// worker "beta" {
//   workerid = "04f97538-270d-4cf3-b782-e09ef35830e9"
// }
