package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type HeroEntryRow struct {
	Id                int             `json:"id"`
	Camp              int             `json:"camp"`
	AbilityType       int             `json:"abilityType"`
	Rank              int             `json:"rank"`
	MaxRank           int             `json:"maxRank"`
	AttackSkill       int             `json:"attackSkill"`
	Skill1            []int           `json:"skill1"`
	Skill2            []int           `json:"skill2"`
	Skill3            []int           `json:"skill3"`
	Skill4            []int           `json:"skill4"`
	InbornPassiveList []int           `json:"inbornPassiveList"`
	Attrs             map[int]float64 `json:"attrs"`
	GrowType          int             `json:"growType"`
	Decomposable      bool            `json:"decomposable"`
	DecompRewards     [][]int         `json:"decompRewards"`
}

type HeroEntry struct {
	Rows map[string]*HeroEntryRow
}

// var HeroEntry = &HeroEntry{}

const (
	CONF_HERO_ENTRY = "HeroEntry"
	FILE_HERO_ENTRY = "hero_entry.json"
)

func init() {
	ConfMgr.addConfMap(CONF_HERO_ENTRY, &HeroEntry{})
}

func (s *HeroEntry) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*HeroEntryRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *HeroEntry) getFileName() string {
	return FILE_HERO_ENTRY
}

func (s *HeroEntry) GetRowByString(k string) (*HeroEntryRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *HeroEntry) GetRowByInt(k int) (*HeroEntryRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *HeroEntry) GetAllRows() map[string]*HeroEntryRow {
	return s.Rows
}
