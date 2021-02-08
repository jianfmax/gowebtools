package hashlist_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSyncMap(t *testing.T) {
	var scene sync.Map
	// 将键值对保存到sync.Map
	scene.Store("greece", 97)
	scene.Store("london", 100)
	scene.Store("egypt", 200)
	// 从sync.Map中根据键取值
	fmt.Println(scene.Load("london"))
	// 根据键删除对应的键值对
	scene.Delete("london")
	// 遍历所有sync.Map中的键值对
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	go scene.Range(func(k, v interface{}) bool {
		fmt.Println("iterate:", k, v)
		time.Sleep(2 * time.Second)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		return true
	})
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(20 * time.Second)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
