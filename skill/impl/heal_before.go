package impl

import "github.com/lgynico/mmo_skill/facade"

type HealBefore struct {
	facade.SkillHandler
}

func NewHealBefore(base facade.SkillHandler) facade.SkillHandler {
	return &HealBefore{
		SkillHandler: base,
	}
}

func (h *HealBefore) OnGiveHealFront(event *facade.GameEvent) bool {
	target := event.Params[0].(facade.BattleUnit)
	rate := float64(target.GetHp()) / float64(target.GetMaxHp())
	if rate > 0.2 {
		return false
	}

	buff := event.Params[1].(facade.Buff)

	//hp低于20% 治疗效果提升30%
	currVal := buff.GetValue()
	currVal = currVal + int(float32(currVal)*0.3)
	buff.SetValue(currVal)
	return true
}
