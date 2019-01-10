package redis

import (
	"github.com/astaxie/goredis"
)

type RedisR interface {
	RedisRpush(key string, value string) error
	RedisRpop(key string) ([]byte, error)
	RedisLLen(key string) (int, error)
	RedisLrange(key string, start int, end int) ([]string, error)
	RedisLpop(key string) ([]byte, error)
	RedisSet(key string, value string) error
	RedisGet(key string) (string, error)
}

type Redis struct{
	host string
	db int
}

type RedisClient struct{
	goredis.Client
}

var GlobalRedis = Redis{"127.0.0.1:6379", 0}
var GlobalRedisClient RedisClient

func init(){
	GlobalRedisClient.Addr = GlobalRedis.host
	GlobalRedisClient.Db = GlobalRedis.db
}

func (*RedisClient) RedisRpush(key string, value string) error{
	err :=GlobalRedisClient.Rpush(key, []byte(value))
	if err != nil{
		return err
	}
	return nil
}

func (*RedisClient) RedisLpop(key string) ([]byte, error){
	data, err :=GlobalRedisClient.Lpop(key)
	if err != nil{
		return nil, err
	}
	return data, nil
}

func (*RedisClient) RedisRpop(key string) ([]byte, error){
	data, err :=GlobalRedisClient.Rpop(key)
	if err != nil{
		return nil, err
	}
	return data, nil
}

func (*RedisClient) RedisLLen(key string) (int, error){
	data, err :=GlobalRedisClient.Llen(key)
	if err != nil{
		return -1, err
	}
	return data, nil
}

func (*RedisClient) RedisLrange(key string, start int, end int) (resp []string, err error){
	data, err :=GlobalRedisClient.Lrange(key, start, end)
	if err != nil{
		return nil, err
	}
	for _, value := range data{
		resp = append(resp, string(value))
	}
	return resp, nil
}

func (*RedisClient) RedisSet(key string, value string) error{
	err :=GlobalRedisClient.Set(key, []byte(value))
	if err != nil{
		return err
	}
	return nil
}

func (*RedisClient) RedisGet(key string) (string, error){
	data, err :=GlobalRedisClient.Get(key)
	if err != nil{
		return "", err
	}
	return string(data), nil
}

