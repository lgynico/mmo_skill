package conf

import (
	"encoding/json"
	"os"
	"strconv"
)

type SkillProjectileRow struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Shape      []int   `json:"shape"`
	Speed      float64 `json:"speed"`
	Time       int     `json:"time"`
	TargetType int     `json:"targetType"`
	EffDeco1   string  `json:"effDeco1"`
	Effect1    int     `json:"effect1"`
	EffDeco2   string  `json:"effDeco2"`
	Effect2    int     `json:"effect2"`
	EffDeco3   string  `json:"effDeco3"`
	Effect3    int     `json:"effect3"`
}

type SkillProjectile struct {
	Rows map[string]*SkillProjectileRow
}

// var SkillProjectile = &SkillProjectile{}

const (
	CONF_SKILL_PROJECTILE = "SkillProjectile"
	FILE_SKILL_PROJECTILE = "skill_projectile.json"
)

func init() {
	ConfMgr.addConfMap(CONF_SKILL_PROJECTILE, &SkillProjectile{})
}

func (s *SkillProjectile) load(fullPath string) {
	// 读取配置文件
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s.Rows = make(map[string]*SkillProjectileRow)
	// 解析json并转换成对象
	err = json.NewDecoder(f).Decode(&s.Rows)

	if err != nil {
		panic(err)
	}
}

func (s *SkillProjectile) getFileName() string {
	return FILE_SKILL_PROJECTILE
}

func (s *SkillProjectile) GetRowByString(k string) (*SkillProjectileRow, bool) {
	row, ok := s.Rows[k]
	return row, ok
}

func (s *SkillProjectile) GetRowByInt(k int) (*SkillProjectileRow, bool) {
	return s.GetRowByString(strconv.Itoa(k))
}

func (s *SkillProjectile) GetAllRows() map[string]*SkillProjectileRow {
	return s.Rows
}
