package dealjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestChangeData(t *testing.T) {
	data := Read("tmp.json")
	fmt.Println(string(data))
	indexMap := make(map[string]string)
	indexMap["value.go_value.times"] = "2006-01-02 15:04:05"
	indexMap["new_time"] = "2006-01-02 15:04:05"
	newData, err := ChangeDataFromIndex(data, indexMap, TimeToDefault)
	fmt.Println(string(newData))
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetKeyValue(t *testing.T) {
	f, _ := os.Open("tmp.json")
	data, _ := ioutil.ReadAll(f)
	dataJson := map[string]interface{}{}
	err := json.Unmarshal(data, &dataJson)
	if err != nil {
		log.Println("error!")
	}
	key := make([]string, 0)
	result := make([]JsonKeyValue, 0)
	tmp, err := getKeyAndValue(dataJson, key, result)
	if err != nil {
		return
	}
	fmt.Println(tmp)
	getValue, err := GetValueByKey(dataJson, []string{"last_name", "first_name"})
	if err != nil {
		log.Println(err)
	}
	log.Println(getValue)
	err = SetValueByKeyAdd(dataJson, []string{"last_name", "first_nameA"}, 2)
	if err != nil {
		log.Println(err)
	}
	getValue, err = GetValueByKey(dataJson, []string{"last_name", "first_name"})
	if err != nil {
		log.Println(err)
	}

	log.Println(getValue)
	newJson, err := json.Marshal(dataJson)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(newJson))
}
