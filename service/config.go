package service

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
	"math"
)

var Configs Configuration
var WhiteListIPs map[string]bool
var BlackListIPs map[string]bool
var RequestFrequency int
var CacheClient RedisClient

func Config(conf *Configuration){
	//Initialize Maps
	WhiteListIPs = make(map[string]bool)
	BlackListIPs = make(map[string]bool)


	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	loadConfig(conf)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		loadConfig(conf)
	})
}

func loadConfig(conf *Configuration) {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Printf("Running Https: %t\nCerts Location: %s\nTarget Url: %s\nRPM Limit: %d\nWhitelist: %s\n",
		conf.Server.Https,
		conf.Server.Certs,
		conf.Server.Uri,
		conf.RateLimit.Rpm,
		conf.RateLimit.Whitelist)

	loadWhitelists(conf.RateLimit.Whitelist)
	calculateRequestFrequency(conf.RateLimit.Rpm)
}

func loadWhitelists(ips string) {
	WhiteListIPs = make(map[string]bool)
	for _, ip := range strings.Split(ips, ",") {
		WhiteListIPs[ip] = true
	}
	log.Printf("Whitelist updated: %s", WhiteListIPs)
}

func calculateRequestFrequency(rpm int) {
	RequestFrequency = int(math.Round(60 / float64(rpm) * 1000))
	log.Printf("Requests below %d milliseconds will be processed.", RequestFrequency)
}
