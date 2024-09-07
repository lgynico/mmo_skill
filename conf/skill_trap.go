package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillTrapRow struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Shape      []float64 `json:"shape"`
	TargetType int       `json:"targetType"`
	Delay      int       `json:"delay"`
	Duration   int       `json:"duration"`
	Interval   int       `json:"interval"`
	Times      int       `json:"times"`
	Buffs      []int     `json:"buffs"`
	MoveTrack  int       `json:"moveTrack"`
	MoveParams []int     `json:"moveParams"`
	DieRemove  bool      `json:"dieRemove"`
}

type SkillTrap struct {
	Rows map[string]*SkillTrapRow
}

// var SkillTrap = &SkillTrap{}

const (
	CONF_SKILL_TRAP = "SkillTrap"
	FILE_SKILL_TRAP = "skill_trap.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_TRAP, &SkillTrap{})
}

func (s *SkillTrap) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillTrapRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillTrap) getFileName() string {
	return FILE_SKILL_TRAP
}

func (s *SkillTrap) GetRowByString(k string) (*SkillTrapRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillTrap) GetRowByInt(k int) (*SkillTrapRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillTrap) GetAllRows() map[string]*SkillTrapRow {
	return s.Rows
}
