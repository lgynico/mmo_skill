package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillEffectRow struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Param1 int    `json:"param1"`
	Param2 int    `json:"param2"`
	Param3 int    `json:"param3"`
	Param4 int    `json:"param4"`
	Param5 int    `json:"param5"`
}

type SkillEffect struct {
	Rows map[string]*SkillEffectRow
}

// var SkillEffect = &SkillEffect{}

const (
	CONF_SKILL_EFFECT = "SkillEffect"
	FILE_SKILL_EFFECT = "skill_effect.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_EFFECT, &SkillEffect{})
}

func (s *SkillEffect) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillEffectRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillEffect) getFileName() string {
	return FILE_SKILL_EFFECT
}

func (s *SkillEffect) GetRowByString(k string) (*SkillEffectRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillEffect) GetRowByInt(k int) (*SkillEffectRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillEffect) GetAllRows() map[string]*SkillEffectRow {
	return s.Rows
}
