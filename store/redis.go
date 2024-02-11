package store

import (
	"context"
	"dna-test/config"
	"dna-test/models"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClientReal struct {
	Client         *redis.Client
	PubSub         *redis.PubSub
	IsEnablePubSub bool
}

func NewRedisConn(config *config.Config) *RedisClientReal {
	return &RedisClientReal{
		Client: redis.NewClient(&redis.Options{
			Addr:     config.CacheHost,
			Password: config.CachePass,
			DB:       0,
		}),
		IsEnablePubSub: true,
	}
}

func NewRedis(ctx context.Context, config *config.Config, out *models.CacheClient) (err error) {
	cacheClient := NewRedisConn(config)

	_, err = cacheClient.Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	if out != nil && cacheClient != nil {
		*out = cacheClient
	}
	return nil
}

func (r *RedisClientReal) SetHash(ctx context.Context, key string, objectType interface{}, expTime time.Duration) error {
	value, err := json.Marshal(objectType)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, value, expTime).Err()
}

func (r *RedisClientReal) GetHash(ctx context.Context, key string, out interface{}) error {
	result, err := r.Client.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}
	return json.Unmarshal(result, &out)
}

func (r *RedisClientReal) GetKeys(ctx context.Context, pattern string, out *[]string) error {
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}

	*out = append(*out, keys...)
	return nil
}

const (
	ListStart = 0
	ListEnd   = -1
)

func UnMarshalStruct(in interface{}, out interface{}) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

func (r *RedisClientReal) GetList(ctx context.Context, list string, out interface{}) error {
	elements, err := r.Client.LRange(ctx, list, ListStart, ListEnd).Result()
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}

	var results []interface{}
	for _, elem := range elements {
		var result interface{}
		err = json.Unmarshal([]byte(elem), &result)
		if err != nil {
			return err
		}

		results = append(results, result)
	}

	return UnMarshalStruct(results, out)
}

func (r *RedisClientReal) PushToList(ctx context.Context, key string, objectType interface{}) error {
	value, err := json.Marshal(objectType)
	if err != nil {
		return err
	}
	return r.Client.RPush(ctx, key, value).Err()
}

func (r *RedisClientReal) DeleteHash(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClientReal) PublishMsg(ctx context.Context, topic string, msg interface{}) error {
	if !r.IsEnablePubSub {
		return nil
	}
	return r.Client.Publish(ctx, topic, msg).Err()
}

func (r *RedisClientReal) SubScribe(ctx context.Context, topics ...string) error {
	if !r.IsEnablePubSub {
		return nil
	}
	r.PubSub = r.Client.Subscribe(ctx, topics...)
	_, err := r.PubSub.Receive(ctx)
	return err
}

func (r *RedisClientReal) ReceiveMsg(ctx context.Context, out *redis.Message) error {
	if r.PubSub == nil {
		return errors.New("no subscriber is provided")
	}
	msg, err := r.PubSub.ReceiveMessage(ctx)
	if err != nil {
		return err
	}
	if out != nil && msg != nil {
		*out = *msg
	}
	return nil
}

func (r *RedisClientReal) CloseSub(ctx context.Context) error {
	if !r.IsEnablePubSub {
		return nil
	}
	return r.PubSub.Unsubscribe(ctx)
}

func (r *RedisClientReal) CloseClient() error {
	return r.Client.Close()
}
