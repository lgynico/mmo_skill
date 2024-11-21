package conf

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
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

type HeroConfAdapter struct {
	*BaseConfAdapter
	HeroRows map[int]*HeroEntryRowEx
}

// var HeroConfAdapter = &HeroConfAdapter{
// 	HeroRows: make(map[int]*HeroEntryRowEx),
// 	RankRows: make(map[int]*HeroRankRowEx),
// 	GrowRows: make(map[int]*HeroGrowConfig),
// }

func (a *HeroConfAdapter) onLoadComplete() {
	a.adaptEntry()
}

func (a *HeroConfAdapter) adaptEntry() {
	a.HeroRows = make(map[int]*HeroEntryRowEx, len(ConfMgr.HeroEntry.Rows))

	for _, row := range ConfMgr.HeroEntry.Rows {

		rowEx := &HeroEntryRowEx{
			HeroEntryRow: row,
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
