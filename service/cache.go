package service

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisClient struct {
	pool *redis.Pool
}

func (rc *RedisClient) OpenDB() {

	// Initialize redis pool
	rc.pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			// Connect to redis
			c, err := redis.Dial("tcp", Configs.Server.RedisUri + ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func (rc *RedisClient) Check() bool {
	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		return false
	} else {
		return true
	}
}

func (rc *RedisClient) SaveVisit(ip string) (bool, error) {
	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("INCR", "ip:" + ip)
	conn.Send("PEXPIRE", "ip:" + ip, RequestFrequency)
	conn.Send("SADD", "visitor_ips", ip)
	res, err := conn.Do("EXEC")

	if err != nil {
		println(res)
		return true, nil
	} else {
		return false, err
	}
}

func (rc *RedisClient) AddIpToSet(ip string, setName string) (bool, error) {
	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	res, err := conn.Do("SADD", setName, ip)

	if err != nil {
		println(res)
		return true, nil
	} else {
		return false, err
	}
}

func (rc *RedisClient) GetSet(setName string) []string {
	var ips []string

	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	ips, err := redis.Strings(conn.Do("SMEMBERS", setName))
	if err != nil {
		log.Printf("error fetching ips from redis : %v", err)
		return nil
	}

	return ips
}

func (rc *RedisClient) EmptySet(setName string) bool {
	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", setName); err == nil {
		return true
	} else {
		log.Printf("error deleting set from redis : %v", err)
		return false
	}

}

func (rc *RedisClient) GetIPVisitCount(ip string) int {
	// Get a connection
	conn := rc.pool.Get()
	defer conn.Close()

	res, err := redis.Int(conn.Do("GET", "ip:" + ip))

	if err != nil {
		log.Printf("error fetching ip visits from redis : %v", err)
		return 0
	}

	return res
}