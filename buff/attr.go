package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
)

type AttrBuff struct {
	facade.BuffHandler
}

func (b *AttrBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	// 加数值
	attrs := unit.GetAttrs()
	size := len(buff.GetConfig().Params)
	for i := 0; i < size; i += 3 {
		attrType := buff.GetConfig().Params[i]
		value := buff.GetLv() * buff.GetConfig().Params[i+1]
		bPerc := buff.GetConfig().Params[i+2] == 1

		oldValue := attrs.Get(attrType)
		if bPerc {
			attrs.AddPerc(attrType, value)
			log.Printf("[%.2f] 单位(%d) 属性 %d 增加 %.2f%% 当前值 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), attrType, float64(value)/100.0, attrs.Get(attrType))
		} else {
			attrs.AddExtr(attrType, value)
			log.Printf("[%.2f] 单位(%d) 属性 %d 增加 %d 当前值 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), attrType, value, attrs.Get(attrType))
		}

		if attrType == consts.MaxHp {
			newValue := attrs.Get(attrType)
			if oldValue != newValue {
				diffV := newValue - oldValue
				unit.SetHp(unit.GetHp() + diffV)
			}
		}
	}
}

func (b *AttrBuff) OnRemoveBuff(buff facade.Buff, unit facade.BattleUnit) {
	b.BuffHandler.OnRemoveBuff(buff, unit)
	// 加数值
	attrs := unit.GetAttrs()
	size := len(buff.GetConfig().Params)
	for i := 0; i < size; i += 3 {
		attrType := buff.GetConfig().Params[i]
		value := buff.GetLv() * buff.GetConfig().Params[i+1]
		bPerc := buff.GetConfig().Params[i+2] == 1
		if bPerc {
			attrs.SubPerc(attrType, value)
			log.Printf("[%.2f] 单位(%d) 减少 %d 增加 %.2f%% 当前值 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), attrType, float64(value)/100.0, attrs.Get(attrType))
		} else {
			attrs.SubExtr(attrType, value)
			log.Printf("[%.2f] 单位(%d) 减少 %d 增加 %d 当前值 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), attrType, value, attrs.Get(attrType))
		}
	}
}
