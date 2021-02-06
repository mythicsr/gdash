package util

import (
	"fmt"
	"github.com/thoas/go-funk"
	"testing"
	"time"
)

func TestConcurrentCall(t *testing.T) {
	var arr []int
	for i := 0; i < 1000; i++ {
		arr = append(arr, 1)
	}

	stime := time.Now()
	results, err := ConcurrentCall(arr, 100, func(input interface{}) (output interface{}) {
		k := input.(int)
		v := Addd(k)
		_ = v
		return v
	})

	if err != nil {
		panic(err)
	}

	vv := 0
	for _, v := range results {
		if v == 2 {
			vv++
		}
	}

	fmt.Println(vv, time.Since(stime).Seconds())
}

func Addd(a int) int {
	sec := time.Duration(funk.RandomInt(0, 100))
	time.Sleep(sec * time.Millisecond)
	return a + 1
}
