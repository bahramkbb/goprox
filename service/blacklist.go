package service

import (
	"log"
	"time"
)

func BlacklistProcessing(){
	for {
		time.Sleep(10 * time.Second)

		log.Print("Processing all Ips to blacklist...")
		for _, val := range CacheClient.GetIpSet("ips") {
			count := CacheClient.GetIPVisitCount(val)

			if count > Configs.RateLimit.Rpm {
				log.Printf("Blacklisting IP: %s, with %d visits.", val, count)

				CacheClient.AddIpToSet(val, "blacklist")
			}
		}

		log.Print("Updating Ip Blacklist...")
		for _, val := range CacheClient.GetIpSet("blacklist") {
			BlackListIPs[val] = true
		}

	}
}