package api

import (
	// "github.com/go-redis/redis"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"encoding/json"
	"errors"
	"time"
)

var (
	Pool   *redis.Pool
)

func InitializeRedisCache() {
	host := "localhost" + ":6379"
	Pool = newRedisPool(host)
}

type RedisQuote struct {
	Price     float64
	Stock     string
	// UserId    string
	// Timestamp int64
	CryptoKey string
}
// func NewRedisClient() {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	pong, err := client.Ping().Result()
// 	glog.Info(pong, err)
// }

func newRedisPool(server string) *redis.Pool {

	return &redis.Pool {
		MaxIdle:     80,
		MaxActive:   10000,
		IdleTimeout: 30 * time.Second,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", server)
			if err != nil {
				glog.Error(err)
				return nil, err
			}

			return conn, err
		},
	}
}

func SetToCache(userId string, qt Quote) error {

	c := Pool.Get()
	defer c.Close()

	q := RedisQuote {
		Price: qt.Price,
		Stock: qt.Stock,
		CryptoKey: qt.CryptoKey,
	}


	quote, err := json.Marshal(q)

	if err != nil {
		glog.Error("Can't marshal the Quote object ", q)
		return errors.New("Can't marshal the Quote object.")
	}

	c.Send("MULTI")
	c.Send("SET", q.Stock, string(quote))
	c.Send("EXPIRE", q.Stock, "60")
	_, err = c.Do("EXEC")

	if err != nil {
		//couldnt set to redis
		glog.Error("Could not SET REDIS value ", err)
		return err
	}

	return nil
}

func GetFromCache(stock string) (Quote, error){
	c := Pool.Get()
	defer c.Close()

	q := Quote{}

	val, err := redis.String(c.Do("GET", stock))

	bytes := ([]byte)(val)

	err = json.Unmarshal(bytes, &q)

	if err != nil {
		glog.Error("Error unmarshling Redis quote ", err)
		return Quote{}, err
	}

	q.Timestamp = getCurrentTs()
	// q.UserId = userId

	glog.Info("Returnig quote from Redis: ", q)

	return q, nil
}