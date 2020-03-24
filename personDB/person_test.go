package personDB

import (
	"RAS/myPostgreSQL"
	"testing"
)

func TestCsvFileToDb(t *testing.T) {
	driverName, dataSourceName, err := myPostgreSQL.LoadPostgreSource()
	if err != nil {
		t.Errorf("myPostgreSQL.LoadPostgreSource error: %v", err)
	}

	db, err := myPostgreSQL.ConnectToDatabase(driverName, dataSourceName)
	if err != nil {
		t.Errorf("ConnectToDatabase error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("db.Ping error: %v", err)
	}

	result, err := DropPersonTable(db)
	if err != nil {
		t.Errorf("DropPersonTable error: %v", err)
	}
	t.Logf("DropPersonTable success: %v", result)

	result, err = CreatePersonTable(db)
	if err != nil {
		t.Errorf("CreatePersonTable error: %v", err)
	}
	t.Logf("CreatePersonTable success: %v", result)

	err = CsvFileToDb("/TestDataset/persons.csv", db)
	if err != nil {
		t.Errorf("CsvFileToDb error: %v", err)
	}
	t.Logf("CsvFileToDb success")

	pp, err := GetAllPersonInfo(db)
	if err != nil {
		t.Errorf("GetAllPersonInfo error: %v", err)
	}
	t.Logf("GetAllPersonInfo success:%v", pp)

	userId := "zhangjun"
	p, err := QueryPerson(userId, db)
	if err != nil {
		t.Errorf("QueryPerson <%s> error: %v", userId, err)
	}
	t.Logf("QueryPerson <%s> success:%v", userId, p)
}
