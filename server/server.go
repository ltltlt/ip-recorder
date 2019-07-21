package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ltltlt/ip-recorder/server/middleware"
)

var (
	username string
	password string
	realm    = "work"
)

func init() {
	var ok1, ok2 bool
	username, ok1 = os.LookupEnv("IP_RECORDER_USERNAME")
	password, ok2 = os.LookupEnv("IP_RECORDER_PASSWORD")

	if !ok1 || !ok2 {
		log.Panicf("you may forgot to setup IP_RECORDER_USERNAME and IP_RECORDER_PASSWORD for auth")
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Oh!! You found it"))
}

func main() {
	server := http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Addr:         ":8080",
	}

	http.HandleFunc("/", middleware.PanicRecover(
		middleware.BasicAuth(
			middleware.AccessLog(handle),
			username, password, realm)))

	panic(server.ListenAndServe())
}
