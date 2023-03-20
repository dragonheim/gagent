package setup

import (
	log "log"
	sync "sync"

	cty "github.com/zclconf/go-cty/cty"

	gs "github.com/dragonheim/gagent/internal/gstructs"

	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"
)

/*
 * Main is the entrypoint for the setup process
 */
func Main(wg *sync.WaitGroup, config gs.GagentConfig) {
	log.Printf("[INFO] Starting setup\n")
	defer wg.Done()

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.SetAttributeValue("name", cty.StringVal(config.Name))
	rootBody.SetAttributeValue("mode", cty.StringVal(config.Mode))
	rootBody.SetAttributeValue("uuid", cty.StringVal(config.UUID))
	rootBody.SetAttributeValue("listenaddr", cty.StringVal("0.0.0.0"))
	rootBody.SetAttributeValue("clientport", cty.NumberIntVal(config.ClientPort))
	rootBody.SetAttributeValue("routerport", cty.NumberIntVal(config.RouterPort))
	rootBody.SetAttributeValue("workerport", cty.NumberIntVal(config.WorkerPort))
	rootBody.AppendNewline()

	clientBlock1 := rootBody.AppendNewBlock("client", []string{config.Name})
	clientBody1 := clientBlock1.Body()
	/*
	 * clientBody1.AppendUnstructuredTokens(
	 * 	hclwrite.TokensForTraversal(hcl.Traversal{
	 * 		hcl.TraverseRoot{
	 * 			Name: hcl.CommentGenerator("comment"),
	 * 		},
	 * 	},
	 * 	))
	 */
	clientBody1.SetAttributeValue("clientid", cty.StringVal(config.UUID))
	rootBody.AppendNewline()

	routerBlock1 := rootBody.AppendNewBlock("router", []string{config.Name})
	routerBody1 := routerBlock1.Body()
	routerBody1.SetAttributeValue("routerid", cty.StringVal(config.UUID))
	routerBody1.SetAttributeValue("address", cty.StringVal("127.0.0.1"))
	routerBody1.SetAttributeValue("clientport", cty.NumberIntVal(config.ClientPort))
	routerBody1.SetAttributeValue("routerport", cty.NumberIntVal(config.RouterPort))
	routerBody1.SetAttributeValue("workerport", cty.NumberIntVal(config.WorkerPort))
	rootBody.AppendNewline()

	workerBlock1 := rootBody.AppendNewBlock("worker", []string{config.Name})
	workerBody1 := workerBlock1.Body()
	workerBody1.SetAttributeValue("workerid", cty.StringVal(config.UUID))
	rootBody.AppendNewline()

	log.Printf("[DEBUG] Configuration file created;\n\n%s\n", f.Bytes())
}
