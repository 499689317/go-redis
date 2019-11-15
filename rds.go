package rds

import (
	"time"

	"github.com/499689317/go-log"
	"github.com/go-redis/redis/v7"
)

var (
	_workqueue = make(chan *work, 2048)
)

const (
	_ = iota
	Key_Del
	String_Set
	Hash_Del
	Hash_Set
	Hash_MSet
	List_LPop
	List_LPush
	List_RPop
	List_RPush
	Set_Add
	Set_Pop
	Set_Rem
	SortSet_Add
	SortSet_Rem
)

type work struct {
	o int
	k string
	f string
	v interface{}
	m map[string]interface{}
	z *redis.Z
	e time.Duration
}

// 消费者，_workloop可以在多个go程中执行
func _workloop() {
	log.Info().Msg("open new goroutine as consumer")
	var e error
	for w := range _workqueue {
		if _client == nil {
			log.Error().Msg("_loop _client is nil")
			_workqueue <- w // TODO w还给queue, 重新排队
			continue
		}

		switch w.o {
		case Key_Del:
			e = _client.DelSync(w.k)
		case String_Set:
			e = _client.SetSync(w.k, w.v, w.e)
		case Hash_Set:
			e = _client.HSetSync(w.k, w.f, w.v)
		case Hash_MSet:
			e = _client.HMSetSync(w.k, w.m)
		case Hash_Del:
			e = _client.HDelSync(w.k)
		case List_LPush:
			e = _client.LPushSync(w.k, w.v)
		case Set_Add:
			e = _client.SAddSync(w.k, w.v)
		case Set_Rem:
			e = _client.SRemSync(w.k)
		case SortSet_Add:
			e = _client.ZAddSync(w.k, w.z)
		case SortSet_Rem:
			e = _client.ZRemSync(w.k)
		}

		if e != nil {
			log.Error().Err(e).Msg("work执行出错")
			log.Warn().Int("w.o", w.o).Str("w.k", w.k).Msg("记录work执行错误log")
			_workqueue <- w
		}
	}
}

// work生产者，TODO 如果_workqueue阻塞，需要超时机制释放资源
func (c *Client) PushWork(o int, k string, v interface{}) {
	if o == 0 || k == "" {
		log.Error().Msg("pushwork param error")
		return
	}
	if v == nil {
		// v为空时，检查对应操作是否正确，记录log
		log.Warn().Int("o", o).Str("k", k).Msg("work value is nil, its ok?")
	}
	select {
	case _workqueue <- &work{o: o, k: k, v: v}:
		log.Info().Int("o", o).Str("k", k).Msg("pushwork is ok")
	case <-time.After(time.Second * 1):
		log.Warn().Int("o", o).Str("k", k).Msg("pushwork is timeout 1s")
	}
}
