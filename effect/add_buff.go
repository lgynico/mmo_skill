package effect

import (
	"sort"

	"github.com/lgynico/mmo_skill/buff"
	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
	"github.com/lgynico/mmo_skill/utils"
)

// 赋予buff
type AddBuff struct {
}

func NewAddBuff() facade.SkillEffect {
	return &AddBuff{}
}

func (e *AddBuff) OnEffect(skill facade.Skill, conf *conf.SkillEffectRow, target facade.BattleUnit) {
	buffUnits := make([]facade.BattleUnit, 0)
	if conf.Param3 != common.BUFF_COND_ANY {
		// 筛选目标类型
		targetUnits := make([]facade.BattleUnit, 0)
		ctx := skill.GetUnit().GetContext()
		for _, target := range ctx.GetUnits() {
			if facade.IsTarget(target, skill.GetUnit(), conf.Param2) {
				targetUnits = append(targetUnits, target)
			}
		}
		buffUnits = e.selectTarget(targetUnits, conf.Param3, conf.Param4)

	} else {
		buffUnits = append(buffUnits, target)
	}

	for _, unit := range buffUnits {
		e.addBuff(unit, skill, conf.Param1)
	}

	// unit := skill.GetUnit()
	// if facade.IsTarget(target, unit, conf.Param2) {
	// 	e.addBuff(target, skill, conf.Param1)
	// } else if facade.IsTarget(unit, unit, conf.Param2) {
	// 	e.addBuff(unit, skill, conf.Param1)
	// }
}

func (e *AddBuff) addBuff(unit facade.BattleUnit, skill facade.Skill, buffId int) {
	if unit == nil /**|| unit.IsDead() */ {
		return
	}

	buff := buff.NewBuff(buffId, skill, 0) // todo
	buff.SetValue(buff.CalValue(skill.GetUnit()))
	buffMgr := unit.GetBuffMgr()
	buffMgr.AddBuff(buff)
	// if buff.GetConfig().Type != int(facade.BUFF_TYPE_DAMAGE) {
	//ctx.PushBuff(unit.GetId(), int32(buff.GetConfig().Id), int32(buff.GetValue()))
	// }

}

func (e *AddBuff) selectTarget(targetUnits []facade.BattleUnit, selectType int, count int) []facade.BattleUnit {
	if count >= len(targetUnits) {
		return targetUnits
	}

	switch selectType {
	case common.BUFF_COND_HP_LEAST:
		return e.selectHpLeast(targetUnits, count)
	case common.BUFF_COND_ALL:
		return e.selectAll(targetUnits)
	case common.BUFF_COND_RANDOM:
		return e.selectRandom(targetUnits, count)
	default:
		return e.selectRandom(targetUnits, count)
	}
}

func (e *AddBuff) selectHpLeast(targetUnits []facade.BattleUnit, count int) []facade.BattleUnit {
	sort.Slice(targetUnits, func(i, j int) bool {
		return targetUnits[i].GetHp() < targetUnits[j].GetHp()
	})

	return targetUnits[:count]
}

func (e *AddBuff) selectRandom(targetUnits []facade.BattleUnit, count int) []facade.BattleUnit {
	ints := make([]int, 0, len(targetUnits))
	for i := 0; i < len(targetUnits); i++ {
		ints = append(ints, i)
	}

	indexs := utils.RandIntsNoRepeat(ints, count)
	units := make([]facade.BattleUnit, 0, count)
	for i := 0; i < len(indexs); i++ {
		index := indexs[i]
		units = append(units, targetUnits[index])
	}

	return units
}

func (e *AddBuff) selectAll(targetUnits []facade.BattleUnit) []facade.BattleUnit {
	return targetUnits
}
