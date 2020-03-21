package personDB

import (
	"RAS/database"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
)

var schemaPerson = `
		CREATE TABLE person (
    		user_id varchar(32) primary key,
			user_name varchar(32),
			department varchar(256),
			mobile varchar(32),
			
			server_ip varchar(32),
			vnc_display int, 
			server_user varchar(32),
			server_passwd varchar(32)
		);
		`

type Person struct {
	UserID     string `db:"user_id"`
	UserName   string `db:"user_name"`
	Department string `db:"department"`
	Mobile     string `db:"mobile"`

	ServerIP     string `db:"server_ip"`
	VncDisplay   int    `db:"vnc_display"`
	ServerUser   string `db:"server_user"`
	ServerPasswd string `db:"server_passwd"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func QueryPerson(id string, db *sqlx.DB) (p Person, err error) {
	err = db.Get(&p, "SELECT * FROM person WHERE user_id=$1", id)
	if err != nil {
		return Person{}, fmt.Errorf("query person in db error: %v", err)
	}
	return p, nil
}

func GetAllPersonInfo(db *sqlx.DB) (pp []Person, err error) {
	pp = []Person{}
	err = db.Select(&pp, "SELECT * FROM person")
	if err != nil {
		return nil, fmt.Errorf("get all person info from db error: %v", err)
	}
	return pp, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
func CsvFileToDb(csvFilePath string, db *sqlx.DB) (err error) {
	//先清空内存结构
	persons := make(map[string]Person)

	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("open csv file: %s error: %v", csvFilePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reader.ReadAll error: %v", err)
	}

	for _, item := range record {
		//从scv文件中逐条解析
		p := Person{}
		p.UserID = item[0]
		p.UserName = item[1]
		p.Department = item[2]
		p.Mobile = item[3]

		p.ServerIP = item[4]
		p.VncDisplay, err = strconv.Atoi(item[5])
		if err != nil {
			return fmt.Errorf("p.VncDisplay strconv.Atoi  error, UserID = %s, VncDisplay = %s: %v",
				item[0], item[5], err)
		}
		p.ServerUser = item[6]
		p.ServerPasswd = item[7]

		//在内存中插入
		persons[p.UserID] = p

		//在数据库中插入
		//假定table person已建立且为空表
		insertPerson := `INSERT INTO person (user_id, user_name, department, mobile, server_ip, vnc_display, server_user, server_passwd) VALUES (:user_id, :user_name, :department, :mobile, :server_ip, :vnc_display, :server_user, :server_passwd)`

		_, err := db.NamedExec(insertPerson, p)
		if err != nil {
			return fmt.Errorf("insert person basic info into database errror, UserID = %s :%v", p.UserID, err)
		}
	}
	return nil
}

func CreatePersonTable(db *sqlx.DB) (result sql.Result, err error) {
	return database.CreateTable(db, schemaPerson)
}

func DropPersonTable(db *sqlx.DB) (result sql.Result, err error) {
	return database.DropTable(db, "person")
}
