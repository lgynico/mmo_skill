package buff

import (
	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/proto"
)

type EnergyBuff struct {
	facade.BuffHandler
}

func (b *EnergyBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	// 扣血
	value := buff.GetValue()
	maxEp := unit.GetAttrs().Get(consts.MaxEnergy)
	change := unit.AddEp(value, maxEp, true)

	if change != 0 {
		ctx := unit.GetContext()
		action := &proto.BattleAction{
			Time:  ctx.GetTimeMillis(),
			Unit:  unit.GetId(),
			Key:   int32(consts.ACTION_ENERGY_CHG),
			Value: int64(value),
		}
		ctx.RecordAction(action)
	}
}

func (b *EnergyBuff) OnCalValue(buff facade.Buff, unit facade.BattleUnit) int {
	return common.EnsureBuffValue(buff, unit, false)
}
