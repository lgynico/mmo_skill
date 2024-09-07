package skill

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
)

type SkillMgr struct {
	facade.SkillMgr
}

func NewSkillMgr(cfg *conf.SkillConfigs, unit facade.BattleUnit) facade.SkillMgr {
	return &SkillMgr{
		SkillMgr: facade.NewSkillMgr(unit, SkillCreator, cfg),
	}
}
