package main

import (
	"RAS/database"
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

var db *sqlx.DB

func main() {
	//log

	///////////////////////////////////////////////////////////////////////////////////////////////
	driverName, dataSourceName, err := database.LoadPostgreSource()
	if err != nil {
		log.Printf("database.LoadPostgreSource error: %v", err)
		return
	}

	db, err = database.ConnectToDatabase(driverName, dataSourceName)
	if err != nil {
		log.Printf("ConnectToDatabase error: %v", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("db.Ping error: %v", err)
		return
	}

	applyLogin.TheDB = db

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
