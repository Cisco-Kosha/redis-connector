package config

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis/v9"
	"os"
)

type Config struct {
	username  string
	password  string
	db        string
	redisHost string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Redis Username")
	flag.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Redis Password")
	flag.StringVar(&conf.db, "db", os.Getenv("DATABASE"), "Redis Database")
	flag.StringVar(&conf.redisHost, "redisHost", os.Getenv("REDIS_HOST"), "Redis Host URL")

	flag.Parse()

	return conf
}

func (c *Config) GetUsername() string {
	return c.username
}

func (c *Config) GetPassword() string {
	return c.password
}

func (c *Config) GetRedisDb() string {
	return c.db
}

func (c *Config) GetRedisHost() string {
	return c.redisHost
}

func (c *Config) GetRedisClient() *redis.Client {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:6379/%d", c.username, c.password, c.redisHost, 0))
	if err != nil {
		panic(err)
	}

	//return redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})

	return redis.NewClient(opt)
}
