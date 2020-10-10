/**
* @Author: jfma
* @Date: 2020/10/09 周五
* @Description: 主函数
* @Version: 1.0
**/

package loadsetting

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// SettingList 所有的设置内容
type SettingList struct {
	settingList []Setting
}

// DataBase 数据库相关的设置
type DataBase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
}

var db *gorm.DB

func LoadDatabase(DataBaseStr, TablePrefix string) {
	var err error
	db, err = gorm.Open(mysql.Open(DataBaseStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: true,
		},
	}) // 连接数据库，设置表名的前缀
	if err != nil {
		log.Error(err)
	}
	sqlDB, err := db.DB() // 设置数据连接池的基本信息
	if err != nil {
		log.Error(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&Setting{}) // 自动迁移数据表
	if err != nil {
		log.Error(err)
	}
}

// GetDB 得到数据库连接
func GetDB() *gorm.DB {
	return db
}

// LoadDatabaseStr 加载mysql数据库
func LoadDatabaseStr(cfg *ini.File) (DataBaseStr string, TablePrefix string) {
	DataBaseSetting := DataBase{}
	err := cfg.Section("database").MapTo(&DataBaseSetting)
	if err != nil {
		panic("加载数据库出错")
	}
	DataBaseStr = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DataBaseSetting.User,
		DataBaseSetting.Password,
		DataBaseSetting.Host,
		DataBaseSetting.Port,
		DataBaseSetting.Name)
	fmt.Println(DataBaseStr)
	return DataBaseStr, DataBaseSetting.TablePrefix
}

// CloseDB 关闭数据库
func CloseDB() {
	sqlDb, _ := db.DB()
	err := sqlDb.Close()
	if err != nil {
		log.Error(err)
	}
}
