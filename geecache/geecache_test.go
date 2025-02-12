package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"cyw": "661",
	"czy": "661",
	"ccc": "690",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(expect, v) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int)
	gee := NewGroup("score", GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		loadCounts[key]++
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 2<<10)

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but %s got", view)
	}
}
