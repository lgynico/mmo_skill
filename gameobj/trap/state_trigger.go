package trap

import (
	"log"

	"github.com/lgynico/mmo_skill/common"
	"github.com/lgynico/mmo_skill/facade"
)

var TriggerState facade.FSMState = nil
var TriggerToEndTransition facade.FSMTransition = nil

func init() {
	TriggerState = &triggerState{facade.NewFSMState()}

	TriggerToEndTransition = &triggerToEndTransition{facade.NewFSMTransition()}
	TriggerState.AddTransition(TriggerToEndTransition)

}

type triggerState struct {
	facade.FSMState
}

func (s *triggerState) OnEnter(obj facade.GameObject, dt int64) {
	log.Printf("[%.2f] 陷阱触发: %s\n", float64(obj.GetContext().GetTimeMillis())/1000, obj.(*Trap).GetConfig().Name)
}

func (s *triggerState) OnUpdate(obj facade.GameObject, dt int64) int64 {
	trap := obj.(*Trap)
	ctx := trap.GetContext()

	// 延迟
	if int64(trap.triggerTime) > ctx.GetTimeMillis() {
		return 0
	}

	// 陷阱移动
	if handler, ok := common.GetTrapMoveHandler(trap.GetConfig().MoveTrack); ok {
		handler.HandleTrapMove(trap, ctx, dt)
	}

	shape := trap.GetBoundingShape(true)
	config := trap.GetConfig()
	for _, unit := range ctx.GetUnits() {
		if unit.IsDead() || !facade.IsTarget(trap.GetUnit(), unit, config.TargetType) {
			continue
		}

		counter, ok := trap.mapObjId2TriggerCounter[unit.GetId()]
		if !ok {
			counter = &triggerCounter{}
			trap.mapObjId2TriggerCounter[unit.GetId()] = counter
		}

		if counter.triggerTimes >= config.Times {
			continue
		}

		if counter.triggerTimestamp+int64(config.Interval) > ctx.GetTimeMillis() {
			continue
		}

		// if !geom.TwoShapeCollision(shape, unit.GetBoundingShape(true)) {
		// 	continue
		// }
		if !shape.IsPointInside(unit.GetPosition()) {
			continue
		}

		// log.Printf("技能命中目标(%d)\n", unit.GetId())
		s.addBuff(trap, unit, ctx)
		counter.triggerTimes++
		counter.triggerTimestamp = ctx.GetTimeMillis()
	}

	return 0
}

func (s *triggerState) addBuff(trap *Trap, unit facade.BattleUnit, ctx facade.Battle) {
	buffMgr := unit.GetBuffMgr()
	for _, buffId := range trap.GetConfig().Buffs {
		// buff := buff.NewBuff(buffId, trap.GetSkill(), trap.GetValue())
		buff := ctx.CreateBuff(buffId, trap.GetSkill(), trap.GetValue())
		buff.SetStartTime(int(ctx.GetTimeMillis()))
		buff.SetValue(buff.CalValue(unit))
		buffMgr.AddBuff(buff)
		// facade.PushAddBuff(ctx.GetMServer(), unit.GetAreaId(), unit, int32(buffId), int32(buff.GetValue()))
	}
}

// =====================================================
// -> 结束
// =====================================================
type triggerToEndTransition struct {
	facade.FSMTransition
}

func (t *triggerToEndTransition) IsValid(obj facade.GameObject) bool {
	trap := obj.(*Trap)
	ctx := trap.GetContext()
	return int64(trap.triggerTime+trap.GetConfig().Interval) >= ctx.GetTimeMillis()
}

func (t *triggerToEndTransition) NextState() facade.FSMState { return EndState }
