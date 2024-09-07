package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type HeroLvlRow struct {
	Id        int     `json:"id"`
	Materials [][]int `json:"materials"`
	TeamLvl   int     `json:"teamLvl"`
}

type HeroLvl struct {
	Rows map[string]*HeroLvlRow
}

// var HeroLvl = &HeroLvl{}

const (
	CONF_HERO_LVL = "HeroLvl"
	FILE_HERO_LVL = "hero_lvl.json"
)

func init() {
	ConfMgr.addConfMap(CONF_HERO_LVL, &HeroLvl{})
}

func (s *HeroLvl) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*HeroLvlRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *HeroLvl) getFileName() string {
	return FILE_HERO_LVL
}

func (s *HeroLvl) GetRowByString(k string) (*HeroLvlRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *HeroLvl) GetRowByInt(k int) (*HeroLvlRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *HeroLvl) GetAllRows() map[string]*HeroLvlRow {
	return s.Rows
}
