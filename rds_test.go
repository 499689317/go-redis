package rds

import (
	"math/rand"
	"testing"
	"time"
	"strconv"
)

type conf struct {
	RdsAddr         string
	RdsPassword     string
	RdsDialTimeout  time.Duration
	RdsReadTimeout  time.Duration
	RdsWriteTimeout time.Duration
	RdsPoolSize     int
	RdsPoolTimeout  time.Duration
}

func (c *conf) Addr() string {
	return c.RdsAddr
}
func (c *conf) Password() string {
	return c.RdsPassword
}
func (c *conf) DialTimeout() time.Duration {
	return c.RdsDialTimeout
}
func (c *conf) ReadTimeout() time.Duration {
	return c.RdsReadTimeout
}
func (c *conf) WriteTimeout() time.Duration {
	return c.RdsWriteTimeout
}
func (c *conf) PoolSize() int {
	return c.RdsPoolSize
}
func (c *conf) PoolTimeout() time.Duration {
	return c.RdsPoolTimeout
}

// func TestRedisClientConnect(t *testing.T) {

// 	c := &conf{
// 		RdsAddr:         "127.0.0.1:6379",
// 		RdsPassword:     "",
// 		RdsDialTimeout:  10 * time.Second,
// 		RdsReadTimeout:  30 * time.Second,
// 		RdsWriteTimeout: 30 * time.Second,
// 		RdsPoolSize:     10,
// 		RdsPoolTimeout:  30 * time.Second,
// 	}
// 	_c, e := NewClient(c)
// 	if e != nil {
// 		t.Errorf("NewClient error %v", e)
// 	}
// 	e = _c.SetSync("test", 123, 30 * time.Second)
// 	t.Error(e)
// 	xx, e := _c.GetSync("test")
// 	if e != nil {
// 		t.Error(e)
// 	}
// 	t.Logf("%s", xx)
// }

var _c *Client
func init() {
	c := &conf{
		RdsAddr:         "127.0.0.1:6379",
		RdsPassword:     "",
		RdsDialTimeout:  10 * time.Second,
		RdsReadTimeout:  30 * time.Second,
		RdsWriteTimeout: 30 * time.Second,
		RdsPoolSize:     10,
		RdsPoolTimeout:  30 * time.Second,
	}
	_c, _ = NewClient(c)
	
}
func BenchmarkRedisClientMethods(b *testing.B) {

	for i := 0; i < b.N; i++ {
		// rand.Seed(time.Now().Unix())
		
		s := rand.New(rand.NewSource(time.Now().Unix()))
		
		// b.Logf("%d", s.Intn(100))
		k := strconv.Itoa(s.Int())
		_c.HSetSync("testmap", k, s.Int())
	}
}
