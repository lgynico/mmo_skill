package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillEntryRow struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Unique        bool   `json:"unique"`
	CoolDown      int    `json:"coolDown"`
	EnterCoolDown int    `json:"enterCoolDown"`
	TargetType    int    `json:"targetType"`
	Lv1           int    `json:"lv1"`
	Lv2           int    `json:"lv2"`
	Lv3           int    `json:"lv3"`
	Lv4           int    `json:"lv4"`
	IsAtk         bool   `json:"isAtk"`
	Range         int    `json:"range"`
	SingTime      int    `json:"singTime"`
	ProcessTime   int    `json:"processTime"`
	ShakeTime     int    `json:"shakeTime"`
}

type SkillEntry struct {
	Rows map[string]*SkillEntryRow
}

// var SkillEntry = &SkillEntry{}

const (
	CONF_SKILL_ENTRY = "SkillEntry"
	FILE_SKILL_ENTRY = "skill_entry.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_ENTRY, &SkillEntry{})
}

func (s *SkillEntry) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillEntryRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillEntry) getFileName() string {
	return FILE_SKILL_ENTRY
}

func (s *SkillEntry) GetRowByString(k string) (*SkillEntryRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillEntry) GetRowByInt(k int) (*SkillEntryRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillEntry) GetAllRows() map[string]*SkillEntryRow {
	return s.Rows
}
