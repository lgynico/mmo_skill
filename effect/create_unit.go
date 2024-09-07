package effect

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
)

// 创建单元
type CreateUnit struct {
}

func NewCreateUnit() facade.SkillEffect {
	return &CreateUnit{}
}

func (e *CreateUnit) OnEffect(skill facade.Skill, conf *conf.SkillEffectRow, target facade.BattleUnit) {
	// TODO
}
