package dealdataforjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// Read 读取一个文件中的内容.
func Read(fileName string) []byte {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return f
}

// GetMap 从文件中读取json文件为map.
func GetMap(fileName string) map[string]interface{} {
	bytes := Read(fileName)
	v := make(map[string]interface{})
	err := json.Unmarshal(bytes, &v)
	if err != nil {
	}
	return v
}

// GetDataFromIndex 根据索引得到值
func GetDataFromIndex(timeTemplate map[string]interface{}, indexMap []string) (string, error) {
	var val = timeTemplate
	for index, value := range indexMap[:len(indexMap)-1] {
		if tmp, ok := val[value].(map[string]interface{}); ok {
			val = tmp
		} else {
			return "", errors.New("没有找到该字段" + indexMap[index])
		}
	}
	return val[indexMap[len(indexMap)-1]].(string), nil
}

// ChangeDataFromAIndex 根据索引修改值,只一个值
func ChangeDataFromAIndex(val map[string]interface{}, indexMap []string, timeTemplate string,
	fun func(timeStr, timeTemplate string) (string, error)) error {
	for index, value := range indexMap[:len(indexMap)-1] {
		if tmp, ok := val[value].(map[string]interface{}); ok {
			val = tmp
		} else {
			return errors.New("没有找到该字段" + indexMap[index])
		}
	}
	var data string
	if tmp, ok := val[indexMap[len(indexMap)-1]].(string); ok {
		data = tmp
	} else {
		return errors.New("没有找到该字段" + indexMap[len(indexMap)-1])
	}
	newStr, err := fun(data, timeTemplate)
	if err != nil {
		return err
	}
	val[indexMap[len(indexMap)-1]] = newStr
	return nil
}

// ChangeDataFromIndex 根据索引修改值的属性
func ChangeDataFromIndex(bytes []byte,
	indexMap map[string]string, fun func(timeStr, timeTemplate string) (string, error)) ([]byte, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}
	for key, value := range indexMap {
		indexStr := strings.Split(key, ".")
		err = ChangeDataFromAIndex(v, indexStr, value, fun)
		if err != nil {
			return nil, err
		}
	}
	newBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return newBytes, nil
}

// TimeToDefault 将不标准的字符串转为标准字符串RFC3339Nano
func TimeToDefault(timeStr, timeTemplate string) (string, error) {
	value, err := time.ParseInLocation(timeTemplate, timeStr, time.Local)
	if err != nil {
		return "", err
	}
	return value.Format(time.RFC3339Nano), nil
}
