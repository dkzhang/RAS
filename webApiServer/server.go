package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	//"github.com/sirupsen/logrus"
	"RAS/webApiServer/applyLogin"
	"RAS/webApiServer/queryIP"
	"net/http"
)

//var log = logrus.New()

func main() {
	//log

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
