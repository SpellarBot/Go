package utils
import (
	"fmt"
	"github.com/go-redis/redis"
	"errors"
)

type RedisClient struct {
	Client *redis.Client
}

func (r *RedisClient) Init(addr string, password string,db int){
	options := redis.Options{Addr:addr,Password:password,DB:db}
	r.Client = redis.NewClient(&options)
}

func (r *RedisClient) Get(key string)(interface{},error){
	return r.Client.Get(key).Val(),nil
}
func (r *RedisClient) Set(key string, value interface{})(error){
	e := r.Client.Set(key,value,0)
	return errors.New(e.Val())
}

func RedisTest(){
	R := RedisClient{}
	R.Init("localhost:6379","",0)
	R.Client.Set("name","zhangsanfeng",30)
	name := R.Client.Get("name").Val()
	fmt.Println(len(name),name)
}
