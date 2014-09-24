package models

import (
	"fmt"
	"testing"
)

// 测试[]int是否可以转成[]interface{}
func TestArrayInt(t *testing.T) {
	var (
		v, vv interface{}
		ai    []int = []int{1, 2, 3}
		ok    bool
	)
	v = ai
	if vv, ok = v.([]interface{}); ok {
		fmt.Println("Convert []int to []interface{}")
		fmt.Println(vv)
	} else {
		fmt.Println("Cannot convert []int to []interface{}")
	}
}
