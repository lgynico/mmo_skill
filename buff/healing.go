package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/facade"
)

type HealingBuff struct {
	facade.BuffHandler
}

func (b *HealingBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	// 加血
	facade.SendEvent(facade.EVENT_GIVE_HEAL_FRONT, buff.GetSrc(), unit, buff)
	value := buff.GetValue()
	unit.TakeHealing(value)

	currHp := unit.GetHp()
	overVal := value - (unit.GetMaxHp() - currHp)
	if overVal > 0 {
		facade.SendEvent(facade.EVENT_OVER_HEAL, buff.GetSrc(), unit, overVal)
	}

	ctx := unit.GetContext()
	facade.SendEvent(facade.EVENT_TAKE_HEAL, unit, ctx, value)
	facade.SendEvent(facade.EVENT_GIVE_HEAL, buff.GetSrc(), ctx, value)
	//ctx.PushHpChange(unit.GetId(), int32(value))
	log.Printf("[%.2f] 单位 %d 加血 %d 总 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), value, unit.GetHp())

	// 统计信息
	ctx.RecordHealing(buff.GetSrc().GetId(), int32(value))
}

func (b *HealingBuff) OnCalValue(buff facade.Buff, unit facade.BattleUnit) int {
	return common.EnsureBuffValue(buff, unit, false)
}
