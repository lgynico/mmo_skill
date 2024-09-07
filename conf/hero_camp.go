package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type HeroCampRow struct {
	Id            int   `json:"id"`
	Restrained    int   `json:"restrained"`
	Refrained     int   `json:"refrained"`
	IsSpecial     bool  `json:"isSpecial"`
	RollbackTrans []int `json:"rollbackTrans"`
	RankType      int   `json:"rankType"`
}

type HeroCamp struct {
	Rows map[string]*HeroCampRow
}

// var HeroCamp = &HeroCamp{}

const (
	CONF_HERO_CAMP = "HeroCamp"
	FILE_HERO_CAMP = "hero_camp.json"
)

func init() {
	ConfMgr.addConfMap(CONF_HERO_CAMP, &HeroCamp{})
}

func (s *HeroCamp) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*HeroCampRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *HeroCamp) getFileName() string {
	return FILE_HERO_CAMP
}

func (s *HeroCamp) GetRowByString(k string) (*HeroCampRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *HeroCamp) GetRowByInt(k int) (*HeroCampRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *HeroCamp) GetAllRows() map[string]*HeroCampRow {
	return s.Rows
}
