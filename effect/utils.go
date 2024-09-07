package effect

import (
	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/facade"
)

func NewSkillEffect(typ int) (facade.SkillEffect, bool) {
	if eff, ok := mapType2SkillEffect[typ]; ok {
		return eff, true
	}

	return nil, false
}

var mapType2SkillEffect = map[int]facade.SkillEffect{
	common.SK_EFF_TYP_CREATE_PROJECTILE: NewCreateProjectile(),
	common.SK_EFF_TYP_CREATE_TRAP:       NewCreateTrap(),
	common.SK_EFF_TYP_CREATE_UNIT:       NewCreateUnit(),
	common.SK_EFF_TYP_ADD_BUFF:          NewAddBuff(),
	common.SK_EFF_TYP_CHANGE_POS:        NewChangePos(),
}
