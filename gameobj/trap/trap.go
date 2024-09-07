package trap

import (
	"github.com/lgynico/mmo_skill/conf"
	"github.com/lgynico/mmo_skill/facade"
)

type triggerCounter struct {
	triggerTimes     int
	triggerTimestamp int64
}

type Trap struct {
	facade.Trap
	createTime              int // 创建时间
	triggerTime             int // 触发时间
	mapObjId2TriggerCounter map[int32]*triggerCounter
	// lastEffTime             int // 最后生效时间
	// effTimesRemain          int // 剩余生效次数
}

func NewTrap(uid int32, config *conf.SkillTrapRowEx, unit facade.BattleUnit) facade.Trap {
	trap := &Trap{
		Trap:                    facade.NewTrap(uid, unit.GetRoleId(), config),
		mapObjId2TriggerCounter: make(map[int32]*triggerCounter),
		// effTimesRemain: config.Times,
		// targetEffTimes: make(map[int]int),
	}

	// trap.AddComponet(facade.NewFSM(trap, InitState))
	trap.SetFSM(facade.NewFSM(trap, InitState))
	trap.SetUnit(unit)
	trap.SetContext(unit.GetContext())
	return trap
}
