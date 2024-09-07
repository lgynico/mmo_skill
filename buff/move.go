package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/geom"
)

type MoveBuff struct {
	facade.BuffHandler
}

func (b *MoveBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	log.Printf("[%.2f] 单位 %d 受到技能位移 %s\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), buff.GetSkill().GetConfig().Name)

	buffCfg := buff.GetConfig()
	heading := b.getHeading(buffCfg.Params[0], unit)
	speed := float64(buffCfg.Params[1]) / 100
	time := int64(buffCfg.Duration)

	unit.SetSkillMove(true)
	unit.SetSkillMoveSpeed(speed)
	unit.SetSkillMoveHeading(heading)
	unit.SetSkillMoveTime(time)

	unit.AddStateModifier(facade.NewAttackeStateModifier(time))
	unit.AddStateModifier(facade.NewMoveStateModifier(time))
	unit.AddStateModifier(facade.NewSpellStateModifier(time))
	unit.AddStateModifier(facade.NewUniqueSpellStateModifier(time))
}

func (b *MoveBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {
	b.BuffHandler.OnRemoveBuff(buff, unit)
}

func (b *MoveBuff) getHeading(direction int, unit facade.BattleUnit) *geom.Vector2d {
	if direction == 0 {
		return unit.GetHeading().Copy()
	} else {
		return unit.GetHeading().Inverse()
	}
}
