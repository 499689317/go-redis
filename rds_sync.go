package rds

import (
	"errors"
	// "sync"
	"time"

	"github.com/499689317/go-log"
	"github.com/go-redis/redis/v7"
)

type Configurable interface {
	Addr() string
	Password() string
	DialTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	PoolSize() int
	PoolTimeout() time.Duration
}
type Client struct {
	// sync.RWMutex
	rdb    *redis.Client
	config Configurable
}

var (
	_client *Client
)

func NewClient(c Configurable) (*Client, error) {

	_client = new(Client)
	_client.rdb = redis.NewClient(&redis.Options{
		Addr:         c.Addr(),
		Password:     c.Password(),
		DialTimeout:  c.DialTimeout(),
		ReadTimeout:  c.ReadTimeout(),
		WriteTimeout: c.WriteTimeout(),
		PoolSize:     c.PoolSize(),
		PoolTimeout:  c.PoolTimeout(),
	})

	e := _client.rdb.Ping().Err()
	if e != nil {
		log.Error().Err(e).Msg("redis connect")
		return nil, errors.New("redis connect failed")
	}
	_client.config = c

	go _workloop()

	log.Info().Str("Addr", c.Addr()).Msg("redis connect ok")
	return _client, nil
}

// Key
func (c *Client) DelSync(keys ...string) error {
	if _client == nil {
		return errors.New("DelSync _client is nil")
	}
	if keys == nil {
		return errors.New("DelSync param error")
	}
	// TODO keys参数会先转为slice，slice作为可变参数时加上...，并且共享同一个内部数组
	return _client.rdb.Del(keys...).Err()
}

// String
func (c *Client) GetSync(key string) (string, error) {
	if _client == nil {
		return "", errors.New("GetSync _client is nil")
	}
	if key == "" {
		return "", errors.New("GetSync param error")
	}
	return _client.rdb.Get(key).Result()
}
func (c *Client) SetSync(key string, value interface{}, expiration time.Duration) error {
	// TODO 如果expiration的值为0，意味着这个键没有超时时长
	if _client == nil {
		return errors.New("SetSync _client is nil")
	}
	if key == "" {
		return errors.New("SetSync param error")
	}
	return _client.rdb.Set(key, value, expiration).Err()
}

// HashMap
func (c *Client) HDelSync(key string, fields ...string) error {
	if _client == nil {
		return errors.New("HDelSync _client is nil")
	}
	if key == "" {
		return errors.New("HDelSync param error")
	}
	return _client.rdb.HDel(key, fields...).Err()
}
func (c *Client) HGetSync(key, field string) (string, error) {
	if _client == nil {
		return "", errors.New("HGetSync _client is nil")
	}
	if key == "" {
		return "", errors.New("HGetSync param error")
	}
	return _client.rdb.HGet(key, field).Result()
}
func (c *Client) HGetAllSync(key string) (map[string]string, error) {
	if _client == nil {
		return nil, errors.New("HGetAllSync _client is nil")
	}
	if key == "" {
		return nil, errors.New("HGetAllSync param error")
	}
	// TODO 根据实际情况返回结果,StringStringMapCmd struct
	return _client.rdb.HGetAll(key).Result()
}
func (c *Client) HMGetSync(key string, fields ...string) ([]interface{}, error) {
	if _client == nil {
		return nil, errors.New("HMGetSync _client is nil")
	}
	if key == "" {
		return nil, errors.New("HMGetSync param error")
	}
	// SliceCmd struct
	return _client.rdb.HMGet(key, fields...).Result()
}
func (c *Client) HSetSync(key, field string, value interface{}) error {
	if _client == nil {
		return errors.New("HSetSync _client is nil")
	}
	if key == "" || field == "" {
		return errors.New("HSetSync param error")
	}
	// BoolCmd struct
	return _client.rdb.HSet(key, field, value).Err()
}
func (c *Client) HMSetSync(key string, fields map[string]interface{}) error {
	if _client == nil {
		return errors.New("HMSetSync _client is nil")
	}
	if key == "" || fields == nil {
		return errors.New("HMSetSync param error")
	}
	// TODO fields为插入的一组或者多组键值对
	return _client.rdb.HMSet(key, fields).Err()
}
func (c *Client) HKeysSync(key string) ([]string, error) {
	if _client == nil {
		return nil, errors.New("HKeysSync _client is nil")
	}
	if key == "" {
		return nil, errors.New("HKeysSync param error")
	}
	// TODO 返回包含map所有key的slice，StringSliceCmd struct
	return _client.rdb.HKeys(key).Result()
}
func (c *Client) HValsSync(key string) ([]string, error) {
	if _client == nil {
		return nil, errors.New("HValsSync _client is nil")
	}
	if key == "" {
		return nil, errors.New("HValsSync param error")
	}
	return _client.rdb.HVals(key).Result()
}
func (c *Client) HLenSync(key string) (int64, error) {
	if _client == nil {
		return 0, errors.New("HLenSync _client is nil")
	}
	if key == "" {
		return 0, errors.New("HLenSync param error")
	}
	return _client.rdb.HLen(key).Result()
}

