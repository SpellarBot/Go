package redis

import (
	"testing"
	"strconv"
	"fmt"
	"time"
)

func TestRedisClient_Get(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	value, err := client.Get("gotest2")
	if err == nil {
		fmt.Println("value: "+value)
	} else {
		fmt.Printf("err: %v", err)
	}
}

func TestRedisClient_Set(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	_, err := client.Set("BOTTLE"+strconv.FormatInt(time.Now().Unix(), 10), 1, 120)
	if err != nil {
		fmt.Printf("err: %v\t", err)
	}
}

func TestRedisClient_Close(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	client.Close()
	msg, err := client.Set("gotest2", 2, 60)
	fmt.Printf("msg: %v\terr: %v\n", msg, err)
}

func TestRedisClient_Init(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()
}

func TestRedisClient_SetQPS(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	client.Set("gotest", 1, 60)

	for {
		time.Sleep(time.Second)
		fmt.Printf("setQPS: %d\n", client.SetQPS())
	}

}

func TestRedisClient_GetQPS(t *testing.T) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	client.Get("gotest")

	for {
		time.Sleep(time.Second)
		fmt.Printf("getQPS: %d\n", client.GetQPS())
	}

}

//16核单协程920000na/op,2000次
func BenchmarkRedisClient_Set(b *testing.B) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		client.Set("BOTTLE"+strconv.Itoa(i), 1, 60)
	}
}
//16核单协程810000na/op,2000次
func BenchmarkRedisClient_Get(b *testing.B) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		client.Get("BOTTLE"+strconv.Itoa(i))
	}
}
//16核并发度100,47500/op,30000次
func BenchmarkRedisClient_Set2(b *testing.B) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	b.SetParallelism(100)

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			client.Set("BOTTLE"+strconv.Itoa(i), 1, 60)
			i++
		}
	})
	b.StopTimer()

}
//16核并发度100,43500na/op,30000次
func BenchmarkRedisClient_Get2(b *testing.B) {
	client := new(RedisClient)
	client.Bootstrap = []string{"10.48.166.16:6379","10.48.166.17:6379","10.48.166.18:6379"}
	client.Password = ""

	client.Init()

	b.SetParallelism(100)

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			client.Get("BOTTLE"+strconv.Itoa(i))
			i++
		}
	})
	b.StopTimer()

}



