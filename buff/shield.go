package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/facade"
)

type ShieldBuff struct {
	facade.BuffHandler
}

func (b *ShieldBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	// 护盾
	value := buff.GetValue()
	unit.AddSp(buff.GetId(), value)

	// ctx := unit.GetContext()
	// action := &proto.BattleAction{
	// 	Time:  int32(ctx.GetTimeMillis()),
	// 	Unit:  unit.GetId(),
	// 	Key:   int32(proto.ActionType_ACT_HURT),
	// 	Value: int64(value),
	// }
	// ctx.RecordAction(action)

	log.Printf("[%.2f] 单位 %d 增加护盾 %d %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), buff.GetId(), value)
}

func (b *ShieldBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {
	b.BuffHandler.OnRemoveBuff(buff, unit)
	unit.RemoveSp(buff.GetId())
	log.Printf("[%.2f] 单位 %d 移除护盾 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), buff.GetId())
}

func (b *ShieldBuff) OnCalValue(buff facade.Buff, unit facade.BattleUnit) int {
	return common.EnsureBuffValue(buff, unit, false)
}
