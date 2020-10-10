/**
* @Author: jfma
* @Date: 2020/10/09 周五
* @Description: 主函数
* @Version: 1.0
**/

package loadsetting

import (
	log "github.com/sirupsen/logrus"
	"time"

	"gorm.io/gorm"
)

// LoadSetting 加载数据的设定
func LoadSetting(db *gorm.DB, reloadTime time.Duration, inputData []*SectionValueMap, outputData chan SectionValueMap) {
	log.SetLevel(log.DebugLevel)
	log.Info("开始进入监听更改情况")
	reloadAllSetting(db, inputData, outputData)
	ticker := time.NewTicker(time.Second * reloadTime)
	go func() {
		for _ = range ticker.C {
			reloadAllSetting(db, inputData, outputData)
		}
	}()
}

// reloadAllSetting 重新加载所有配置
func reloadAllSetting(db *gorm.DB, inputData []*SectionValueMap, outputData chan SectionValueMap) {
	for _, data := range inputData {
		reloadOneSetting(db, data, outputData)
	}
}

// reloadOneSetting 重新加载某个设置
func reloadOneSetting(db *gorm.DB, sectionMap *SectionValueMap, outputData chan SectionValueMap) {
	sqlData := make([]Setting, 0)
	db.Where("section = ? ", sectionMap.Section).Find(&sqlData)
	newValueMap, err := SectionListToMap(sqlData)
	if err != nil {
		log.Error("出现错误，无法将setting类型转化为Map类型")
		return
	}
	if !MapSimilar(sectionMap.Value, newValueMap.Value) {
		outputData <- newValueMap
		*sectionMap = newValueMap
		log.Debug(sectionMap.Section, ". 设置内容有更改，新内容: ", sectionMap, ".")
	} else {
		log.Debug(sectionMap.Section, ". 设置内容无更改.")
	}
}
