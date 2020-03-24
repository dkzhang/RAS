package myPostgreSQL

import "testing"

func TestConnectToDatabase(t *testing.T) {
	db, err := ConnectToDatabase("postgres", "user=postgres dbname=mydb password=Jim980911 sslmode=disable")
	if err != nil {
		t.Errorf("ConnectToDatabase error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Errorf("db.Ping error: %v", err)
	}
}
