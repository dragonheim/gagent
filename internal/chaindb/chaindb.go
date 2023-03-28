package chaindb

import (
	sha256 "crypto/sha256"
	fmt "fmt"
	ioutil "io/ioutil"
	log "log"
	time "time"

	gstructs "github.com/dragonheim/gagent/internal/gstructs"
	cty "github.com/zclconf/go-cty/cty"

	hclsimple "github.com/hashicorp/hcl/v2/hclsimple"
	hclwrite "github.com/hashicorp/hcl/v2/hclwrite"
)

type GagentDb struct {
	ChainRow []*GagentDbRow `hcl:"timestamp,block"`
}

type GagentDbRow struct {
	Timestamp  time.Time             `hcl:"timestamp"`
	DBName     string                `hcl:"chainid,optional"`
	Agent      gstructs.AgentDetails `hcl:"agent,block"`
	DbCurrHash [32]byte              `hcl:"currhash"`
	DbPrevHash [32]byte              `hcl:"prevhash"`
}

/*
 * Initialize the database
 */
func NewGagentDb() *GagentDb {
	return &GagentDb{
		ChainRow: make([]*GagentDbRow, 0),
	}
}

/*
 * Load the database from disk
 */
func (db *GagentDb) LoadHCL(ChainDBPath string) error {
	err := hclsimple.DecodeFile(ChainDBPath, nil, db)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] DB values: %v\n", db)
	return nil
}

/*
 * Write the database to an HCL file
 */
func (db *GagentDb) WriteHCL(ChainDBPath string) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, row := range db.ChainRow {
		rowBlock := rootBody.AppendNewBlock("row", []string{})
		rowBody := rowBlock.Body()

		rowBody.SetAttributeValue("timestamp", cty.StringVal(row.Timestamp.Format(time.RFC3339)))
		rowBody.SetAttributeValue("chainid", cty.StringVal(row.DBName))
		rowBody.SetAttributeValue("currhash", cty.StringVal(fmt.Sprintf("%x", row.DbCurrHash)))
		rowBody.SetAttributeValue("prevhash", cty.StringVal(fmt.Sprintf("%x", row.DbPrevHash)))

		agentBlock := rowBody.AppendNewBlock("agent", []string{})
		agentBody := agentBlock.Body()
		agentBody.SetAttributeValue("name", cty.StringVal(row.Agent.Client))
		agentBody.SetAttributeValue("version", cty.StringVal(row.Agent.Shasum))
	}

	return ioutil.WriteFile(ChainDBPath, f.Bytes(), 0600)
}

/*
 * Add a new row to the chaindb
 */
func (db *GagentDb) AddRow(row *GagentDbRow) {
	row.Timestamp = time.Now()
	db.ChainRow = append(db.ChainRow, row)
	db.SetCurrHash()
	db.SetPrevHash()
}

/*
 * Set current hash of the database
 */
func (db *GagentDb) SetCurrHash() {
	row := db.ChainRow[len(db.ChainRow)-1]
	row.DbCurrHash = sha256.Sum256([]byte(fmt.Sprintf("%v", db)))
}

/*
 * Set previous hash of the database
 */
func (db *GagentDb) SetPrevHash() {
	row := db.ChainRow[len(db.ChainRow)-1]
	if len(db.ChainRow) > 1 {
		row.DbPrevHash = db.ChainRow[len(db.ChainRow)-2].DbCurrHash
	}
}
