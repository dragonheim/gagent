/*
 * This is the name of this node and is only used
 * for logging purposes.
 *
 * Optional.
 */
name = "gagent-zulu.example.org"

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
 * Required.
 */
mode = "router"

/*
 * This is the UUID used throughout the G'Agent system
 * to uniquely identify this node.
 *
 * Required.
 */
// uuid = "04f97538-270d-4ca3-b782-e09ef35830e9"

/*
 * This is a list of known G'Agent routers. At least
 * one router is required for workers and clients. If
 * there is more than one router, clients and workers
 * will connect to them in sequential order.
 */
// router "alpha" {
// 	routerid = "04f97538-270d-4cb3-b782-e09ef35830e9"
// 	address = "gagent-alpha.example.org"
// 	tags = [ "a", "b", "c", "d" ]
// }
// 
// router "beta" {
// 	routerid = "04f97538-270d-4cc3-b782-e09ef35830e9"
// 	address = "gagent-beta.example.org"
// 	tags = [ "a", "c", "e", "g" ]
// }
// 
// router "charlie" {
// 	routerid = "04f97538-270d-4cd3-b782-e09ef35830e9"
// 	address = "gagent-charlie.example.org"
// 	tags = [ "b", "d", "f", "h" ]
// }

/*
 * This is a list of known G'Agent workers. This is only
 * used by routers to determine which workers are
 * allowed to accept and respond to agents.
 *
 * At least one worker is reuqired for routers.
 */
// worker "alpha" {
// 	workerid = "04f97538-270d-4ce3-b782-e09ef35830e9"
// 	address = "gagent-alpha.example.org"
// 	tags = [ "a", "b", "c", "d" ]
// }
// 
// worker "beta" {
// 	workerid = "04f97538-270d-4cf3-b782-e09ef35830e9"
// 	address = "gagent-beta.example.org"
// 	tags = [ "a", "c", "e", "g" ]
// }
// 
// worker "charlie" {
// 	workerid = "04f97538-270d-4c04-b782-e09ef35830e9"
// 	address = "gagent-charlie.example.org"
// 	tags = [ "b", "d", "f", "h" ]
// }

