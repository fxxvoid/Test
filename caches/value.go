package caches

import (
	"asterism/helpers"
	"sync/atomic"
	"time"
)

const (
	NeverDie = 0
)

type value struct {
	Data []byte
	//寿命
	Ttl int64
	//创建时间
	Ctime int64
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		Data: helpers.Copy(data),
		Ttl: ttl,
		Ctime: time.Now().Unix(),
	}
}

func (v *value) alive() bool {
	return v.Ttl == NeverDie || time.Now().Unix()-v.Ctime<v.Ttl
}

func (v *value) visit() []byte {
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}