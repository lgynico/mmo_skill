package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type HeroGrowRow struct {
	Id         int `json:"id"`
	Type       int `json:"type"`
	LvlSection int `json:"lvlSection"`
	LvlGrow    int `json:"lvlGrow"`
}

type HeroGrow struct {
	Rows map[string]*HeroGrowRow
}

// var HeroGrow = &HeroGrow{}

const (
	CONF_HERO_GROW = "HeroGrow"
	FILE_HERO_GROW = "hero_grow.json"
)

func init() {
	ConfMgr.addConfMap(CONF_HERO_GROW, &HeroGrow{})
}

func (s *HeroGrow) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*HeroGrowRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *HeroGrow) getFileName() string {
	return FILE_HERO_GROW
}

func (s *HeroGrow) GetRowByString(k string) (*HeroGrowRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *HeroGrow) GetRowByInt(k int) (*HeroGrowRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *HeroGrow) GetAllRows() map[string]*HeroGrowRow {
	return s.Rows
}
