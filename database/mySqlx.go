package database

import (
	"RAS/myUtils"
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

func LoadPostgreSource() (driverName, dataSourceName string, err error) {
	//Load IdKey from file
	filename := "/IdKey/Database/ras_pg.json"
	idKey, err := myUtils.LoadIdKey(filename)
	if err != nil {
		return "", "", fmt.Errorf("load PostgreSQL source from file error: %v", err)
	}

	dataSourceName = fmt.Sprintf(idKey.SecretId, idKey.SecretKey)
	return "postgres", dataSourceName, nil
	//"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	//"host=ras-pg user=postgres password=%s dbname=ras sslmode=disable"
}
