package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDatabase(driverName, dataSourceName string) (db *sqlx.DB, err error) {
	db, err = sqlx.Open(driverName, dataSourceName)
	return db, err
}

func CreateTable(db *sqlx.DB, schema string) (result sql.Result, err error) {
	result, err = db.Exec(schema)
	return result, err
}

func DropTable(db *sqlx.DB, tableName string) (result sql.Result, err error) {
	exec := `DROP Table ` + tableName
	result, err = db.Exec(exec)
	return result, err
}

func LoadDatabaseSourceConfig(filename string) (driverName, dataSourceName string) {
	fmt.Printf("config filename = %s", filename)
	return "postgres", "user=postgres dbname=ras password=Jim980911 sslmode=disable"
}
