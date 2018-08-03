package main

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func filterTypes(typeStrs []string, opt cmp.Option) cmp.Option {
	typeSet := make(map[string]struct{})
	for _, t := range typeStrs {
		typeSet[t] = struct{}{}
	}
	return cmp.FilterValues(func(a, b interface{}) bool {
		if _, ok := typeSet[reflect.TypeOf(a).String()]; ok {
			return true
		}
		if _, ok := typeSet[reflect.TypeOf(b).String()]; ok {
			return true
		}
		return false
	}, opt)
}

func cmpSym(cmpr func(a, b interface{}) bool) cmp.Option {
	return cmp.Comparer(func(any1, any2 interface{}) bool {
		return cmpr(any1, any2) || cmpr(any2, any1) || cmp.Equal(any1, any2)
	})
}

func myCmpInt(a1 int, any2 interface{}) bool {
	switch a2 := any2.(type) {
	case string:
		return strconv.Itoa(a1) == a2
	}
	return cmp.Equal(a1, any2)
}

func myCmpAny(any1 interface{}, any2 interface{}) bool {
	switch a1 := any1.(type) {
	case int:
		return myCmpInt(a1, any2)
	case float64:
		return myCmpInt(int(a1), any2)
	}
	return false
}

func myCmp() cmp.Option {
	return filterTypes([]string{"float64", "string"}, cmpSym(myCmpAny))
}

func TestGoCmp(t *testing.T) {
	{
		a := map[string]interface{}{
			"foo": 100,
			"bar": "200",
		}
		b := map[string]interface{}{
			"foo": float64(100),
			"bar": 200,
		}
		t.Log(cmp.Diff(a, b))
		assert.False(t, cmp.Equal(a, b))

		opt := myCmp()
		t.Log(cmp.Diff(a, b, opt))
		assert.True(t, cmp.Equal(a, b, opt))
	}
}
