package service

import (
	"log"
	"time"
)

func BlacklistProcessing(){
	for {
		time.Sleep(10 * time.Second)

		log.Print("Getting all Ips...")

		test := CacheClient.GetAllIps()

		log.Print(test)
	}
}