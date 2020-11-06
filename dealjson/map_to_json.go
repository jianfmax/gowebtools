package dealdataforjson

import "fmt"

func Convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = Convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

func ConvertStr(m map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = Convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
