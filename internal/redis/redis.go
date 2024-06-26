package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedis interface {
	GetClient() *redis.Client
	Connect() error
	Ping(ctx context.Context) error
}

var Redis IRedis

type RedisClient struct {
	Client *redis.Client
	config Config
}

type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	PoolTimeout  int
	IdleTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

func NewRedis(config Config) (IRedis, error) {
	r := &RedisClient{
		config: config,
	}
	if err := r.Connect(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}

func (r *RedisClient) Connect() error {
	Client := redis.NewClient(&redis.Options{
		Addr:         r.config.Addr,
		Password:     r.config.Password,
		DB:           r.config.DB,
		PoolSize:     r.config.PoolSize,
		PoolTimeout:  time.Duration(r.config.PoolTimeout) * time.Second,
		IdleTimeout:  time.Duration(r.config.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(r.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(r.config.WriteTimeout) * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	str, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Print(str)
	r.Client = Client
	return nil
}

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.Client.Ping(ctx).Result()
	return err
}
