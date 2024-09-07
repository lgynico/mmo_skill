package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillLevelRow struct {
	Id          int    `json:"id"`
	Type        int    `json:"type"`
	Events      []int  `json:"events"`
	Range       int    `json:"range"`
	EffDeco1    string `json:"effDeco1"`
	Effect1     int    `json:"effect1"`
	EffDeco2    string `json:"effDeco2"`
	Effect2     int    `json:"effect2"`
	EffDeco3    string `json:"effDeco3"`
	Effect3     int    `json:"effect3"`
	PasEvent1   int    `json:"pasEvent1"`
	PasEffDeco1 string `json:"pasEffDeco1"`
	PasEffect1  int    `json:"pasEffect1"`
	PasEvent2   int    `json:"pasEvent2"`
	PasEffDeco2 string `json:"pasEffDeco2"`
	PasEffect2  int    `json:"pasEffect2"`
	PasEvent3   int    `json:"pasEvent3"`
	PasEffDeco3 string `json:"pasEffDeco3"`
	PasEffect3  int    `json:"pasEffect3"`
}

type SkillLevel struct {
	Rows map[string]*SkillLevelRow
}

// var SkillLevel = &SkillLevel{}

const (
	CONF_SKILL_LEVEL = "SkillLevel"
	FILE_SKILL_LEVEL = "skill_level.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_LEVEL, &SkillLevel{})
}

func (s *SkillLevel) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillLevelRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillLevel) getFileName() string {
	return FILE_SKILL_LEVEL
}

func (s *SkillLevel) GetRowByString(k string) (*SkillLevelRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillLevel) GetRowByInt(k int) (*SkillLevelRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillLevel) GetAllRows() map[string]*SkillLevelRow {
	return s.Rows
}
