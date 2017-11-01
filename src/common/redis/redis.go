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

type RedisClient struct {
	Bootstrap []string
	Password  string
	DB        int
	Logger    func(string)

	conn      *goredis.ClusterClient
	ticker    *time.Ticker
	getQPS    int32
	getQPSCnt int32
	setQPS    int32
	setQPSCnt int32
}

func (r *RedisClient) Init() {
	if r.Logger == nil {
		r.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if r.conn == nil {
		r.conn = goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs:    r.Bootstrap,
			Password: r.Password,
		})
		pong, err := r.conn.Ping().Result()
		r.Logger(fmt.Sprint("[RedisClient][Init] Ping Redis pong=[%v] err=[%v]", pong, err))
	}
	if r.ticker == nil {
		go r.countQPS()
	}
}

func (r *RedisClient) Set(key string, value interface{}, expire int) (bool, error) {
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

func (r *RedisClient) IncrBy(key string, incr int64) (int64, error) {
	if r.conn == nil {
		return -1, fmt.Errorf("[RedisClient][Set] Error! Connect to Redis Failed! Have You Init()?")
	}
	atomic.AddInt32(&r.setQPSCnt, 1)
	v, err := r.conn.IncrBy(key, incr).Result()
	if err == nil {
		return v, nil
	} else {
		return -1, err
	}
}

func (r *RedisClient) Get(key string) (string, error) {
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

// hash multi set
func (r *RedisClient) HMSet(key string, fields map[string]interface{}, expire int) (bool, error) {
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
func (r *RedisClient) HSet(key, field string, value interface{}) (bool, error) {
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
func (r *RedisClient) HIncrBy(key, field string, incr int64) (int64, error) {
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
func (r *RedisClient) HIncrByFloat(key, field string, incr float64) (float64, error) {
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
func (r *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
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

func (r *RedisClient) Close() {
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
}

func (r *RedisClient) SetQPS() int32 {
	return r.setQPS
}

func (r *RedisClient) GetQPS() int32 {
	return r.getQPS
}

func (r *RedisClient) countQPS() {
	r.ticker = time.NewTicker(time.Duration(1) * time.Second)
	for _ = range r.ticker.C {
		r.getQPS = r.getQPSCnt
		r.getQPSCnt = 0
		r.setQPS = r.setQPSCnt
		r.setQPSCnt = 0
	}
}
