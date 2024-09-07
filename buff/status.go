package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/facade"
)

type StatusBuff struct {
	facade.BuffHandler
	stateModifiers []facade.StateModifier
}

func NewStatusBuff(handler facade.BuffHandler) facade.BuffHandler {
	return &StatusBuff{
		BuffHandler:    handler,
		stateModifiers: make([]facade.StateModifier, 0),
	}
}

func (b *StatusBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	stat := buff.GetConfig().Params[0]
	log.Printf("[%.2f] 单位 %d 设置状态：%d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), stat)
	b.addModifier(unit, stat, facade.BUFF_STAT_MOVE_BAN, -1, facade.NewMoveStateModifier)
	b.addModifier(unit, stat, facade.BUFF_STAT_ATTACK_BAN, -1, facade.NewAttackeStateModifier)
	b.addModifier(unit, stat, facade.BUFF_STAT_SKILL_BAN, -1, facade.NewSpellStateModifier)
	b.addModifier(unit, stat, facade.BUFF_STAT_CHARM, -1, facade.NewCharmStateModifier)
	b.addModifier(unit, stat, facade.BUFF_STAT_FINAL_SKILL_BAN, -1, facade.NewUniqueSpellStateModifier)
}

func (b *StatusBuff) addModifier(unit facade.BattleUnit, stat, test int, time int64, applyFunc func(int64) facade.StateModifier) {
	if (stat & test) == test {
		modifier := applyFunc(time)
		unit.AddStateModifier(modifier)
		b.stateModifiers = append(b.stateModifiers, modifier)
	}
}

func (b *StatusBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {
	b.BuffHandler.OnRemoveBuff(buff, unit)
	log.Printf("[%.2f] 单位 %d 移除状态：%d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), buff.GetConfig().Params[0])
	// unit.UnsetBuffState(buff.GetConfig().Params[0])
	for _, modifier := range b.stateModifiers {
		unit.RemoveStateModifier(modifier)
	}
}
