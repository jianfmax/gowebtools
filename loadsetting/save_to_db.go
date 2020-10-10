/**
* @Author: jfma
* @Date: 2020/10/09 周五
* @Description: 主函数
* @Version: 1.0
**/

package loadsetting

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

// SaveToDb 将文件的内容保存到数据库中
func SaveToDb(fileName string, db *gorm.DB) error {
	log.Println("将ini文件中的数据保存在数据库中！")
	Cfg, err := ini.Load(fileName)
	if err != nil {
		panic(err)
	}
	sectionList := Cfg.SectionStrings()
	for _, section := range sectionList {
		log.Println(section)
		keyList := Cfg.Section(section).KeyStrings()
		for _, key := range keyList {
			value := Cfg.Section(section).Key(key).String()
			log.Println(key + ": " + value)
			SaveOneSetting(Setting{
				ID:      0,
				Section: section,
				Key:     key,
				Value:   value,
				Remark:  "",
			}, db)
		}
	}
	return nil
}

// SaveOneSetting 保存一个设置
func SaveOneSetting(data Setting, db *gorm.DB) {
	tmp := db.Where("`section` = ? and `key` = ?", data.Section, data.Key).First(&Setting{})
	if tmp.RowsAffected >= 1 {
		return
	}
	result := db.Create(&data)
	if result.Error != nil {
		log.Error(result.Error)
	}
}

// SaveToSectionValueMap 将ini文件中的数据保存在[]*SectionValueMap中
func SaveToSectionValueMap(fileName string) []*SectionValueMap {
	log.Println("将ini文件中的数据保存在数据库中！")
	Cfg, err := ini.Load(fileName)
	if err != nil {
		panic(err)
	}
	data := make([]*SectionValueMap, 0)
	sectionList := Cfg.SectionStrings()
	for _, section := range sectionList {
		newValue := &SectionValueMap{
			Section: section,
			Value:   map[string]string{},
		}
		keyList := Cfg.Section(section).KeyStrings()
		for _, key := range keyList {
			value := Cfg.Section(section).Key(key).String()
			log.Println(key + ": " + value)
			newValue.Value[key] = value
		}
		data = append(data, newValue)
	}
	return data
}
