package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type MonsterEntryRow struct {
	Id     int    `json:"id"`
	Lv     int    `json:"Lv"`
	Name   string `json:"name"`
	Rank   int    `json:"rank"`
	HeroId int    `json:"heroId"`
}

type MonsterEntry struct {
	Rows map[string]*MonsterEntryRow
}

// var MonsterEntry = &MonsterEntry{}

const (
	CONF_MONSTER_ENTRY = "MonsterEntry"
	FILE_MONSTER_ENTRY = "monster_entry.json"
)

func init() {
	ConfMgr.addConfMap(CONF_MONSTER_ENTRY, &MonsterEntry{})
}

func (s *MonsterEntry) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*MonsterEntryRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *MonsterEntry) getFileName() string {
	return FILE_MONSTER_ENTRY
}

func (s *MonsterEntry) GetRowByString(k string) (*MonsterEntryRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *MonsterEntry) GetRowByInt(k int) (*MonsterEntryRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *MonsterEntry) GetAllRows() map[string]*MonsterEntryRow {
	return s.Rows
}
