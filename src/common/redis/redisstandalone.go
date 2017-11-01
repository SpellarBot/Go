// Redis Client
// @Author: Golion
// @Date: 2017.5

package redis

import (
	"fmt"
	"sync/atomic"
	"time"

	goredis "github.com/go-redis/redis"
)

type RedisStandaloneClient struct {
	Addr string
	Password  string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	PoolSize int

	DB        int
	Logger    func(string)

	conn      *goredis.Client
	ticker    *time.Ticker
	getQPS    int32
	getQPSCnt int32
	setQPS    int32
	setQPSCnt int32
}

func (r *RedisStandaloneClient) Init() {
	if r.Logger == nil {
		r.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if r.conn == nil {
		option := &goredis.Options{
			Addr:     r.Addr,
			Password: r.Password, // no password set
			DB:       r.DB,  // use default DB
		}
		if r.ReadTimeout > 0 {
			option.ReadTimeout = r.ReadTimeout
		}
		if r.WriteTimeout > 0 {
			option.WriteTimeout = r.WriteTimeout
		}
		if r.PoolSize > 0 {
			option.PoolSize = r.PoolSize
		}

		r.conn = goredis.NewClient(option)

		pong, err := r.conn.Ping().Result()
		r.Logger(fmt.Sprint("[RedisClient][Init] Ping Redis pong=[%v] err=[%v]", pong, err))
	}
	if r.ticker == nil {
		go r.countQPS()
	}
}

func (r *RedisStandaloneClient) Set(key string, value interface{}, expire int) (bool, error) {
	if r.conn == nil {
		return false, fmt.Errorf("[RedisClient][Set] Error! Connect to Redis Failed! Have You Init()?")
	}
	atomic.AddInt32(&r.setQPSCnt, 1)
	err := r.conn.Set(key, value, time.Duration(expire)*time.Second).Err()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (r *RedisStandaloneClient) IncrBy(key string, incr int64) (int64, error) {
	if r.conn == nil {
		return 0, fmt.Errorf("[RedisClient][Set] Error! Connect to Redis Failed! Have You Init()?")
	}
	atomic.AddInt32(&r.setQPSCnt, 1)
	v, err := r.conn.IncrBy(key, incr).Result()
	if err == nil {
		return v, nil
	} else {
		return 0, err
	}
}

func (r *RedisStandaloneClient) Get(key string) (string, error) {
	if r.conn == nil {
		return "", fmt.Errorf("[RedisClient][Get] Error! Connect to Redis Failed! Have You Init()?")
	}
	atomic.AddInt32(&r.getQPSCnt, 1)
	val, err := r.conn.Get(key).Result()
	if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func (r *RedisStandaloneClient) Expire(key string, expiration time.Duration) (bool, error) {
	if r.conn == nil {
		return false, fmt.Errorf("[RedisClient][Set] Error! Connect to Redis Failed! Have You Init()?")
	}
	v, err := r.conn.Expire(key, expiration).Result()
	return v, err
}

// hash multi set
func (r *RedisStandaloneClient) HMSet(key string, fields map[string]interface{}, expire int) (bool, error) {
	if r.conn == nil {
		return false, fmt.Errorf("[RedisClient][HMSet] Error! Connect to Redis Failed! Have You Init()?")
	}

	atomic.AddInt32(&r.setQPSCnt, 1)
	err := r.conn.HMSet(key, fields).Err()

	var ok = true

	if expire > 0 {
		expireErr := r.conn.PExpire(key, time.Duration(expire)*time.Second).Err()

		if expireErr != nil {
			r.conn.Del(key)
			ok = false
		}
	}

	if ok && err == nil {
		return ok, nil
	} else {
		return false, err
	}
}

// hash field set, if key is not exist will creat and set
func (r *RedisStandaloneClient) HSet(key, field string, value interface{}) (bool, error) {
	if r.conn == nil {
		return false, fmt.Errorf("[RedisClient][HSet] Error! Connect to Redis Failed! Have You Init()?")
	}

	atomic.AddInt32(&r.setQPSCnt, 1)
	err := r.conn.HSet(key, field, value).Err()

	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// hash field incrby, if key is not exist will create and incr, return current value
func (r *RedisStandaloneClient) HIncrBy(key, field string, incr int64) (int64, error) {
	if r.conn == nil {
		return -1, fmt.Errorf("[RedisClient][HIncrBy] Error! Connect to Redis Failed! Have You Init()?")
	}

	atomic.AddInt32(&r.setQPSCnt, 1)
	v, err := r.conn.HIncrBy(key, field, incr).Result()

	if err == nil {
		return v, nil
	} else {
		return -1, err
	}
}
func (r *RedisStandaloneClient) HIncrByFloat(key, field string, incr float64) (float64, error) {
	if r.conn == nil {
		return -1.0, fmt.Errorf("[RedisClient][HIncrBy] Error! Connect to Redis Failed! Have You Init()?")
	}

	atomic.AddInt32(&r.setQPSCnt, 1)
	v, err := r.conn.HIncrByFloat(key, field, incr).Result()

	if err == nil {
		return v, nil
	} else {
		return -1.0, err
	}
}

// hash multi fields get
func (r *RedisStandaloneClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	if r.conn == nil {
		return []interface{}{}, fmt.Errorf("[RedisClient][HMGet] Error! Connect to Redis Failed! Have You Init()?")
	}

	atomic.AddInt32(&r.getQPSCnt, 1)
	val, err := r.conn.HMGet(key, fields...).Result()
	if err != nil {
		return []interface{}{}, err
	} else {
		if val[0] == nil || val[1] == nil {
			return []interface{}{}, fmt.Errorf(fmt.Sprintf("itemid = %s no data in redis.", []byte(key)))
		}
		return val, nil
	}
}

func (r *RedisStandaloneClient) Close() {
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
}

func (r *RedisStandaloneClient) SetQPS() int32 {
	return r.setQPS
}

func (r *RedisStandaloneClient) GetQPS() int32 {
	return r.getQPS
}

func (r *RedisStandaloneClient) countQPS() {
	r.ticker = time.NewTicker(time.Duration(1) * time.Second)
	for _ = range r.ticker.C {
		r.getQPS = r.getQPSCnt
		r.getQPSCnt = 0
		r.setQPS = r.setQPSCnt
		r.setQPSCnt = 0
	}
}

//如果key不存在，返回nil
func (r *RedisStandaloneClient) MGet(keys ...string) ([]interface{}, error) {
	if r.conn == nil {
		return []interface{}{}, fmt.Errorf("[RedisClient][MGet] Error! Connect to Redis Failed! Have You Init()?")
	}
	atomic.AddInt32(&r.getQPSCnt, 1)
	val, err := r.conn.MGet(keys...).Result()
	if err != nil {
		return []interface{}{}, err
	} else {
		return val, nil
	}
}
