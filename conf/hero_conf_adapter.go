package conf

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/lgynico/mmo_skill/utils"
)

func init() {
	// Must init after skill configs
	// see SkillConfAdapter
}

type SkillConfig struct {
	*SkillEntryRow
	*SkillLevelRowEx
}

func NewSkillConfig(entryCfg *SkillEntryRow, levelCfg *SkillLevelRowEx) *SkillConfig {
	return &SkillConfig{
		SkillEntryRow:   entryCfg,
		SkillLevelRowEx: levelCfg,
	}
}

type SkillConfigs struct {
	Attack  *SkillConfig
	Normals []*SkillConfig
}

type HeroEntryRowEx struct {
	*HeroEntryRow
	RankRow         *HeroRankRow
	CampRow         *HeroCampRow
	AttackSkillCfg  *SkillConfig
	NormalSkillCfgs []*treemap.Map
}

func (cfg *HeroEntryRowEx) GetSkillCfgs(lv int) *SkillConfigs {
	attack := cfg.AttackSkillCfg
	normals := make([]*SkillConfig, 0, len(cfg.NormalSkillCfgs))

	for _, m := range cfg.NormalSkillCfgs {
		k, v := m.Floor(lv)
		if k == nil {
			_, v = m.Min()
		}
		skillCfg := v.(*SkillConfig)
		normals = append(normals, skillCfg)
	}

	return &SkillConfigs{
		Attack:  attack,
		Normals: normals,
	}
}

type HeroRankRowEx struct {
	*HeroRankRow
	PrevRow, NextRow *HeroRankRowEx
}

type HeroGrowConfig struct {
	holder *treemap.Map
	Type   int
}

func (c *HeroGrowConfig) GetTotalGrow(lvl int) float64 {
	grow := float64(1)
	for {
		key, val := c.holder.Floor(lvl)
		if val == nil {
			break
		}

		floorLv := key.(int)
		cfg := val.(*HeroGrowRow)

		grow *= math.Pow(float64(cfg.LvlGrow)/10000, float64(lvl-floorLv+1))
		lvl = floorLv - 1

	}

	return grow
}

type HeroConfAdapter struct {
	*BaseConfAdapter
	HeroRows map[int]*HeroEntryRowEx
	RankRows map[int]*HeroRankRowEx
	GrowRows map[int]*HeroGrowConfig
}

// var HeroConfAdapter = &HeroConfAdapter{
// 	HeroRows: make(map[int]*HeroEntryRowEx),
// 	RankRows: make(map[int]*HeroRankRowEx),
// 	GrowRows: make(map[int]*HeroGrowConfig),
// }

func (a *HeroConfAdapter) RandHeroEntry(camp, rank int) (*HeroEntryRow, bool) {
	heroCfgs := make([]*HeroEntryRow, 0, len(ConfMgr.HeroEntry.Rows))
	for _, heroCfg := range ConfMgr.HeroEntry.Rows {
		if (camp == 0 || camp == heroCfg.Camp) && heroCfg.Rank == rank {
			heroCfgs = append(heroCfgs, heroCfg)
		}
	}

	if len(heroCfgs) == 0 {
		return nil, false
	}

	idx := utils.RandIntByCrypto(0, len(heroCfgs))
	return heroCfgs[idx], true
}

func (a *HeroConfAdapter) onLoadComplete() {
	a.adaptEntry()
	a.adaptRank()
	a.adaptGrow()
}

func (a *HeroConfAdapter) adaptGrow() {
	a.GrowRows = make(map[int]*HeroGrowConfig, len(ConfMgr.HeroGrow.Rows))

	for _, growCfg := range ConfMgr.HeroGrow.Rows {
		cfg, ok := a.GrowRows[growCfg.Type]
		if !ok {
			cfg = &HeroGrowConfig{
				holder: treemap.NewWithIntComparator(),
				Type:   growCfg.Type,
			}
			a.GrowRows[cfg.Type] = cfg
		}

		cfg.holder.Put(growCfg.LvlSection, growCfg)
	}
}

