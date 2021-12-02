package chaindb

import (
	sha "crypto/sha256"
	fmt "fmt"
	log "log"
	time "time"

	gstructs "git.dragonheim.net/dragonheim/gagent/internal/gstructs"

	hclsimple "github.com/hashicorp/hcl/v2/hclsimple"
	// hclwrite "github.com/hashicorp/hcl/v2/hclwrite"
)

type GagentDb struct {
	chainRow []*gagentDbRow `hcl:"timestamp,block"`
}

type gagentDbRow struct {
	timestamp  time.Time             `hcl:"timestamp"`
	DBName     string                `hcl:"chainid,optional"`
	Agent      gstructs.AgentDetails `hcl:"agent,block"`
	dbCurrHash [32]byte              `hcl:"currhash"`
	dbPrevHash [32]byte              `hcl:"prevhash"`
}

/*
 * Initialize the database
 */
func (db *GagentDb) Init() {
	db.chainRow = make([]*gagentDbRow, 0)
}

/*
 * Load the database from disk
 */
func (db *GagentDb) Load() error {
	err := hclsimple.DecodeFile("chaindb.hcl", nil, &db)
	log.Printf("[DEBUG] DB values: %v\n", db)
	return err
}

/*
 * Add a new row to the chaindb
 */
func (db *GagentDb) AddRow(row *gagentDbRow) error {
	row.timestamp = time.Now()
	db.chainRow = append(db.chainRow, row)

	return nil
}

/*
 * Set current hash of the database
 */
func (db *GagentDb) SetCurrHash() {
	db.chainRow[len(db.chainRow)-1].dbCurrHash = [32]byte{}
	foo := sha.Sum256([]byte(fmt.Sprintf("%v", db)))
	db.chainRow[len(db.chainRow)-1].dbCurrHash = foo
}

/*
 * Set previous hash of the database
 */
func (db *GagentDb) SetPrevHash() {
	db.chainRow[len(db.chainRow)-1].dbPrevHash = db.chainRow[len(db.chainRow)-1].dbCurrHash
}
