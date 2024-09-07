package buff

import "github.com/lgynico/mmo_skill/facade"

type DefaultBuff struct {
	facade.BuffHandler
}

func (b *DefaultBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
}

func (b *DefaultBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {

}

func (b *DefaultBuff) OnCalValue(buff facade.Buff, unit facade.BattleUnit) int {
	return 0
}
