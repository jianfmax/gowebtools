package tests

import (
	"github.com/jianfmax/gowebtools/loadsetting"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"testing"
	"time"
)

func TestTmp(t *testing.T) {
	// 加载数据库
	iniFile, err := ini.Load("../conf.ini")
	if err != nil {
		panic(err)
	}
	dataBaseStr, TablePrefix := loadsetting.LoadDatabaseStr(iniFile)
	loadsetting.LoadDatabase(dataBaseStr, TablePrefix)
	err = loadsetting.SaveToDb("../conf.ini", loadsetting.GetDB())
	if err != nil {
		log.Println(err)
	}
	// 监听数据库的变化情况
	output := make(chan loadsetting.SectionValueMap, 100)
	inputData := loadsetting.SaveToSectionValueMap("../conf.ini")
	loadsetting.LoadSetting(loadsetting.GetDB(), time.Duration(10), inputData, output)
	for data := range output {
		log.Println(data.Section)
		log.Println(data.Value)
	}
}
