package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillHaloRow struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Radius     float64 `json:"radius"`
	Cycle      int     `json:"cycle"`
	Buff       int     `json:"buff"`
	TargetType int     `json:"targetType"`
}

type SkillHalo struct {
	Rows map[string]*SkillHaloRow
}

// var SkillHalo = &SkillHalo{}

const (
	CONF_SKILL_HALO = "SkillHalo"
	FILE_SKILL_HALO = "skill_halo.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_HALO, &SkillHalo{})
}

func (s *SkillHalo) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillHaloRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillHalo) getFileName() string {
	return FILE_SKILL_HALO
}

func (s *SkillHalo) GetRowByString(k string) (*SkillHaloRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillHalo) GetRowByInt(k int) (*SkillHaloRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillHalo) GetAllRows() map[string]*SkillHaloRow {
	return s.Rows
}
