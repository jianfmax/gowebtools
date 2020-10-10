/**
* @Author: jfma
* @Date: 2020/10/10
* @Description: 所有用到的模型
* @Version: 1.0
 */

package loadsetting

import (
	"fmt"
	"github.com/pkg/errors"
)

type SectionValueMap struct {
	Section string
	Value   map[string]string
}

func (receiver SectionValueMap) String() string {
	return fmt.Sprintf("Section: %v. Value: %v", receiver.Section, receiver.Value)
}

// Setting 保存设置内容的结构体
type Setting struct {
	ID      uint   `gorm:"primarykey"`
	Section string `gorm:"type:varchar(255);uniqueIndex:idx_name"`
	Key     string `gorm:"type:varchar(255);uniqueIndex:idx_name"`
	Value   string
	Remark  string
}

func SectionListToMap(list []Setting) (SectionValueMap, error) {
	if len(list) == 0 {
		return SectionValueMap{}, nil
	}
	section := list[0].Section
	sectionValueMap := SectionValueMap{
		Section: section,
		Value:   map[string]string{},
	}
	for _, value := range list {
		if section != value.Section {
			return sectionValueMap, errors.WithStack(errors.New("传入的section不一致"))
		}
		sectionValueMap.Value[value.Key] = value.Value
	}
	return sectionValueMap, nil
}
