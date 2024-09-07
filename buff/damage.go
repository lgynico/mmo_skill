package buff

import (
	"log"

	"github.com/lgynico/mmo_skill/common"

	"github.com/lgynico/mmo_skill/consts"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/proto"
)

type DamageBuff struct {
	facade.BuffHandler
}

func (b *DamageBuff) OnHandleBuff(buff facade.Buff, unit facade.BattleUnit) {
	// 扣血
	value := buff.GetValue()
	hp, sp := unit.TakeDamage(value)
	if sp > 0 {
		log.Printf("[%.2f] 单位 %d 扣除护盾值 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), sp)
	}

	if hp > 0 {
		ctx := unit.GetContext()
		action := &proto.BattleAction{
			Time:  ctx.GetTimeMillis(),
			Unit:  unit.GetId(),
			Key:   int32(consts.ACTION_HURT),
			Value: int64(value),
		}
		ctx.RecordAction(action)

		log.Printf("[%.2f] 单位 %d 扣血 %d 余 %d\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId(), value, unit.GetHp())

		facade.SendEvent(facade.EVENT_TAKE_DAMAGE, unit, value)
		facade.SendEvent(facade.EVENT_GIVE_DAMAGE, buff.GetSrc(), value)

		if unit.IsDead() {
			log.Printf("[%.2f] 单位 %d 死亡\n", float64(unit.GetContext().GetTimeMillis())/1000, unit.GetId())
			action := &proto.BattleAction{
				Time: ctx.GetTimeMillis(),
				Unit: unit.GetId(),
				Key:  int32(consts.ACTION_DAED),
			}
			ctx.RecordAction(action)

			facade.SendEvent(facade.EVENT_KILL_TARGET, buff.GetSrc(), buff.GetSkill())
		}

		skill := buff.GetSkill()
		if skill.GetConfig().IsAtk {
			facade.EnergyRecover(skill.GetUnit(), skill, consts.HitEnergy)
		}

		// 吸血
		b.lifeSteal(buff.GetSrc(), hp, skill)
	}

	// 统计信息
	ctx := unit.GetContext()
	ctx.RecordDamage(buff.GetSrc().GetId(), int32(value))
	ctx.RecordSuffering(unit.GetId(), int32(value))
}

func (b *DamageBuff) OnCalValue(buff facade.Buff, unit facade.BattleUnit) int {
	return common.EnsureBuffValue(buff, unit, true)
}

func (b *DamageBuff) lifeSteal(srcUnit facade.BattleUnit, hp int, skill facade.Skill) {
	lifeStealRate := float64(srcUnit.GetAttrs().Get(consts.LifeSteal)) / 100.0
	if lifeStealRate > 1 {
		lifeStealRate = 1
	}

	addHp := int(float64(hp) * lifeStealRate)

	if addHp > 0 {
		buff := NewBuff(30, skill, 0) // todo
		buff.SetValue(addHp)
		srcUnit.GetBuffMgr().AddBuff(buff)
	}
}
