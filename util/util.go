package util

import (
	"errors"
	"github.com/thoas/go-funk"
	"sync"
)

type ProcessFunc func(param interface{}) (key interface{}, value interface{})

func ProcessMap(fnParams []interface{}, fn ProcessFunc) (retMap map[interface{}]interface{}, err error) {
	var wg sync.WaitGroup
	var m sync.Mutex
	retMap = make(map[interface{}]interface{})

	if !funk.IsIteratee(fnParams) {
		return nil, errors.New("can not iterate")
	}

	for _, param := range fnParams {
		wg.Add(1)
		go func(param interface{}) {
			k, v := fn(param)
			m.Lock()
			retMap[k] = v
			m.Unlock()
			wg.Done()
		}(param)
	}
	wg.Wait()
	return retMap, nil
}
