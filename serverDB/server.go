package serverDB

import (
	"RAS/myPostgreSQL"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

var schemaServer = `
		CREATE TABLE server_info (
    		name varchar(64) primary key,
			ip varchar(64),
			ssh_ip varchar(64),
			password varchar(64),
			description varchar(256)
		);
		`

type ServerInfo struct {
	Name        string `db:"name"`
	IP          string `db:"ip"`
	SshIP       string `db:"ssh_ip"`
	Password    string `db:"password"`
	Description string `db:"description"`
}

func QueryServer(name string, db *sqlx.DB) (s ServerInfo, err error) {
	err = db.Get(&s, "SELECT * FROM server_info WHERE name=$1", name)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("query server info in db error: %v", err)
	}
	return s, nil
}

func GetAllServerInfo(db *sqlx.DB) (ss []ServerInfo, err error) {
	ss = []ServerInfo{}
	err = db.Select(&ss, "SELECT * FROM server_info")
	if err != nil {
		return nil, fmt.Errorf("get all server info from db error: %v", err)
	}
	return ss, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func CsvFileToDb(csvFilePath string, db *sqlx.DB) (err error) {
	//先清空内存结构
	servers := make(map[string]ServerInfo)

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
		si := ServerInfo{}
		si.Name = item[0]
		si.IP = item[1]
		si.SshIP = item[2]
		si.Password = item[3]
		si.Description = item[4]

		//在内存中插入
		servers[si.Name] = si

		//在数据库中插入
		//假定table server_info已建立且为空表
		insertServer := `INSERT INTO server_info (name, ip, ssh_ip, password, description) VALUES (:name, :ip, :ssh_ip, :password, :description)`

		_, err := db.NamedExec(insertServer, si)
		if err != nil {
			return fmt.Errorf("insert server basic info into myPostgreSQL errror, UserID = %s :%v", si.Name, err)
		}
	}
	return nil
}

func CreateServerTable(db *sqlx.DB) (result sql.Result, err error) {
	return myPostgreSQL.CreateTable(db, schemaServer)
}

func DropServerTable(db *sqlx.DB) (result sql.Result, err error) {
	return myPostgreSQL.DropTable(db, "server_info")
}
