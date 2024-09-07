package impl

import "github.com/lgynico/mmo_skill/facade"

type OverHealing struct {
	facade.SkillHandler
}

func NewOverHealing(base facade.SkillHandler) facade.SkillHandler {
	return &OverHealing{
		SkillHandler: base,
	}
}

func (h *OverHealing) OnOverHeal(event *facade.GameEvent) bool {
	// 过量的50%转化为护盾
	target := event.Params[0].(facade.BattleUnit)
	val := event.Params[1].(int) / 2

	ctx := target.GetContext()
	buff := ctx.CreateBuff(6320, h.GetSkill(), val)
	target.GetBuffMgr().AddBuff(buff)
	return true
}
