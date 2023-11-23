package redis

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"time"
//
//	"github.com/go-redis/redis/v8"
//	"medbuddy-backend/internal/config"
//	log "github.com/sirupsen/logrus"
//)
//
//var (
//	Rds *redis.Client
//	Ctx = context.Background()
//)
//
//func SetupRedis() {
//	logger := log.NewLogger()
//	getConfig := config.GetConfig()
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     fmt.Sprintf("%v:%v", getConfig.Redis.Redishost, getConfig.Redis.Redisport),
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//
//	if err := rdb.Ping(Ctx).Err(); err != nil {
//		log.Fatalln("Redis db errors: ", err)
//	}
//	pong, _ := rdb.Ping(Ctx).Result()
//	fmt.Println("Redis says: ", pong)
//	Rds = rdb
//	logger.Info("Redis CONNECTION ESTABLISHED")
//}
//
//type Redis struct {
//	Rdb *redis.Client
//}
//
//func GetRedisDb() *Redis {
//	return &Redis{Rdb: Rds}
//}
//
//func (rdb *Redis) RedisSet(key string, value interface{}) errors {
//	serialized, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	return rdb.Rdb.Set(Ctx, key, serialized, 24*time.Hour).Err()
//}
//
//func (rdb *Redis) RedisGet(key string) ([]byte, errors) {
//	serialized, err := rdb.Rdb.Get(Ctx, key).Bytes()
//	return serialized, err
//}
//
//func (rdb *Redis) RedisDelete(key string) (int64, errors) {
//	deleted, err := rdb.Rdb.Del(Ctx, key).Result()
//	if err != nil {
//		return 0, err
//	}
//	return deleted, nil
//}