// List
func (c *Client) LLenSync(key string) (int64, error) {
	if _client == nil {
		return 0, errors.New("LLenSync _client is nil")
	}
	if key == "" {
		return 0, errors.New("LLenSync param error")
	}
	return _client.rdb.LLen(key).Result()
}
func (c *Client) LPopSync(key string) (string, error) {
	if _client == nil {
		return "", errors.New("LPopSync _client is nil")
	}
	if key == "" {
		return "", errors.New("LPopSync param error")
	}
	return _client.rdb.LPop(key).Result()
}
func (c *Client) LPushSync(key string, values ...interface{}) error {
	if _client == nil {
		return errors.New("LPushSync _client is nil")
	}
	if key == "" {
		return errors.New("LPushSync param error")
	}
	return _client.rdb.LPush(key, values...).Err()
}
func (c *Client) RPopSync() {

}
func (c *Client) RPushSync() {

}

// Set, TODO 与List的区别再于Set不能存重复数据
func (c *Client) SAddSync(key string, members ...interface{}) error {
	if _client == nil {
		return errors.New("SAddSync _client is nil")
	}
	if key == "" {
		return errors.New("SAddSync param error")
	}
	return _client.rdb.SAdd(key, members...).Err()
}
func (c *Client) SPopSync(key string) (string, error) {
	if _client == nil {
		return "", errors.New("SPopSync _client is nil")
	}
	if key == "" {
		return "", errors.New("SPopSync param error")
	}
	// StringCmd struct
	return _client.rdb.SPop(key).Result()
}
func (c *Client) SRemSync(key string, members ...interface{}) error {
	if _client == nil {
		return errors.New("SRemSync _client is nil")
	}
	if key == "" {
		return errors.New("SRemSync param error")
	}
	return _client.rdb.SRem(key, members...).Err()
}

// SortSet, TODO score: member
// type Z struct {
// 		Score float64
// 		Member interface{}
// }
func (c *Client) ZAddSync(key string, members ...*redis.Z) error {
	if _client == nil {
		return errors.New("ZAddSync _client is nil")
	}
	if key == "" || members == nil {
		return errors.New("ZAddSync param error")
	}
	return _client.rdb.ZAdd(key, members...).Err()
}
func (c *Client) ZCountSync(key, min, max string) (int64, error) {
	if _client == nil {
		return 0, errors.New("ZCountSync _client is nil")
	}
	if key == "" || min == "" || max == "" {
		return 0, errors.New("ZCountSync param error")
	}
	return _client.rdb.ZCount(key, min, max).Result()
}
func (c *Client) ZRankSync(key, member string) (int64, error) {
	if _client == nil {
		return 0, errors.New("ZRankSync _client is nil")
	}
	if key == "" || member == "" {
		return 0, errors.New("ZRankSync param error")
	}
	// TODO 获取成员在Set中位置，可应用于排行榜
	return _client.rdb.ZRank(key, member).Result()
}
func (c *Client) ZRemSync(key string, members ...interface{}) error {
	if _client == nil {
		return errors.New("ZRemSync _client is nil")
	}
	if key == "" || members == nil {
		return errors.New("ZRemSync param error")
	}
	return _client.rdb.ZRem(key, members...).Err()
}
