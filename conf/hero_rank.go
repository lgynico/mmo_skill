package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type HeroRankRow struct {
	Id             int     `json:"id"`
	Next           int     `json:"next"`
	RankupCond     [][]int `json:"rankupCond"`
	RankupCondSp   [][]int `json:"rankupCondSp"`
	RollbackSelf   int     `json:"rollbackSelf"`
	RollbackSelfSp int     `json:"rollbackSelfSp"`
	RollbackFodder int     `json:"rollbackFodder"`
	MaxLvl         int     `json:"maxLvl"`
	TeamLvl        int     `json:"TeamLvl"`
	AttrCeo        int     `json:"attrCeo"`
	Score          int     `json:"score"`
}

type HeroRank struct {
	Rows map[string]*HeroRankRow
}

// var HeroRank = &HeroRank{}

const (
	CONF_HERO_RANK = "HeroRank"
	FILE_HERO_RANK = "hero_rank.json"
)

func init() {
	ConfMgr.addConfMap(CONF_HERO_RANK, &HeroRank{})
}

func (s *HeroRank) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*HeroRankRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *HeroRank) getFileName() string {
	return FILE_HERO_RANK
}

func (s *HeroRank) GetRowByString(k string) (*HeroRankRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *HeroRank) GetRowByInt(k int) (*HeroRankRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *HeroRank) GetAllRows() map[string]*HeroRankRow {
	return s.Rows
}
