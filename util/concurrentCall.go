package util

import (
	"errors"
	"github.com/thoas/go-funk"
	"reflect"
	"sync"
)

type ConcurrentCallFunc func(input interface{}) (output interface{})

// 切片 fnInputs
//
// 并发执行 fn(input)
//
// num 每次并发数
//
// 返回 fnOutputs
func ConcurrentCall(fnInputs interface{}, num int, fn ConcurrentCallFunc) (fnOutputs []interface{}, err error) {
	inputs := ToInterfaceSlice(fnInputs)
	if funk.IsEmpty(inputs) {
		return nil, errors.New("fnInputs is empty")
	}
	if !funk.IsIteratee(inputs) {
		return nil, errors.New("can not iterate")
	}

	ss := funk.Chunk(inputs, num).([][]interface{})
	for _, s := range ss {
		outputs, err := concurrentCall(s, fn)
		if err != nil {
			return nil, err
		}
		fnOutputs = append(fnOutputs, outputs...)
	}
	return
}

func concurrentCall(inputs []interface{}, fn ConcurrentCallFunc) (fnOutputs []interface{}, err error) {
	var wg sync.WaitGroup
	var m sync.Mutex
	cnt := len(inputs)

	fnOutputs = make([]interface{}, cnt)
	wg.Add(cnt)

	for idx, in := range inputs {
		go func(idx int, in interface{}) {
			defer wg.Done()
			v := fn(in)
			m.Lock()
			defer m.Unlock()
			fnOutputs[idx] = v
		}(idx, in)
	}
	wg.Wait()

	return fnOutputs, nil
}

// 任意 slice 转换为一个 []interface{}
func ToInterfaceSlice(slice interface{}) []interface{} {
	itSlice := make([]interface{}, 0)
	if !funk.IsIteratee(slice) {
		return nil
	}

	var (
		arrValue = reflect.ValueOf(slice)
		arrType  = arrValue.Type()
	)

	if arrType.Kind() == reflect.Slice || arrType.Kind() == reflect.Array {
		for i := 0; i < arrValue.Len(); i++ {
			itSlice = append(itSlice, arrValue.Index(i).Interface())
		}
	}
	return itSlice
}
