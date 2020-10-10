/**
* @Author: jfma
* @Date: 2020/10/9
* @Description: 计算两个map[string]string 是否相同
* @Version: 1.0
 */

package loadsetting

import (
	"fmt"
	"sort"
)

// MapSimilar 两个map 中的内容是否相同
func MapSimilar(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	data1 := make([]string, 0)
	for k, v := range map1 {
		data1 = append(data1, fmt.Sprintf("%v:%v", k, v))
	}
	data2 := make([]string, 0)
	for k, v := range map2 {
		data2 = append(data2, fmt.Sprintf("%v:%v", k, v))
	}
	sort.Strings(data1)
	sort.Strings(data2)
	for index, data := range data1 {
		if data != data2[index] {
			return false
		}
	}
	return true
}