func (a *HeroConfAdapter) adaptEntry() {
	a.HeroRows = make(map[int]*HeroEntryRowEx, len(ConfMgr.HeroEntry.Rows))

	for _, row := range ConfMgr.HeroEntry.Rows {
		rankRow, ok := ConfMgr.HeroRank.GetRowByInt(row.Rank)
		if !ok {
			panic(fmt.Errorf("config not exists: name = %s, row = %d", CONF_HERO_RANK, row.Rank))
		}

		campRow, ok := ConfMgr.HeroCamp.GetRowByInt(row.Camp)
		if !ok {
			panic(fmt.Errorf("config not exists: name = %s, row = %d", CONF_HERO_CAMP, row.Camp))
		}

		rowEx := &HeroEntryRowEx{
			HeroEntryRow: row,
			RankRow:      rankRow,
			CampRow:      campRow,
		}

		a.adaptSkills(rowEx)

		a.HeroRows[row.Id] = rowEx
	}
}

func (a *HeroConfAdapter) adaptSkills(heroCfg *HeroEntryRowEx) {
	// TODO 先过滤木桩
	if heroCfg.HeroEntryRow.Id == 99999 {
		return
	}
	a.fillAttackSkill(heroCfg)
	a.fillNormalSkills(heroCfg)
}

func (a *HeroConfAdapter) fillAttackSkill(heroCfg *HeroEntryRowEx) {
	skillCfg, ok := ConfMgr.SkillConfAdapter.EntryRows[heroCfg.AttackSkill]
	if !ok {
		panic(fmt.Sprintf("config not exists: name=%s, row=%d", CONF_SKILL_ENTRY, heroCfg.AttackSkill))
	}
	heroCfg.AttackSkillCfg = NewSkillConfig(skillCfg.SkillEntryRow, skillCfg.Lv1Row)
}

func (a *HeroConfAdapter) fillNormalSkills(heroCfg *HeroEntryRowEx) {
	heroCfg.NormalSkillCfgs = make([]*treemap.Map, 0)
	if m := a.newSkillMap(heroCfg.Skill1); m != nil {
		heroCfg.NormalSkillCfgs = append(heroCfg.NormalSkillCfgs, m)
	}
	if m := a.newSkillMap(heroCfg.Skill2); m != nil {
		heroCfg.NormalSkillCfgs = append(heroCfg.NormalSkillCfgs, m)
	}
	if m := a.newSkillMap(heroCfg.Skill3); m != nil {
		heroCfg.NormalSkillCfgs = append(heroCfg.NormalSkillCfgs, m)
	}
	if m := a.newSkillMap(heroCfg.Skill4); m != nil {
		heroCfg.NormalSkillCfgs = append(heroCfg.NormalSkillCfgs, m)
	}
}

func (a *HeroConfAdapter) newSkillMap(skill []int) *treemap.Map {
	if len(skill) == 0 {
		return nil
	}

	skillCfg, ok := ConfMgr.SkillConfAdapter.EntryRows[skill[0]]
	if !ok {
		panic(fmt.Sprintf("config not exists: name=%s, row=%d", CONF_SKILL_ENTRY, skill[0]))
	}

	m := treemap.NewWithIntComparator()
	for i := 1; i < len(skill); i++ {
		lv := skill[i]
		var levelCfg *SkillLevelRowEx
		switch i {
		case 1:
			levelCfg = skillCfg.Lv1Row
		case 2:
			levelCfg = skillCfg.Lv2Row
		case 3:
			levelCfg = skillCfg.Lv3Row
		case 4:
			levelCfg = skillCfg.Lv4Row
		}
		m.Put(lv, NewSkillConfig(skillCfg.SkillEntryRow, levelCfg))
	}

	return m
}

func (a *HeroConfAdapter) adaptRank() {
	a.RankRows = make(map[int]*HeroRankRowEx, len(ConfMgr.HeroRank.Rows))

	for _, rankRow := range ConfMgr.HeroRank.Rows {
		rowEx := &HeroRankRowEx{
			HeroRankRow: rankRow,
		}

		a.RankRows[rankRow.Id] = rowEx
	}

	for _, rowEx := range a.RankRows {
		if rowEx.Next != 0 {
			nextRow, ok := a.RankRows[rowEx.Next]
			if !ok {
				panic(fmt.Errorf("config not exists: name=%s, row=%d", CONF_HERO_RANK, rowEx.Next))
			}
			rowEx.NextRow = nextRow
			nextRow.PrevRow = rowEx
		}
	}
}
