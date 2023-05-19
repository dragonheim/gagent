package chaindb

import (
	"bytes"
	"os"
	"testing"
	"time"

	gstructs "github.com/dragonheim/gagent/internal/gstructs"
)

const testChainDBPath = "/tmp/test_chaindb.hcl"

func TestGagentDb(t *testing.T) {
	// Create a new GagentDb
	db := NewGagentDb()

	// Add a row to the database
	row := &GagentDbRow{
		DBName: "testDB",
		Agent: gstructs.AgentDetails{
			Client: "testAgent",
			Shasum: "v1.0.0",
		},
	}
	db.AddRow(row)

	// Check if the row was added correctly
	if len(db.ChainRow) != 1 {
		t.Errorf("Expected length of ChainRow to be 1, but got %d", len(db.ChainRow))
	}

	// Check if the timestamp was set correctly
	if db.ChainRow[0].Timestamp.After(time.Now()) {
		t.Error("Timestamp is incorrectly set in the future")
	}

	// Write the database to an HCL file
	err := db.WriteHCL(testChainDBPath)
	if err != nil {
		t.Errorf("Error writing HCL file: %v", err)
	}

	// Load the database from the HCL file
	loadedDb := NewGagentDb()
	err = loadedDb.LoadHCL(testChainDBPath)
	if err != nil {
		t.Errorf("Error loading HCL file: %v", err)
	}

	// Check if the loaded database is the same as the original one
	if !bytes.Equal(loadedDb.ChainRow[0].DbCurrHash[:], db.ChainRow[0].DbCurrHash[:]) {
		t.Error("Loaded database has a different current hash than the original one")
	}

	if !bytes.Equal(loadedDb.ChainRow[0].DbPrevHash[:], db.ChainRow[0].DbPrevHash[:]) {
		t.Error("Loaded database has a different previous hash than the original one")
	}

	// Clean up the test HCL file
	err = os.Remove(testChainDBPath)
	if err != nil {
		t.Errorf("Error cleaning up test HCL file: %v", err)
	}
}
