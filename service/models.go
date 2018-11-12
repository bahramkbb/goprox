package service

type Configuration struct {
	RateLimit RateLimitConfiguration
	Server ServerConfiguration
}

type RateLimitConfiguration struct {
	Rpm int
	Whitelist string
}

type ServerConfiguration struct {
	Uri string
	RedisUri string
}
