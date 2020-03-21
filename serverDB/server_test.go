package serverDB

import (
	"RAS/database"
	"testing"
)

func TestCsvFileToDb(t *testing.T) {
	driverName, dataSourceName, err := database.LoadPostgreSource()
	if err != nil {
		t.Errorf("database.LoadPostgreSource error: %v", err)
	}

	db, err := database.ConnectToDatabase(driverName, dataSourceName)
	if err != nil {
		t.Errorf("ConnectToDatabase error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("db.Ping error: %v", err)
	}

	result, err := DropServerTable(db)
	if err != nil {
		t.Errorf("DropServerTable error: %v", err)
	}
	t.Logf("DropServerTable success: %v", result)

	result, err = CreateServerTable(db)
	if err != nil {
		t.Errorf("CreateServerTable error: %v", err)
	}
	t.Logf("CreateServerTable success: %v", result)

	err = CsvFileToDb("/TestDataset/servers.csv", db)
	if err != nil {
		t.Errorf("CsvFileToDb error: %v", err)
	}
	t.Logf("CsvFileToDb success")

	pp, err := GetAllServerInfo(db)
	if err != nil {
		t.Errorf("GetAllServerInfo error: %v", err)
	}
	t.Logf("GetAllServerInfo success:%v", pp)

	serverName := "TestServer001"
	p, err := QueryServer(serverName, db)
	if err != nil {
		t.Errorf("QueryServer <%s> error: %v", serverName, err)
	}
	t.Logf("QueryServer <%s> success:%v", serverName, p)
}
