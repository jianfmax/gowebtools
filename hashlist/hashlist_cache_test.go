/**
* @Author: jianfma
* @Date: 2021/2/7
* @Description: TODO
* @Version: 1.0
 */

package hashlist_test

import (
	"fmt"
	"github.com/jianfmax/gowebtools/hashlist"
	"testing"
	"time"
)

func TestHashListCache(t *testing.T) {
	liveTime := 4 * time.Second
	hList := hashlist.New().MaxSize(10000).PurgeOnTime("4s").SetLiveTime(&liveTime).AddDeleteFunc(nil).HashList().New()
	err := hList.Set("hello", "world")
	if err != nil {
		t.Error("添加失败")
		fmt.Println("添加失败")
	}
	fmt.Println("添加成功")
	fmt.Println(hList.Size())
	hList.Set("hello", "world1")
	data, _ := hList.Get("hello")
	fmt.Println(data)
	hList.Set("hello", "world2")
	hList.Set("hello", "world3")
	fmt.Println(hList.Size())
	hList.Set("hello4", "world4")
	hList.Set("hello5", "world5")
	fmt.Println(hList.Size())
	data, _ = hList.Get("hello")
	fmt.Println(data)
	dataDeep := hList.GetThroughKeyValue("hello")
	fmt.Printf("dataDeep: %v\n", dataDeep)
	dataList := hList.GetByFunc("hello", func(k1, k2 interface{}) bool {
		return true
	})
	fmt.Println(dataList)
	dataAllList := hList.GetAll()
	fmt.Println(dataAllList)
	dataAllKeyList := hList.GetAllKey()
	fmt.Println(dataAllKeyList)
	dataRemove := hList.Remove("hello5")
	fmt.Println(dataRemove)
	data = hList.Has("hello5")
	fmt.Printf("Has: %v\n", data)
	data = hList.HasThroughKeyValue("hello5")
	fmt.Printf("HasThroughKeyValue: %v\n", data)
	data = hList.HasByFunc("hello5", func(k1, k2 interface{}) bool {
		return true
	})
	fmt.Printf("HasByFunc: %v\n", data)
	rangeList := make(chan interface{}, 0)
	go hList.Range(func(k, v interface{}, r chan interface{}) bool {
		r <- k
		r <- v
		return true
	}, rangeList)
	for i := 0; i < hList.Size()*2; i++ {
		d := <-rangeList
		fmt.Printf("range %v = %v \n", i, d)
	}
	time.Sleep(4 * time.Second)
	fmt.Println(hList.Size())
	time.Sleep(5 * time.Second)
	fmt.Println(hList.Size())
	dataAllList = hList.GetAll()
	fmt.Printf("dataAllList: %v\n", dataAllList)
}
