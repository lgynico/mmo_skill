package impl

import "github.com/lgynico/mmo_skill/facade"

type BurningRage struct {
	facade.SkillHandler
	isTrigger bool
}

func NewBurningRage(base facade.SkillHandler) facade.SkillHandler {
	return &BurningRage{
		SkillHandler: base,
		isTrigger:    false,
	}
}

func (s *BurningRage) OnKillTarget(event *facade.GameEvent) bool {
	if s.isTrigger {
		return true
	}

	skill := event.Params[0].(facade.Skill)
	if skill != s.GetSkill() {
		return true
	}

	s.isTrigger = true
	return false
}
