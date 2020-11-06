package dealdataforjson

import (
	"errors"
	"log"
	"reflect"
)

// 得到所有的key:value

// JsonKeyValue 保存所有的键值对
type JsonKeyValue struct {
	Key   []string
	Value interface{}
	Type  string
}

// GetAllKeyValue 得到所有的key Value
func GetAllKeyValue(mapObj map[string]interface{}) ([]JsonKeyValue, error) {
	result, err := getKeyAndValue(ConvertStr(mapObj), make([]string, 0), make([]JsonKeyValue, 0))
	return result, err
}

func getKeyAndValue(mapObj interface{}, nowKey []string, allResult []JsonKeyValue) ([]JsonKeyValue, error) {
	result := allResult
	switch mapValue := mapObj.(type) {
	case string, int, float32, float64, bool, byte:
		log.Printf("key: %v; valueTyep: %v; value: %v", nowKey, reflect.TypeOf(mapObj), mapValue)
		result = append(result, JsonKeyValue{
			Key:   nowKey,
			Value: mapValue,
			Type:  (reflect.TypeOf(mapObj)).String(),
		})
	case []interface{}:
		for _, value := range mapValue {
			tmp, err := getKeyAndValue(value, nowKey, allResult)
			if err != nil {
				return result, err
			}
			result = append(result, tmp...)
		}
	case map[string]interface{}:
		for key, value := range mapValue {
			lastKey := append(nowKey, key)
			tmpMap, err := getKeyAndValue(value, lastKey, allResult)
			if err != nil {
				return result, err
			}
			result = append(result, tmpMap...)
		}
	default:
		return result, errors.New("出错啦")
	}
	return result, nil
}

// GetValueByKey 得到一个key下的值
func GetValueByKey(mapObj interface{}, keyArray []string) (JsonKeyValue, error) {
	if len(keyArray) == 0 {
		return JsonKeyValue{}, errors.New("不存在该值")
	}
	if len(keyArray) == 1 {
		switch mapKey := mapObj.(type) {
		case string, int, float32, float64, bool, byte:
			return JsonKeyValue{}, errors.New("不存在该值")
		case []interface{}:
			for _, value := range mapKey {
				jsonKeyValue, err := GetValueByKey(value, keyArray)
				if err != nil {
					continue
				}
				return jsonKeyValue, nil
			}
			return JsonKeyValue{}, errors.New("不存在该值")
			/*switch arrayValue := value.(type) {
				case map[string]interface{}:
					if tmp, ok := arrayValue[keyArray[0]]; ok {
						return JsonKeyValue{Value: tmp, Type: reflect.TypeOf(tmp).String()}, nil
					}
					return JsonKeyValue{}, errors.New("不存在该值")
				default:
					return JsonKeyValue{}, errors.New("不存在该值")
				}
			}
			return JsonKeyValue{}, errors.New("不存在该值")*/
		case map[string]interface{}:
			if tmp, ok := mapKey[keyArray[0]]; ok {
				return JsonKeyValue{Value: tmp, Type: reflect.TypeOf(tmp).String()}, nil
			}
			return JsonKeyValue{}, errors.New("不存在该值")
		}
	}
	switch mapKey := mapObj.(type) {
	case string, int, float32, float64, bool, byte:
		return JsonKeyValue{}, errors.New("不存在该值")
	case []interface{}:
		for _, value := range mapKey {
			jsonKeyValue, err := GetValueByKey(value, keyArray)
			if err != nil {
				continue
			}
			return jsonKeyValue, nil
		}
		return JsonKeyValue{}, errors.New("不存在该值")
		/*for _, value := range mapKey {
			switch arrayValue := value.(type) {
			case map[string]interface{}:
				if tmp, ok := arrayValue[keyArray[0]]; ok {
					return JsonKeyValue{Value: tmp, Type: reflect.TypeOf(tmp).String()}, nil
				}
				return JsonKeyValue{}, errors.New("不存在该值")
			default:
				return JsonKeyValue{}, errors.New("不存在该值")
			}
		}
		return JsonKeyValue{}, errors.New("不存在该值")*/
	case map[string]interface{}:
		if tmp, ok := mapKey[keyArray[0]]; ok {
			jsonKeyValue, err := GetValueByKey(tmp, keyArray[1:])
			if err != nil {
				return JsonKeyValue{}, errors.New("不存在该值")
			}
			return jsonKeyValue, nil
		}
		return JsonKeyValue{}, errors.New("不存在该值")
	}
	return JsonKeyValue{}, errors.New("不存在该值")
}

// SetValueByKey 设置一个key的value, 如果不存在这个key，就直接报错。
func SetValueByKey(mapObj interface{}, keyArray []string, newValue interface{}) error {
	if len(keyArray) == 0 {
		return errors.New("设置该值出错")
	}
	if len(keyArray) == 1 {
		switch mapKey := mapObj.(type) {
		case string, int, float32, float64, bool, byte:
			return errors.New("设置该值出错")
		case []interface{}:
			for _, value := range mapKey {
				err := SetValueByKey(value, keyArray, newValue)
				if err != nil {
					continue
				}
				return nil
			}
			return errors.New("设置该值出错")
		case map[string]interface{}:
			if _, ok := mapKey[keyArray[0]]; ok {
				mapKey[keyArray[0]] = newValue
				return nil
			}
			return errors.New("设置该值出错")
		}
	}
	switch mapKey := mapObj.(type) {
	case string, int, float32, float64, bool, byte:
		return errors.New("设置该值出错")
	case []interface{}:
		for _, value := range mapKey {
			err := SetValueByKey(value, keyArray, newValue)
			if err != nil {
				continue
			}
			return nil
		}
		return errors.New("设置该值出错")
	case map[string]interface{}:
		if tmp, ok := mapKey[keyArray[0]]; ok {
			err := SetValueByKey(tmp, keyArray[1:], newValue)
			if err != nil {
				return errors.New("设置该值出错")
			}
			return nil
		}
		return errors.New("设置该值出错")
	}
	return errors.New("设置该值出错")
}

// SetValueByKeyAdd 设置一个key的value, 如果只有最后一个key不存在，则添加最后一个值，如果其余的key不存在，则报错。
func SetValueByKeyAdd(mapObj interface{}, keyArray []string, newValue interface{}) error {
	if len(keyArray) == 0 {
		return errors.New("设置该值出错")
	}
	if len(keyArray) == 1 {
		switch mapKey := mapObj.(type) {
		case string, int, float32, float64, bool, byte:
			return errors.New("设置该值出错")
		case []interface{}:
			for _, value := range mapKey {
				err := SetValueByKeyAdd(value, keyArray, newValue)
				if err != nil {
					continue
				}
				return nil
			}
			return errors.New("设置该值出错")
		case map[string]interface{}:
			mapKey[keyArray[0]] = newValue
			return nil
		}
	}
	switch mapKey := mapObj.(type) {
	case string, int, float32, float64, bool, byte:
		return errors.New("设置该值出错")
	case []interface{}:
		for _, value := range mapKey {
			err := SetValueByKeyAdd(value, keyArray, newValue)
			if err != nil {
				continue
			}
			return nil
		}
		return errors.New("设置该值出错")
	case map[string]interface{}:
		if tmp, ok := mapKey[keyArray[0]]; ok {
			err := SetValueByKeyAdd(tmp, keyArray[1:], newValue)
			if err != nil {
				return errors.New("设置该值出错")
			}
			return nil
		}
		return errors.New("设置该值出错")
	}
	return errors.New("设置该值出错")
}
