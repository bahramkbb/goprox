package service

import (
	"log"
	"time"
)

func BlacklistProcessing(){
	for {
		time.Sleep(10 * time.Second)

		log.Print("IP processing started...")

		//Empty the old black list for checking
		CacheClient.EmptySet("blacklist")

		for _, val := range CacheClient.GetSet("visitor_ips") {

			if _, ok := WhiteListIPs[val]; ok {
				continue
			}

			count := CacheClient.GetIPVisitCount(val)

			if count > Configs.RateLimit.Rpm {

				if _, ok := BlackListIPs[val]; ok {
					log.Printf("Permanently Blocking: %s", val)
					CacheClient.AddIpToSet(val, "permanent_blacklist")
					continue
				}

				CacheClient.AddIpToSet(val, "blacklist")
			}
		}

		BlackListIPs = make(map[string]bool)
		PermanentBlackListIPs = make(map[string]bool)

		log.Print("Blacklist IPs:")
		for _, val := range CacheClient.GetSet("blacklist") {
			BlackListIPs[val] = true
		}
		log.Print(BlackListIPs)

		log.Print("Permanent Blacklist IPs:")
		for _, val := range CacheClient.GetSet("permanent_blacklist") {
			PermanentBlackListIPs[val] = true
		}
		log.Print(PermanentBlackListIPs)

	}
}