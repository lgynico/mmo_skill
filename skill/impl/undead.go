package impl

import "github.com/lgynico/mmo_skill/facade"

type UnDead struct {
	facade.SkillHandler
	isTrigger bool
}

func NewUnDead(base facade.SkillHandler) facade.SkillHandler {
	return &UnDead{
		SkillHandler: base,
		isTrigger:    false,
	}
}

func (h *UnDead) OnTakeDamage(event *facade.GameEvent) bool {
	if h.isTrigger {
		return true
	}

	// 判断是否死亡了
	unit := event.Target.(facade.BattleUnit)
	if unit.IsDead() {
		unit.TakeHealing(1)
		h.isTrigger = true
		return false
	}

	return true
}
