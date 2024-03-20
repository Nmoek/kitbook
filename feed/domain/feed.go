package domain

import (
	"errors"
	"fmt"
	"github.com/ecodeclub/ekit"
	"time"
)

type FeedEvent struct {
	Id int64

	Uid   int64
	Type  string
	Ctime time.Time

	Ext ExtendFields
}

var errKeyNotFound = errors.New("没有找到对应key")

type ExtendFields map[string]string

func (f ExtendFields) Get(key string) ekit.AnyValue {
	val, ok := f[key]
	if !ok {
		return ekit.AnyValue{
			Err: fmt.Errorf("%w, key[%v] \n", errKeyNotFound, key),
		}
	}

	return ekit.AnyValue{Val: val}
}
