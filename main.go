package main

import (
	"github.com/bahramkbb/goprox/service"
	"log"
	"net/http"
)

func main(){
	log.Print("Loading Server configuration:")
	service.Config(&service.Configs)

	log.Print("Connecting to redis...")
	service.CacheClient = service.RedisClient{}
	service.CacheClient.OpenDB()

	log.Print("Starting webserver...")
	http.HandleFunc("/", service.Proxy)

	go service.BlacklistProcessing()

	log.Fatal(http.ListenAndServe(":8000", nil))
}