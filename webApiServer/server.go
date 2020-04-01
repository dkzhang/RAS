package main

import (
	"RAS/myPostgreSQL"
	"RAS/myRedis"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"log"

	//"github.com/sirupsen/logrus"
	"RAS/webApiServer/applyLogin"
	"RAS/webApiServer/queryIP"
	"net/http"
)

//var log = logrus.New()

var dbPostgreSQL *sqlx.DB
var dbRedis *myRedis.Redis

func main() {
	//log

	///////////////////////////////////////////////////////////////////////////////////////////////
	driverName, dataSourceName, err := myPostgreSQL.LoadPostgreSource()
	if err != nil {
		log.Printf("myPostgreSQL.LoadPostgreSource error: %v", err)
		return
	}

	dbPostgreSQL, err = myPostgreSQL.ConnectToDatabase(driverName, dataSourceName)
	if err != nil {
		log.Printf("ConnectToDatabase error: %v", err)
		return
	}
	defer dbPostgreSQL.Close()

	err = dbPostgreSQL.Ping()
	if err != nil {
		log.Printf("dbPostgreSQL.Ping error: %v", err)
		return
	}

	applyLogin.TheDB = dbPostgreSQL

	///////////////////////////////////////////////////////////////////////////////////////////////
	redisHost, redisPasswd, err := myRedis.LoadRedisSource()
	redisOpts := &myRedis.RedisOpts{
		Host:     redisHost,
		Password: redisPasswd,
	}
	dbRedis = myRedis.NewRedis(redisOpts)
	applyLogin.TheRedis = dbRedis
	queryIP.TheRedis = dbRedis
	///////////////////////////////////////////////////////////////////////////////////////////////
	mux := httprouter.New()

	mux.GET("/QueryIP", queryIP.GetIpInfo)
	mux.POST("/ApplyLogin", applyLogin.PostApplyLogin)

	///////////////////////////////////////////////////////////////////////////////////////////////
	htxyServer := http.Server{
		Addr:    "0.0.0.0:8083",
		Handler: mux,
	}
	///////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("The web server is running......")
	//htxyServer.ListenAndServe()
	htxyServer.ListenAndServeTLS("/TLS/2_ras.gribgp.com.crt", "/TLS/3_ras.gribgp.com.key")
}
