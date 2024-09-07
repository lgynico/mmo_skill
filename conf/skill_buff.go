package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillBuffRow struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Params      []int  `json:"params"`
	IsOuterVal  bool   `json:"isOuterVal"`
	Layer       int    `json:"layer"`
	Duration    int    `json:"duration"`
	EffTimes    int    `json:"effTimes"`
	EffInterval int    `json:"effInterval"`
	Tag         string `json:"tag"`
	SerieTag    string `json:"serieTag"`
	Priority    int    `json:"priority"`
	MutexTag    string `json:"mutexTag"`
	ImmuneTag   string `json:"immuneTag"`
	SubBuff     []int  `json:"subBuff"`
}

type SkillBuff struct {
	Rows map[string]*SkillBuffRow
}

// var SkillBuff = &SkillBuff{}

const (
	CONF_SKILL_BUFF = "SkillBuff"
	FILE_SKILL_BUFF = "skill_buff.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_BUFF, &SkillBuff{})
}

func (s *SkillBuff) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillBuffRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillBuff) getFileName() string {
	return FILE_SKILL_BUFF
}

func (s *SkillBuff) GetRowByString(k string) (*SkillBuffRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillBuff) GetRowByInt(k int) (*SkillBuffRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillBuff) GetAllRows() map[string]*SkillBuffRow {
	return s.Rows
}
